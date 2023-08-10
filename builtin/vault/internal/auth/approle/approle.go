// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package approle

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"

	"github.com/hashicorp/waypoint/builtin/vault/internal/auth"
)

type approleMethod struct {
	logger    hclog.Logger
	mountPath string

	roleId   string
	secretId string
}

// NewApproleAuthMethod reads the user configuration and returns a configured
// AuthMethod
func NewApproleAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &approleMethod{
		logger:    conf.Logger,
		mountPath: conf.MountPath,
	}

	roleIdRaw, ok := conf.Config["role_id"]
	if !ok {
		return nil, errors.New("missing 'role_id' value")
	}
	a.roleId, ok = roleIdRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role_id' value into string")
	}

	secretIdRaw, ok := conf.Config["secret_id"]
	if !ok {
		return nil, errors.New("missing 'secret_id' value")
	}
	a.secretId, ok = secretIdRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'secret_id' value into string")
	}

	return a, nil
}

func (a *approleMethod) Authenticate(ctx context.Context, client *api.Client) (retPath string, header http.Header, retData map[string]interface{}, retError error) {
	a.logger.Trace("beginning authentication")
	return fmt.Sprintf("auth/%s/login", a.mountPath), nil, map[string]interface{}{
		"role_id":   a.roleId,
		"secret_id": a.secretId,
	}, nil
}

func (a *approleMethod) NewCreds() chan struct{} {
	return nil
}

func (a *approleMethod) CredSuccess() {
}

func (a *approleMethod) Shutdown() {
}
