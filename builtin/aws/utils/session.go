package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetSession(c *SessionConfig) (*session.Session, error) {
	config := aws.NewConfig().WithRegion(c.Region)

	return session.NewSessionWithOptions(session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	})
}

type SessionConfig struct {
	Region string
}
