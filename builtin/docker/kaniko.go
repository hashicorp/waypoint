package docker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"

	"github.com/docker/distribution/reference"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/assets"
	"github.com/hashicorp/waypoint/internal/ociregistry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (b *Builder) buildWithKaniko(
	ctx context.Context,
	ui terminal.UI,
	sg terminal.StepGroup,
	log hclog.Logger,
	dockerfilePath string,
	contextDir string,
	buildArgs map[string]*string,
	ai *AccessInfo,
) (*Image, error) {
	step := sg.Add("Building Docker image with kaniko...")
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	target := &Image{
		Image:    ai.Image,
		Tag:      ai.Tag,
		Location: &Image_Registry{Registry: &empty.Empty{}},
	}

	var auth string

	if ai.Auth != nil {
		switch sv := ai.Auth.(type) {
		case *AccessInfo_Encoded:
			user, pass, err := CredentialsFromConfig(sv.Encoded)
			if err != nil {
				return nil, err
			}
			auth = ociregistry.BasicAuth(user, pass)
		case *AccessInfo_Header:
			auth = sv.Header
		}
	}

	// Determine the host that we're setting auth for. We have to parse the
	// image for this cause it may not contain a host. Luckily Docker has
	// libs to normalize this all for us.
	log.Trace("determining host for auth configuration", "image", target.Name())
	ref, err := reference.ParseNormalizedNamed(target.Name())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}
	host := reference.Domain(ref)
	log.Trace("auth host", "host", host)

	var os ociregistry.Server
	os.DisableEntrypoint = b.config.DisableCEB
	os.Auth = auth
	os.Logger = log
	os.Upstream = "http://" + host

	data, err := assets.Asset("ceb/ceb")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to restore custom entry point binary: %s", err)
	}

	refPath := reference.Path(ref)

	step.Done()
	step = sg.Add("Testing registry and uploading entrypoint layer")

	err = os.SetupEntrypointLayer(refPath, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error setting up entrypoint layer: %s", err)
	}

	li, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}

	defer li.Close()
	go http.Serve(li, &os)

	port := li.Addr().(*net.TCPAddr).Port

	localRef := fmt.Sprintf("localhost:%d/%s:%s", port, refPath, ai.Tag)

	// Start constructing our arg string for img
	args := []string{
		"/kaniko/executor",
		"--context", "dir://" + contextDir,
		"-f", dockerfilePath,
		"-d", localRef,
	}

	// If we have build args we append each
	for k, v := range buildArgs {
		// v should always not be nil but guard just in case to avoid a panic
		if v != nil {
			args = append(args, "--build-arg", k+"="+*v)
		}
	}

	log.Debug("executing kaniko", "args", args)

	step.Done()
	step = sg.Add("Executing kaniko...")

	// NOTE(mitchellh): we can probably use the img Go pkg directly one day.
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// Command output should go to the step
	cmd.Stdout = step.TermOutput()
	cmd.Stderr = cmd.Stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	step.Done()

	step = sg.Add("Image pushed to '%s:%s'", ai.Image, ai.Tag)
	step.Done()

	return target, nil
}
