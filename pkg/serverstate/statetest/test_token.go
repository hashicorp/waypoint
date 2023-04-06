// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	tests["token"] = []testFunc{
		TestTokenSignature,
	}
}

func TestTokenSignature(t *testing.T, factory Factory, rf RestartFactory) {
	ctx := context.Background()
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	body1 := []byte("test1")
	keyId := "k1"

	// Can generate a signature
	sig1, err := s.TokenSignature(ctx, body1, keyId)
	require.NoError(err)
	require.NotEmpty(sig1)

	// Good signature verification succeeds
	{
		valid, err := s.TokenSignatureVerify(ctx, body1, sig1, keyId)
		require.NoError(err)
		require.True(valid)
	}

	// Tampered body verification fails
	{
		valid, err := s.TokenSignatureVerify(ctx, []byte("test2"), sig1, keyId)
		require.NoError(err)
		require.False(valid)
	}

	// Tampered signature verification fails
	{
		valid, err := s.TokenSignatureVerify(ctx, body1, []byte("tampered signature"), keyId)
		require.NoError(err)
		require.False(valid)
	}

	// Different key sig verify fails
	{
		_, err := s.TokenSignatureVerify(ctx, body1, sig1, "k2")
		require.Error(err)
	}

}
