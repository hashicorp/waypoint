package server

import (
	"crypto/rand"

	"github.com/oklog/ulid"
)

//go:generate sh -c "protoc -I../../vendor/proto/api-common-protos -Iproto/ proto/*.proto --go_out=plugins=grpc:gen/"

var ulidReader = ulid.Monotonic(rand.Reader, 1)

// Id returns a unique Id that can be used for new values. This generates
// a ulid value but the ID itself should be an internal detail. An error will
// be returned if the ID could be generated.
func Id() (string, error) {
	id, err := ulid.New(ulid.Now(), ulidReader)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
