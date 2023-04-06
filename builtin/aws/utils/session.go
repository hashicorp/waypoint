// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
)

func GetSession(c *SessionConfig) (*session.Session, error) {
	config := aws.NewConfig().WithRegion(c.Region)

	if c.Logger != nil {
		l := c.Logger

		switch {
		case l.IsDebug():
			config = config.WithLogLevel(aws.LogDebug)
		case l.IsTrace():
			config = config.WithLogLevel(aws.LogDebugWithRequestRetries)
		}
	}

	return session.NewSessionWithOptions(session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	})
}

type SessionConfig struct {
	Region string
	Logger hclog.Logger
}
