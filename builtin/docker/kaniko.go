// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package docker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/docker/distribution/reference"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/assets"
	"github.com/hashicorp/waypoint/internal/pkg/epinject/ociregistry"
	"github.com/pkg/errors"
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
		Location: &Image_Registry{Registry: &Image_RegistryLocation{}},
	}

	var os ociregistry.Server
	os.DisableEntrypoint = b.config.DisableCEB
	os.Logger = log

	if ai.Auth != nil {
		switch sv := ai.Auth.(type) {
		case *AccessInfo_Encoded:
			user, pass, err := CredentialsFromConfig(sv.Encoded)
			if err != nil {
				return nil, err
			}
			os.AuthConfig.Username = user
			os.AuthConfig.Password = pass
		case *AccessInfo_Header:
			os.AuthConfig.Auth = sv.Header
		case *AccessInfo_UserPass_:
			os.AuthConfig.Username = sv.UserPass.Username
			os.AuthConfig.Password = sv.UserPass.Password
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
		os.Upstream = "http://" + host
	} else {
		os.Upstream = "https://" + host
	}

	refPath := reference.Path(ref)

	err = os.Negotiate(ref.Name())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to negotiate with upstream")
	}

	if !b.config.DisableCEB {
		// For Kaniko we can use our runtime arch because the image we build
		// always matches the architecture of our Kaniko environment.
		assetName, ok := assets.CEBArch[runtime.GOARCH]
		if !ok {
			return nil, status.Errorf(codes.FailedPrecondition,
				"automatic entrypoint injection not supported for architecture: %s", runtime.GOARCH)
		}

		data, err := assets.Asset(assetName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to restore custom entry point binary: %s", err)
		}

		step.Done()
		step = sg.Add("Testing registry and uploading entrypoint layer")

		err = os.SetupEntrypointLayer(refPath, data)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error setting up entrypoint layer to host %q, err: %s", os.Upstream, err)
		}
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

	if b.config.Target != "" {
		args = append(args, "--target", b.config.Target)
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

	// Command output should go to the step
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
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
