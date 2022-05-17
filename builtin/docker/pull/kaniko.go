package dockerpull

import (
	"context"
	"fmt"
	"github.com/docker/distribution/reference"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	wpdocker "github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/internal/assets"
	"github.com/hashicorp/waypoint/internal/pkg/epinject/ociregistry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func (b *Builder) pullWithKaniko(
	ctx context.Context,
	ui terminal.UI,
	sg terminal.StepGroup,
	log hclog.Logger,
	ai *wpdocker.AccessInfo,
) (*wpdocker.Image, error) {
	step := sg.Add("Pulling Docker image with Kaniko...")
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	target := &wpdocker.Image{
		Image:    b.config.Image,
		Tag:      b.config.Tag,
		Location: &wpdocker.Image_Docker{Docker: &empty.Empty{}},
	}

	var oci ociregistry.Server
	oci.DisableEntrypoint = b.config.DisableCEB
	oci.Logger = log

	if ai.Auth != nil {
		switch sv := ai.Auth.(type) {
		case *wpdocker.AccessInfo_Encoded:
			user, pass, err := wpdocker.CredentialsFromConfig(sv.Encoded)
			if err != nil {
				return nil, err
			}
			oci.AuthConfig.Username = user
			oci.AuthConfig.Password = pass
		case *wpdocker.AccessInfo_Header:
			oci.AuthConfig.Auth = sv.Header
		case *wpdocker.AccessInfo_UserPass_:
			oci.AuthConfig.Username = sv.UserPass.Username
			oci.AuthConfig.Password = sv.UserPass.Password
		default:
			return nil, status.Error(codes.Unauthenticated, "Unexpected auth type.")
		}
	}

	// Determine the host that we're setting auth for. We have to parse the
	// image for this cause it may not contain a host. Luckily Docker has
	// libs to normalize this all for us.
	log.Trace("determining host for auth configuration", "image", target.Name())
	ref, err := reference.ParseNormalizedNamed(target.Image)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}

	host := reference.Domain(ref)
	if host == "docker.io" {
		// The normalized name parse above will turn short names like "foo/bar"
		// into "docker.io/foo/bar" but the actual registry host for these
		// is "index.docker.io".
		host = "index.docker.io"
	}
	log.Trace("auth host", "host", host)

	if ai.Insecure {
		oci.Upstream = "http://" + host
	} else {
		oci.Upstream = "https://" + host
	}

	refPath := reference.Path(ref)

	if !b.config.DisableCEB {
		step.Update("Injecting entrypoint...")
		// For Kaniko we can use our runtime arch because the image we build
		// always matches the architecture of our Kaniko environment.
		assetName, ok := assets.CEBArch[runtime.GOARCH]
		if !ok {
			return nil, status.Errorf(codes.FailedPrecondition,
				"automatic entrypoint injection not supported for architecture: %s", runtime.GOARCH)
		}

		data, err := assets.Asset(assetName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to restore custom entrypoint binary: %s", err)
		}

		step = sg.Add("Testing registry and uploading entrypoint layer")

		err = oci.SetupEntrypointLayer(refPath, data)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error setting up entrypoint layer: %s", err)
		}
	}

	// Setting up local registry to which Kaniko will push
	li, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}

	defer li.Close()
	go http.Serve(li, &oci)

	port := li.Addr().(*net.TCPAddr).Port

	localRef := fmt.Sprintf("localhost:%d/%s:%s", port, refPath, ai.Tag)

	dockerfileBS := []byte(fmt.Sprintf("FROM %s:%s\n", target.Image, target.Tag))
	err = os.WriteFile("Dockerfile", dockerfileBS, 0644)
	if err != nil {
		return nil, err
	}
	// Start constructing our arg string for img
	args := []string{
		"/kaniko/executor",
		"-f", filepath.Dir("Dockerfile"),
		"-d", localRef,
	}

	log.Debug("executing kaniko", "args", args)
	step.Update("Executing Kaniko...")

	// Command output should go to the step
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = step.TermOutput()
	cmd.Stderr = cmd.Stdout

	// check for error from executor run
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	step.Done()

	step = sg.Add("Image pull completed.")
	step.Done()

	return target, nil
}
