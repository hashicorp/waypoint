// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package tokenutil

import (
	"crypto/subtle"

	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var (
	ErrInvalidToken = errors.New("invalid token detected")
)

const (
	// TokenMagic is used as a byte sequence prepended to the encoded TokenTransport to identify
	// the token as valid before attempting to decode it. This is mostly a nicity to improve
	// understanding of the token data and error messages.
	TokenMagic = "wp24"
)

// TokenDecode provides the ability to decode a waypoint token into
// the embedded protobuf. WARNING: This function is unable to
// verify the token, only the server that generate the server can do that.
func TokenDecode(token string) (*pb.TokenTransport, *pb.Token, error) {
	data, err := base58.Decode(token)
	if err != nil {
		return nil, nil, err
	}

	if len(data) < len(TokenMagic) || subtle.ConstantTimeCompare(data[:len(TokenMagic)], []byte(TokenMagic)) != 1 {
		return nil, nil, errors.Wrapf(ErrInvalidToken, "bad magic")
	}

	var tt pb.TokenTransport
	err = proto.Unmarshal(data[len(TokenMagic):], &tt)
	if err != nil {
		return nil, nil, err
	}

	// Decode the actual token structure
	var body pb.Token
	err = proto.Unmarshal(tt.Body, &body)
	if err != nil {
		return nil, nil, err
	}
	return &tt, &body, nil
}

// StripCreds removes the credentials from the given TokenTransport and repackages
// it up as a string. This doesn't invalidate the signature, because the signature
// is only against the body field of the transport.
func StripCreds(tt *pb.TokenTransport) (string, error) {
	tmp := pb.TokenTransport{
		Body:      tt.Body,
		Signature: tt.Signature,
		KeyId:     tt.KeyId,
		Metadata:  tt.Metadata,
	}

	data, err := proto.Marshal(&tmp)
	if err != nil {
		return "", err
	}

	magicified := append([]byte(TokenMagic), data...)

	return base58.Encode(magicified), nil
}
