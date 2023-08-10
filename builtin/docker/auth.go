// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package docker

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/types"
	"github.com/docker/distribution/reference"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CredentialsFromConfig returns the username and password present in the encoded
// auth string. This encoded auth string is one that users can pass as authentication
// information to registry.
func CredentialsFromConfig(encodedAuth string) (string, string, error) {
	// Create a reader that base64 decodes our encoded auth and then
	// JSON decodes that.
	var authCfg types.AuthConfig

	dec := json.NewDecoder(
		base64.NewDecoder(base64.URLEncoding, strings.NewReader(encodedAuth)),
	)

	if err := dec.Decode(&authCfg); err != nil {
		return "", "", status.Errorf(codes.FailedPrecondition,
			"Failed to decode encoded_auth: %s", err)
	}

	return authCfg.Username, authCfg.Password, nil
}

// TempDockerConfig creates a new Docker configuration with the
// configured auth in it. It saves this Docker config to a temporary path
// and returns the path to that Docker file.
//
// We have to do this because `img` doesn't support setting auth for
// a single operation. Therefore, we must set auth in the Docker config,
// but we don't want to pollute any concurrent runs or the main file. So
// we create a copy.
//
// This can return ("", nil) if there is no custom Docker config necessary.
//
// Callers should defer file deletion for this temporary file.
func TempDockerConfig(
	log hclog.Logger,
	target *Image,
	encodedAuth string,
) (string, error) {
	if encodedAuth == "" {
		return "", nil
	}

	// Create a reader that base64 decodes our encoded auth and then
	// JSON decodes that.
	var authCfg types.AuthConfig
	var rdr io.Reader = strings.NewReader(encodedAuth)
	rdr = base64.NewDecoder(base64.URLEncoding, rdr)
	dec := json.NewDecoder(rdr)
	if err := dec.Decode(&authCfg); err != nil {
		return "", status.Errorf(codes.FailedPrecondition,
			"Failed to decode encoded_auth: %s", err)
	}

	// Determine the host that we're setting auth for. We have to parse the
	// image for this cause it may not contain a host. Luckily Docker has
	// libs to normalize this all for us.
	log.Trace("determining host for auth configuration", "image", target.Name())
	ref, err := reference.ParseNormalizedNamed(target.Name())
	if err != nil {
		return "", status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}
	host := reference.Domain(ref)
	log.Trace("auth host", "host", host)

	// Parse our old Docker config and add the auth.
	log.Trace("loading Docker configuration")
	file, err := config.Load(config.Dir())
	if err != nil {
		return "", err
	}

	if file.AuthConfigs == nil {
		file.AuthConfigs = map[string]types.AuthConfig{}
	}
	file.AuthConfigs[host] = authCfg

	// Create a temporary directory for our config
	td, err := ioutil.TempDir("", "wp-docker-config")
	if err != nil {
		return "", status.Errorf(codes.Internal,
			"Failed to create temporary directory for Docker config: %s", err)
	}

	// Create a temporary file and write our Docker config to it
	f, err := os.Create(filepath.Join(td, "config.json"))
	if err != nil {
		return "", status.Errorf(codes.Internal,
			"Failed to create temporary file for Docker config: %s", err)
	}
	defer f.Close()
	if err := file.SaveToWriter(f); err != nil {
		return "", status.Errorf(codes.Internal,
			"Failed to create temporary file for Docker config: %s", err)
	}

	log.Info("temporary Docker config created for auth",
		"auth_host", host,
		"path", td,
	)

	return td, nil
}

// TempDockerConfigWithPassword creates a new Docker configuration with the
// configured auth in it. It saves this Docker config to a temporary path
// and returns the path to that Docker file.
//
// We have to do this because `img` doesn't support setting auth for
// a single operation. Therefore, we must set auth in the Docker config,
// but we don't want to pollute any concurrent runs or the main file. So
// we create a copy.
//
// This can return ("", nil) if there is no custom Docker config necessary.
//
// Callers should defer file deletion for this temporary file.
func TempDockerConfigWithPassword(
	log hclog.Logger,
	target *Image,
	username string,
	password string,
) (string, error) {
	if password == "" {
		return "", nil
	}

	// Create a reader that base64 decodes our encoded auth and then
	// JSON decodes that.
	var authCfg types.AuthConfig

	authCfg.Username = username
	authCfg.Password = password

	// Determine the host that we're setting auth for. We have to parse the
	// image for this cause it may not contain a host. Luckily Docker has
	// libs to normalize this all for us.
	log.Trace("determining host for auth configuration", "image", target.Name())
	ref, err := reference.ParseNormalizedNamed(target.Name())
	if err != nil {
		return "", status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}
	host := reference.Domain(ref)
	log.Trace("auth host", "host", host)

	// Parse our old Docker config and add the auth.
	log.Trace("loading Docker configuration")
	file, err := config.Load(config.Dir())
	if err != nil {
		return "", err
	}

	if file.AuthConfigs == nil {
		file.AuthConfigs = map[string]types.AuthConfig{}
	}
	file.AuthConfigs[host] = authCfg

	// Create a temporary directory for our config
	td, err := ioutil.TempDir("", "wp-docker-config")
	if err != nil {
		return "", status.Errorf(codes.Internal,
			"Failed to create temporary directory for Docker config: %s", err)
	}

	// Create a temporary file and write our Docker config to it
	f, err := os.Create(filepath.Join(td, "config.json"))
	if err != nil {
		return "", status.Errorf(codes.Internal,
			"Failed to create temporary file for Docker config: %s", err)
	}
	defer f.Close()
	if err := file.SaveToWriter(f); err != nil {
		return "", status.Errorf(codes.Internal,
			"Failed to create temporary file for Docker config: %s", err)
	}

	log.Info("temporary Docker config created for auth",
		"auth_host", host,
		"path", td,
	)

	return td, nil
}
