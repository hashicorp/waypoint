package lambda

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/datadir"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
)

type BuilderConfig struct {
	Runtime string `hcl:"runtime"`
	Setup   string `hcl:"setup"`
}

type Builder struct {
	config BuilderConfig

	runtime string
	sum     string
	preRef  string

	id string

	preZip, appZip, libZip string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Config() interface{} {
	return &b.config
}

func (b *Builder) BuildFunc() interface{} {
	return b.Build
}

var (
	ErrMissingRuntime  = errors.New("missing runtime indicator file")
	ErrInvalidRuntime  = errors.New("invalid runtime")
	ErrUnexpectedError = errors.New("unexpected error")
)

var supportedRuntimes = map[string]struct{}{
	"ruby2.5": struct{}{},
}

// These are paths that we want to remove to try to make the build more reproducible and smaller.
// NOTE: The following paths must be left in to keep thing works properly:
// * ruby/2.5.0/extensions/x86_64-linux/.*/gem.build_complete
var pruneRegexp = []*regexp.Regexp{
	regexp.MustCompile("ruby/2.5.0/extensions/x86_64-linux/.*/gem_make.out"),
	regexp.MustCompile("ruby/2.5.0/extensions/x86_64-linux/.*/mkmf.log"),
	regexp.MustCompile("ruby/2.5.0/gems/.*/Makefile"),
}

const runtimePath = ".devflow/runtime"

type AppInfo struct {
	Runtime string `json:"runtime"`
	PreZip  string `json:"pre_zip"`
	LibZip  string `json:"lib_zip"`
	AppZip  string `json:"app_zip"`

	BuildID     string `json:"build_id"`
	MetadataSum string `json:"metadata_sum"`
}

func (b *Builder) AppInfo() *AppInfo {
	return &AppInfo{
		Runtime: b.runtime,
		PreZip:  b.preZip,
		LibZip:  b.libZip,
		AppZip:  b.appZip,

		BuildID:     b.id,
		MetadataSum: b.sum,
	}
}

func (b *Builder) ExtractRuntime(ctx context.Context, src *component.Source) (string, error) {
	runtime := b.config.Runtime

	_, ok := supportedRuntimes[runtime]
	if !ok {
		return "", errors.Wrapf(ErrInvalidRuntime, "runtime: %s", runtime)
	}

	return runtime, nil
}

func (b *Builder) randomId() (string, error) {
	u, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

type MappedFile struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

func (b *Builder) HashData(parts ...[]byte) (string, error) {
	h := sha256.New()

	for _, part := range parts {
		h.Write(part)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (b *Builder) HashDir(dir string) (string, error) {
	h := sha256.New()

	filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		if filepath.Base(path) == "data" || filepath.Base(path) == "cache" {
			return filepath.SkipDir
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		defer f.Close()

		io.Copy(h, f)

		return nil
	})

	return hex.EncodeToString(h.Sum(nil)), nil
}

func HashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	h := sha256.New()

	io.Copy(h, f)

	return h.Sum(nil), nil
}

func (b *Builder) dockerClient(ctx context.Context) (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	cli.NegotiateAPIVersion(ctx)

	return cli, nil
}

func (b *Builder) HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := sha256.New()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (b *Builder) Build(
	ctx context.Context,
	L hclog.Logger,
	src *component.Source,
	dir *datadir.Component,
) (component.Artifact, error) {

	subDir, err := datadir.NewScopedDir(dir, "lambda-build")
	if err != nil {
		return nil, err
	}

	scratchDir := subDir.DataDir()

	err = b.BuildPreImage(ctx, L, src, scratchDir)
	if err != nil {
		return nil, err
	}

	cli, err := b.dockerClient(ctx)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("devflow-%s-%s", src.App, b.sum)

	L.Info("running container application creation", "name", name)
	cfg := container.Config{
		AttachStdout: true,
		AttachStderr: true,
		Image:        b.preRef,
	}

	cfg.Cmd = append(cfg.Cmd, "/builder")

	absPath, err := filepath.Abs(src.Path)
	if err != nil {
		return nil, err
	}

	hostCfg := container.HostConfig{
		Binds: []string{absPath + ":/input"},
	}
	networkCfg := network.NetworkingConfig{}

	body, err := cli.ContainerCreate(ctx, &cfg, &hostCfg, &networkCfg, name)
	if err != nil {
		return nil, err
	}

	defer cli.ContainerRemove(ctx, body.ID, types.ContainerRemoveOptions{Force: true})

	err = cli.ContainerStart(ctx, body.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	opts := types.ContainerAttachOptions{
		Logs:   true,
		Stream: true,
		Stdout: true,
		Stderr: true,
	}

	resp, err := cli.ContainerAttach(ctx, body.ID, opts)
	if err != nil {
		return nil, err
	}

	var lg LogPrinter
	lg.Prefix = "[app-builder] "

	go func() {
		L.Info("forwarding logs from prehook container")
		lg.Display(resp.Reader)
		L.Info("logs from prehook container done")
	}()

	c, errc := cli.ContainerWait(ctx, body.ID, container.WaitConditionNotRunning)

	select {
	case <-ctx.Done():
		cli.ContainerKill(ctx, body.ID, "SIGKILL")
	case serr := <-errc:
		L.Error("error waiting for container", "error", serr)
		err = serr
	case resp := <-c:
		L.Info("container finished", "code", resp.StatusCode)

		switch resp.StatusCode {
		case 0:
			// ok!
		default:
			// unexpected error
			return nil, ErrUnexpectedError
		}
	}

	rid, err := b.randomId()
	if err != nil {
		return nil, err
	}

	ref := "devflow.local/tmp:" + rid

	L.Info("container finished and temp committed", "ref", ref)
	_, err = cli.ContainerCommit(ctx, body.ID, types.ContainerCommitOptions{
		Reference: ref,
	})
	if err != nil {
		return nil, err
	}

	defer cli.ImageRemove(ctx, ref, types.ImageRemoveOptions{PruneChildren: true})

	L.Info("extracting application from container", "ref", body.ID)

	diff, err := cli.ContainerDiff(ctx, body.ID)
	if err != nil {
		return nil, err
	}

	var files []MappedFile

	for _, change := range diff {
		if change.Kind == 2 { // skip anything deleted
			continue
		}

		path := change.Path
		src := path

		if path[0] == '/' {
			path = path[1:]
		}

		include := true
		for _, prefix := range ignorePrefixes {
			if strings.HasPrefix(path, prefix) {
				include = false
				break
			}
		}

		if include {
			for _, re := range pruneRegexp {
				if re.MatchString(path) {
					include = false
					break
				}
			}
		}

		if include {
			dest := path

			for _, shift := range shiftPrefix {
				if strings.HasPrefix(dest, shift.prefix) {
					dest = shift.shift + dest[len(shift.prefix):]
					break
				}
			}

			files = append(files, MappedFile{Source: src, Dest: dest})
		}
	}

	id, err := b.randomId()
	if err != nil {
		return nil, ErrInvalidRuntime
	}

	b.id = id

	zipFiles, err := b.extractFileDiff(ctx, L, src, scratchDir, id, id+"-lib", cli, ref, files)
	if err != nil {
		return nil, err
	}

	b.appZip = zipFiles[0]

	if len(zipFiles) > 1 {
		b.libZip = zipFiles[1]
	}

	appHash, err := b.HashFile(zipFiles[0])
	if err != nil {
		return nil, err
	}

	layerHash, _ := b.HashFile(zipFiles[1])

	L.Info("extracted app", "id", id, "app-hash", appHash[:10], "layer-hash", layerHash[:10])

	return b.AppInfo(), nil
}

func (b *Builder) BuildPreImage(ctx context.Context, L hclog.Logger, src *component.Source, scratchDir string) error {
	runtime, err := b.ExtractRuntime(ctx, src)
	if err != nil {
		return err
	}

	b.runtime = runtime

	L.Info("detected runtime", "runtime", runtime)

	cli, err := b.dockerClient(ctx)
	if err != nil {
		return err
	}

	sum, err := b.HashData([]byte(runtime), []byte(b.config.Setup))
	if err != nil {
		return err
	}

	b.sum = sum

	L.Info("devflow metadata hash", "hash", sum)

	diffFile := filepath.Join(scratchDir, "diff-"+sum+".json")

	var (
		changes   []string
		validDiff bool
		validImg  bool
	)

	df, err := os.Open(diffFile)
	if err != nil {
		L.Info("missing container diff listing for prepackage")
	} else {
		defer df.Close()

		err = json.NewDecoder(df).Decode(&changes)
		if err != nil {
			L.Info("invalid encoding of container diff", "error", err)
		} else {
			validDiff = true
		}
	}

	df.Close()

	ref := fmt.Sprintf("devflow.local/%s:pre-%s", src.App, sum)

	b.preRef = ref

	insp, _, err := cli.ImageInspectWithRaw(ctx, ref)
	if err == nil {
		validImg = true
	}

	if validDiff && validImg {
		L.Info("reusing prepackage image", "id", insp.ID, "created-at", insp.Created)
	} else {
		L.Info("creating prepackage layer")

		name := fmt.Sprintf("devflow-pre-%s-%s", src.App, sum)

		L.Info("running container for prehooks", "name", name, "src", src.Path)

		prePath := filepath.Join(scratchDir, "pre.sh")

		err = ioutil.WriteFile(prePath, []byte(b.config.Setup+"\n"), 0755)
		if err != nil {
			return err
		}

		absPath, err := filepath.Abs(scratchDir)
		if err != nil {
			return err
		}

		hostCfg := container.HostConfig{
			Binds: []string{absPath + ":/input"},
		}
		networkCfg := network.NetworkingConfig{}

		cfg := container.Config{
			AttachStdout: true,
			AttachStderr: true,
			Image:        "robloweco/lambda:" + runtime,
		}

		cfg.Cmd = append(cfg.Cmd, "/builder", "-pre", "/input/pre.sh")

		body, err := cli.ContainerCreate(ctx, &cfg, &hostCfg, &networkCfg, name)
		if err != nil {
			return err
		}

		defer cli.ContainerRemove(ctx, body.ID, types.ContainerRemoveOptions{Force: true})

		err = cli.ContainerStart(ctx, body.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}

		opts := types.ContainerAttachOptions{
			Logs:   true,
			Stream: true,
			Stdout: true,
			Stderr: true,
		}

		resp, err := cli.ContainerAttach(ctx, body.ID, opts)
		if err != nil {
			return err
		}

		var lg LogPrinter
		lg.Prefix = "[pre-builder] "

		go func() {
			L.Info("forwarding logs from prehook container")
			lg.Display(resp.Reader)
			L.Info("logs from prehook container done")
		}()

		c, errc := cli.ContainerWait(ctx, body.ID, container.WaitConditionNotRunning)

		var packagePre bool

		select {
		case <-ctx.Done():
			cli.ContainerKill(ctx, body.ID, "SIGKILL")
		case serr := <-errc:
			L.Error("error waiting for container", "error", serr)
			err = serr
		case resp := <-c:
			L.Info("container finished", "code", resp.StatusCode)

			switch resp.StatusCode {
			case 0:
				// ok!
				packagePre = true
			case 10:
				// No prehooks defined
			default:
				// unexpected error
				err = ErrUnexpectedError
			}
		}

		L.Info("prehooks container finished", "package-pre", packagePre)
		_, err = cli.ContainerCommit(ctx, body.ID, types.ContainerCommitOptions{
			Reference: ref,
		})
		if err != nil {
			return err
		}

		L.Info("prehook image created", "ref", ref)

		diff, err := cli.ContainerDiff(ctx, body.ID)
		if err != nil {
			return err
		}

		changes = nil

		for _, change := range diff {
			if change.Kind == 2 { // skip anything deleted
				continue
			}

			changes = append(changes, change.Path)
		}

		df, err = os.Create(diffFile)
		if err != nil {
			return err
		}

		defer df.Close()

		err = json.NewEncoder(df).Encode(changes)
		if err != nil {
			return err
		}

		df.Close()
	}

	var files []MappedFile

	for _, path := range changes {
		src := path

		if path[0] == '/' {
			path = path[1:]
		}

		include := true
		for _, prefix := range ignorePrefixes {
			if strings.HasPrefix(path, prefix) {
				include = false
				break
			}
		}

		if include {
			for _, re := range pruneRegexp {
				if re.MatchString(path) {
					include = false
					break
				}
			}
		}

		if include {
			dest := path

			for _, shift := range shiftPrefix {
				if strings.HasPrefix(dest, shift.prefix) {
					dest = shift.shift + dest[len(shift.prefix):]
					break
				}
			}

			files = append(files, MappedFile{Source: src, Dest: dest})
		}
	}

	zipFiles, err := b.extractFileDiff(ctx, L, src, scratchDir, sum+"-pre", "", cli, ref, files)
	if err != nil {
		return err
	}

	b.preZip = zipFiles[0]

	appHash, err := b.HashFile(zipFiles[0])
	if err != nil {
		return err
	}

	L.Info("extracted pre-layer", "id", sum, "hash", appHash[:10])

	return nil
}

func (b *Builder) extractFileDiff(ctx context.Context, L hclog.Logger, src *component.Source, scratchDir, id, secondaryId string, cli *client.Client, image string, files []MappedFile) ([]string, error) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Source < files[j].Source
	})

	cfg := container.Config{
		AttachStdout: true,
		AttachStderr: true,
		Image:        image,
	}

	localPath := filepath.Join(scratchDir, id)

	w, err := os.Create(localPath)
	if err != nil {
		return nil, err
	}
	defer w.Close()

	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		return nil, err
	}

	w.Close()

	cfg.Cmd = append(cfg.Cmd, "/builder", "-extract", "/output/"+id)
	if secondaryId != "" {
		cfg.Cmd = append(cfg.Cmd, "-layer", "/output/"+secondaryId)
	}

	scratch, err := filepath.Abs(scratchDir)
	if err != nil {
		return nil, err
	}

	hostCfg := container.HostConfig{
		Binds: []string{scratch + ":/output"},
	}
	networkCfg := network.NetworkingConfig{}

	name := fmt.Sprintf("devflow-%s-%s-extract", src.App, id)

	body, err := cli.ContainerCreate(ctx, &cfg, &hostCfg, &networkCfg, name)
	if err != nil {
		return nil, err
	}

	err = cli.ContainerStart(ctx, body.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	defer cli.ContainerRemove(ctx, body.ID, types.ContainerRemoveOptions{Force: true})

	opts := types.ContainerAttachOptions{
		Logs:   true,
		Stream: true,
		Stdout: true,
		Stderr: true,
	}

	resp, err := cli.ContainerAttach(ctx, body.ID, opts)
	if err != nil {
		return nil, err
	}

	logsDone := make(chan struct{})

	var lg LogPrinter
	lg.Prefix = "[extract-files] "

	go func() {
		L.Info("forwarding logs from prehook container")
		lg.Display(resp.Reader)
		logsDone <- struct{}{}
		L.Info("logs from prehook container done")
	}()

	c, errc := cli.ContainerWait(ctx, body.ID, container.WaitConditionNotRunning)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case serr := <-errc:
		L.Error("error waiting for container", "error", serr)
		return nil, serr
	case resp := <-c:
		L.Info("container finished", "code", resp.StatusCode)

		switch resp.StatusCode {
		case 0:
			// ok!
		default:
			// unexpected error
			return nil, ErrUnexpectedError
		}
	}

	var zipFiles []string

	primary := filepath.Join(scratchDir, id+".zip")
	if _, err := os.Stat(primary); err != nil {
		return nil, ErrUnexpectedError
	}

	L.Info("gathering zip files", "id", id, "secondary-id", secondaryId)

	zipFiles = append(zipFiles, primary)

	if secondaryId != "" {
		secondary := filepath.Join(scratchDir, secondaryId+".zip")
		if _, err := os.Stat(secondary); err == nil {
			zipFiles = append(zipFiles, secondary)
		}
	}

	return zipFiles, nil
}
