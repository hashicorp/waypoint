package lambda

import (
	"archive/tar"
	"archive/zip"
	"context"
	"crypto/rand"
	"crypto/sha1"
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
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/datadir"
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

func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
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

func (b *Builder) AppInfo() *AppInfo {
	return &AppInfo{
		Runtime: b.runtime,
		PreZip:  b.preZip,
		LibZip:  b.libZip,
		AppZip:  b.appZip,

		BuildId:     b.id,
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

type builderInput struct {
	Schema string `json:"jsonschema"`
	Id     int    `json:"id"`
	Method string `json:"method"`
	Params struct {
		Version    string `json:"__protocol_version"`
		Capability struct {
			Language             string `json:"language"`
			DependencyManager    string `json:"dependency_manager"`
			ApplicationFramework string `json:"application_framework"`
		} `json:"capability"`
		SourceDir     string   `json:"source_dir"`
		ArtifactsDir  string   `json:"artifacts_dir"`
		ScratchDir    string   `json:"scratch_dir"`
		ManifestPath  string   `json:"manifest_path"`
		Runtime       string   `json:"runtime"`
		Optimizations string   `json:"optimizations"`
		Options       string   `json:"options"`
		SearchPaths   []string `json:"executable_search_paths"`
		Mode          string   `json:"mode"`
	} `json:"params"`
}

func (b *Builder) Build(
	ctx context.Context,
	L hclog.Logger,
	src *component.Source,
	dir *datadir.Component,
) (*AppInfo, error) {

	rtc, err := FindRuntimeConfig(b.config.Runtime, src.Path)
	if err != nil {
		return nil, err
	}

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
		AttachStdin:  true,
		OpenStdin:    true,
		StdinOnce:    true,
		Image:        b.preRef,
	}

	cfg.Cmd = append(cfg.Cmd, "/bin/sh", "-c", "GLOBIGNORE=*/.devflow; mkdir -p /tmp/src /tmp/scratch; for i in /input/*; do ln -s $i /tmp/src/; done; lambda-builders")

	absPath, err := filepath.Abs(src.Path)
	if err != nil {
		return nil, err
	}

	outputPath, err := filepath.Abs(filepath.Join(scratchDir, "output"))
	if err != nil {
		return nil, err
	}

	hostCfg := container.HostConfig{
		Binds: []string{absPath + ":/input", outputPath + ":/output"},
	}
	networkCfg := network.NetworkingConfig{}

	body, err := cli.ContainerCreate(ctx, &cfg, &hostCfg, &networkCfg, name)
	if err != nil {
		return nil, err
	}

	defer cli.ContainerRemove(ctx, body.ID, types.ContainerRemoveOptions{Force: true})

	opts := types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
		Stdin:  true,
	}

	resp, err := cli.ContainerAttach(ctx, body.ID, opts)
	if err != nil {
		return nil, err
	}

	var lg LogPrinter
	lg.Prefix = "[app-builder] "

	go func() {
		L.Info("forwarding logs from build container")
		lg.Display(resp.Reader)
		L.Info("logs from container done")
	}()

	// Ok, rotate the .devflow directory to avoid pulling it into the build.
	err = cli.ContainerStart(ctx, body.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	var bc builderInput
	bc.Schema = "2.0"
	bc.Id = 1
	bc.Method = "LambdaBuilder.build"
	bc.Params.Version = "0.3"
	bc.Params.Capability.Language = rtc.Language
	bc.Params.Capability.DependencyManager = rtc.DepManager
	bc.Params.Capability.ApplicationFramework = rtc.AppFramework
	bc.Params.SourceDir = "/tmp/src"
	bc.Params.ScratchDir = "/tmp/scratch"
	bc.Params.ArtifactsDir = "/output"
	bc.Params.ManifestPath = "/input/Gemfile"
	bc.Params.Runtime = b.config.Runtime
	bc.Params.SearchPaths = rtc.SearchPaths

	err = json.NewEncoder(resp.Conn).Encode(&bc)
	if err != nil {
		return nil, err
	}

	resp.CloseWrite()

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

	id, err := b.randomId()
	if err != nil {
		return nil, ErrInvalidRuntime
	}

	b.id = id

	b.appZip = filepath.Join(scratchDir, id+".zip")

	af, err := os.Create(b.appZip)
	if err != nil {
		return nil, err
	}

	defer af.Close()

	az := zip.NewWriter(af)

	defer az.Close()

	b.libZip = filepath.Join(scratchDir, id+"-lib.zip")

	lf, err := os.Create(b.libZip)
	if err != nil {
		return nil, err
	}

	defer lf.Close()

	lz := zip.NewWriter(lf)

	defer lz.Close()

	filepath.Walk(outputPath, func(entryPath string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		if outputPath == entryPath {
			return nil
		}

		path := entryPath[len(outputPath):]
		if path == "" {
			return nil
		}

		if path[0] == '/' {
			path = path[1:]
		}

		for _, re := range pruneRegexp {
			if re.MatchString(path) {
				return nil
			}
		}

		hdr, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		hdr.Modified = time.Time{}
		hdr.ModifiedDate = 0
		hdr.ModifiedTime = 0

		hdr.Name = path

		var ew io.Writer

		for _, prefix := range layerPrefixes {
			if strings.HasPrefix(path, prefix) {

				ew, err = lz.CreateHeader(hdr)
				if err != nil {
					return err
				}

				break
			}
		}

		if ew == nil {
			ew, err = az.CreateHeader(hdr)
			if err != nil {
				return err
			}
		}

		f, err := os.Open(entryPath)
		if err != nil {
			return err
		}

		defer f.Close()

		io.Copy(ew, f)

		return nil
	})

	appHash, err := b.HashFile(b.appZip)
	if err != nil {
		return nil, err
	}

	layerHash, _ := b.HashFile(b.libZip)

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
			Image:        "lambci/lambda:build-" + runtime,
		}

		cfg.Cmd = append(cfg.Cmd, "/bin/bash", "-c", "cd /tmp && cp /input/pre.sh . && bash ./pre.sh")

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

	sourceMap := map[string]string{}

	for _, file := range files {
		sourceMap[file.Source] = file.Dest
		fmt.Fprintln(w, file.Source)
	}

	// err = json.NewEncoder(w).Encode(files)
	// if err != nil {
	// return nil, err
	// }

	w.Close()

	cfg.Cmd = append(cfg.Cmd, "/bin/sh", "-c", fmt.Sprintf(
		`xargs -rd '\n' sh -c 'exec find "$@" -prune ! -type d' sh < /output/%s| tar -T - -cvpPnf /output/main.tar`, id))

	// cfg.Cmd = append(cfg.Cmd, "zip", "-@", "/output/files.zip")
	// if secondaryId != "" {
	// cfg.Cmd = append(cfg.Cmd, "-layer", "/output/"+secondaryId)
	// }

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

	L.Info("extracting files into lambda zips")

	var zipFiles []string

	primary := filepath.Join(scratchDir, id+".zip")

	zipFiles = append(zipFiles, primary)

	pf, err := os.Create(primary)
	if err != nil {
		return nil, err
	}

	defer pf.Close()

	pfz := zip.NewWriter(pf)

	defer pfz.Close()

	var (
		sf  *os.File
		sfz *zip.Writer
	)

	if secondaryId != "" {
		secondary := filepath.Join(scratchDir, secondaryId+".zip")
		zipFiles = append(zipFiles, secondary)

		sf, err = os.Create(secondary)
		if err != nil {
			return nil, err
		}

		defer sf.Close()

		sfz = zip.NewWriter(sf)

		defer sfz.Close()
	}

	tf, err := os.Open(filepath.Join(scratchDir, "main.tar"))
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(tf)

	h := sha1.New()

	var (
		layerFiles int
	)

	for {
		thdr, err := tr.Next()
		if err != nil {
			break
		}

		hdr, err := zip.FileInfoHeader(thdr.FileInfo())
		if err != nil {
			return nil, err
		}

		path := sourceMap[thdr.Name]

		if sfz != nil && strings.HasPrefix(path, "_layer/") {
			hdr.Name = path[len("_layer/"):]

			body, err := sfz.CreateHeader(hdr)
			if err != nil {
				return nil, err
			}

			io.Copy(io.MultiWriter(body, h), tr)

			layerFiles++

		} else {
			hdr.Name = path

			body, err := pfz.CreateHeader(hdr)
			if err != nil {
				return nil, err
			}

			io.Copy(io.MultiWriter(body, h), tr)
		}
	}

	L.Info("extracted files", "cas", hex.EncodeToString(h.Sum(nil)))

	return zipFiles, nil
}

var (
	_ component.Builder      = (*Builder)(nil)
	_ component.Configurable = (*Builder)(nil)
)
