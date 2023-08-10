// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/docker/cli/cli/config"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/registry"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wpdockerclient "github.com/hashicorp/waypoint/builtin/docker/client"
)

func (r *Registry) pushWithDocker(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	source *Image,
	target *Image,
	authConfig *Auth,
) error {
	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "unable to create output for logs:%s", err)
	}

	sg := ui.StepGroup()
	defer sg.Wait()
	step := sg.Add("Initializing Docker client...")
	defer func() { step.Abort() }()

	cli, err := wpdockerclient.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "unable to create Docker client:%s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	step.Update("Tagging Docker image: %s => %s:%s", source.Name(), r.config.Image, r.config.Tag)
	err = cli.ImageTag(ctx, source.Name(), target.Name())
	if err != nil {
		return status.Errorf(codes.Internal, "unable to tag image:%s", err)
	}

	step.Done()

	if r.config.Local {
		return nil
	}

	ref, err := reference.ParseNormalizedNamed(target.Name())
	if err != nil {
		return status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}

	var encodedAuth = ""

	if r.config.EncodedAuth != "" {
		encodedAuth = r.config.EncodedAuth
	} else if encodedAuth == "" && r.config.Password != "" {
		// If there was no explicit encoded auth but there is a password, make the username+password
		// into an encoded auth string.
		var authConfig types.AuthConfig

		authConfig.Username = r.config.Username
		authConfig.Password = r.config.Password

		buf, err := json.Marshal(authConfig)
		if err != nil {
			return status.Errorf(codes.Internal, "unable to generate authentication info for registry: %s", err)
		}
		encodedAuth = base64.URLEncoding.EncodeToString(buf)
	} else if r.config.Auth == nil && r.config.EncodedAuth == "" {
		// Resolve the Repository name from fqn to RepositoryInfo
		repoInfo, err := registry.ParseRepositoryInfo(ref)
		if err != nil {
			return status.Errorf(codes.Internal, "unable to parse repository info from image name: %s", err)
		}

		var server string

		if repoInfo.Index.Official {
			info, err := cli.Info(ctx)
			if err != nil || info.IndexServerAddress == "" {
				server = registry.IndexServer
			} else {
				server = info.IndexServerAddress
			}
		} else {
			server = repoInfo.Index.Name
		}

		var errBuf bytes.Buffer
		cf := config.LoadDefaultConfigFile(&errBuf)
		if errBuf.Len() > 0 {
			// NOTE(mitchellh): I don't know why we ignore this, but we always have.
			log.Warn("error loading Docker config file", "err", err)
		}

		authConfig, _ := cf.GetAuthConfig(server)
		buf, err := json.Marshal(authConfig)
		if err != nil {
			return status.Errorf(codes.Internal, "unable to generate authentication info for registry: %s", err)
		}
		encodedAuth = base64.URLEncoding.EncodeToString(buf)
	} else if (*r.config.Auth != Auth{}) {
		if authConfig.Hostname != "" {
			return status.Errorf(codes.InvalidArgument, "hostname not supported for registry")
		}
		authBytes, err := json.Marshal(types.AuthConfig{
			Username:      authConfig.Username,
			Password:      authConfig.Password,
			Email:         authConfig.Email,
			Auth:          authConfig.Auth,
			ServerAddress: authConfig.ServerAddress,
			IdentityToken: authConfig.IdentityToken,
			RegistryToken: authConfig.RegistryToken,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal auth info to json: %s", err)
		}
		encodedAuth = base64.URLEncoding.EncodeToString(authBytes)
	}

	step = sg.Add("Pushing Docker image...")

	options := types.ImagePushOptions{
		RegistryAuth: encodedAuth,
	}

	responseBody, err := cli.ImagePush(ctx, reference.FamiliarString(ref), options)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to push image to registry: %s", err)
	}

	defer responseBody.Close()

	var termFd uintptr
	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
	}

	err = jsonmessage.DisplayJSONMessagesStream(responseBody, step.TermOutput(), termFd, true, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to stream Docker logs to terminal: %s", err)
	}

	step.Done()
	return nil
}
