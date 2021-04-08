package epinject

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hashicorp/go-hclog"
	"github.com/oklog/ulid"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// AlterEntrypointImg has the same signature as AlterEntrypoint but uses
// "img" under the hood to perform the entrypoint modification.
func AlterEntrypointImg(
	ctx context.Context,
	image string,
	cb func(cur []string) (*NewEntrypoint, error),
) (string, error) {
	L := hclog.FromContext(ctx).With("image", image)
	L.Debug("altering entrypoint of docker image using img")

	// Create a temporary directory. We do this in case img creates state
	// (it does not in the current directory but in case it ever does)
	td, err := ioutil.TempDir("", "wp-epinject")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(td)

	// Determine the existing entrypoint
	L.Debug("executing img inspect to determine existing entrypoint")
	var buf bytes.Buffer
	cmd := exec.CommandContext(ctx,
		"img",
		"inspect",
		image,
	)
	cmd.Dir = td
	cmd.Stdout = &buf
	cmd.Stderr = cmd.Stdout
	if err := cmd.Run(); err != nil {
		L.Warn("img inspect failed", "err", err, "output", buf.String())
		return "", err
	}

	// Parse the image
	var imageSpec ocispec.Image
	if err := json.Unmarshal(buf.Bytes(), &imageSpec); err != nil {
		return "", err
	}

	L.Debug("extracted existing entrypoint", "entrypoint", imageSpec.Config.Entrypoint)

	// Determine the new entrypoint configuration based on the existing
	newEp, err := cb(imageSpec.Config.Entrypoint)
	if err != nil {
		return "", err
	}
	if newEp.Entrypoint == nil {
		newEp.Entrypoint = imageSpec.Config.Entrypoint
	}
	if newEp.NewImage == "" {
		newEp.NewImage = image
	}

	// Create a random name for our new container image
	u, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return "", err
	}
	name := strings.ToLower("epinject-" + u.String())

	// Build our template data. The entrypoint generates the actual
	// entrypoint string to use in the Dockerfile.
	var tplData tplData
	tplData.Base = image
	if len(newEp.Entrypoint) > 0 {
		v, err := json.Marshal(newEp.Entrypoint)
		if err != nil {
			return "", err
		}

		tplData.Entrypoint = string(v)
	}
	if len(imageSpec.Config.Cmd) > 0 {
		v, err := json.Marshal(imageSpec.Config.Cmd)
		if err != nil {
			return "", err
		}

		tplData.Cmd = string(v)
	}

	// For every file, we copy it into our temporary directory so it can be copied.
	L.Debug("copying files for injection", "n", len(newEp.InjectFiles))
	idx := 0
	for containerPath, finfo := range newEp.InjectFiles {
		// Local path for our file copy
		localPath := filepath.Join(td, fmt.Sprintf("%d", idx))
		idx++
		L.Trace("staging file for copy", "from", localPath, "to", containerPath)

		f, err := os.Create(localPath)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(f, finfo.Reader); err != nil {
			f.Close()
			return "", err
		}
		if err := f.Sync(); err != nil {
			f.Close()
			return "", err
		}

		// We have to chmod the file because Buildkit inherits the file
		// permissions of the source file. This doesn't work on Windows but
		// we don't support img on Windows so that's okay.
		if err := f.Chmod(finfo.Info.Mode()); err != nil {
			f.Close()
			return "", err
		}

		f.Close()

		tplData.Copy = append(tplData.Copy, tplDataFile{
			From: "./" + filepath.Base(localPath),
			To:   containerPath,
		})
	}

	// Render our Dockerfile into our directory
	dockerfilePath := filepath.Join(td, "Dockerfile")
	f, err := os.Create(dockerfilePath)
	if err != nil {
		return "", err
	}
	tpl := template.Must(template.New("dockerfile").Parse(dockerfileTemplate))
	err = tpl.Execute(f, &tplData)
	f.Close()
	if err != nil {
		return "", err
	}

	// Execute the build
	cmd = exec.CommandContext(ctx,
		"img",
		"build",
		"-f", dockerfilePath,
		"-t", name,
		".",
	)
	cmd.Dir = td
	cmd.Stdout = &buf
	cmd.Stderr = cmd.Stdout
	buf.Reset()
	L.Debug("executing img build for injection", "args", cmd.Args)
	if err := cmd.Run(); err != nil {
		L.Warn("failed to inject", "err", err, "output", buf.String())
		return "", err
	}

	// Retag to the final name
	cmd = exec.CommandContext(ctx,
		"img",
		"tag",
		name,
		newEp.NewImage,
	)
	cmd.Dir = td
	cmd.Stdout = &buf
	cmd.Stderr = cmd.Stdout
	buf.Reset()
	L.Debug("executing img tag to rename", "args", cmd.Args)
	if err := cmd.Run(); err != nil {
		L.Warn("failed to tag", "err", err, "output", buf.String())
		return "", err
	}

	// Remove the temporary image
	cmd = exec.CommandContext(ctx,
		"img",
		"rm",
		name,
	)
	cmd.Dir = td
	cmd.Stdout = &buf
	cmd.Stderr = cmd.Stdout
	buf.Reset()
	L.Debug("removing build image", "args", cmd.Args)
	if err := cmd.Run(); err != nil {
		L.Warn("failed to remove build image", "err", err, "output", buf.String())
		return "", err
	}

	// Create the image
	return newEp.NewImage, nil
}

const dockerfileTemplate = `
FROM {{.Base}}

{{range $file := .Copy}}
COPY {{.From}} {{.To}}
{{end}}

{{if .Cmd}}
CMD {{.Cmd}}
{{end}}

{{if .Entrypoint}}
ENTRYPOINT {{.Entrypoint}}
{{end}}
`

type tplData struct {
	Base       string // Base image
	Copy       []tplDataFile
	Cmd        string
	Entrypoint string
}

type tplDataFile struct {
	From string
	To   string
}
