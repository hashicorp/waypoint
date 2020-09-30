package protocol

import (
	"errors"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var (
	ErrClientOutdated = errors.New("client protocol version is too outdated for the server")
	ErrServerOutdated = errors.New("server minimum protocol version is too outdated for the client")
)

// Negotiate takes two protocol versions and determines the value to use.
// If negotiation is impossible, an error is returned. The error value is
// one of the exported variables in this file.
func Negotiate(client, server *pb.VersionInfo_ProtocolVersion) (uint32, error) {
	// If the client is too old, then it is an error
	if client.Current < server.Minimum {
		return 0, ErrClientOutdated
	}

	// If the server is too old, also an error
	if server.Current < client.Minimum {
		return 0, ErrServerOutdated
	}

	// Determine our shared protocol number. We use the maximum protocol
	// that we both support.
	version := server.Current
	if version > client.Current {
		version = client.Current
	}

	return version, nil
}
