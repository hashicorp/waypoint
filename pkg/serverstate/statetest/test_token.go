package statetest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	tests["token"] = []testFunc{
		TestTokenSignature,
	}
}

func TestTokenSignature(t *testing.T, factory Factory, rf RestartFactory) {
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	body1 := []byte("test1")
	keyId := "k1"

	// Can generate a signature
	sig1, err := s.TokenSignature(body1, keyId)
	require.NoError(err)
	require.NotEmpty(sig1)

	// Good signature verification succeeds
	{
		valid, err := s.TokenSignatureVerify(body1, sig1, keyId)
		require.NoError(err)
		require.True(valid)
	}

	// Tampered body verification fails
	{
		valid, err := s.TokenSignatureVerify([]byte("test2"), sig1, keyId)
		require.NoError(err)
		require.False(valid)
	}

	// Tampered signature verification fails
	{
		valid, err := s.TokenSignatureVerify(body1, []byte("tampered signature"), keyId)
		require.NoError(err)
		require.False(valid)
	}

	// Different key sig verify fails
	{
		_, err := s.TokenSignatureVerify(body1, sig1, "k2")
		require.Error(err)
	}

}
