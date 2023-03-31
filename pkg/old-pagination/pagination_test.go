package pagination

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeAndParsePageToken(t *testing.T) {
	require := require.New(t)

	t.Run("Decodes and parses properly encoded + formatted page token", func(t *testing.T) {
		paginationToken, _ := EncodeAndSerializePageToken("key", "value")
		key, value, err := DecodeAndParsePageToken(paginationToken)
		require.NoError(err)
		require.Equal("key", key)
		require.Equal("value", value)
	})

	t.Run("Errors if pagination token is incorrectly formatted", func(t *testing.T) {
		t.Run("Not base64 encoded", func(t *testing.T) {
			_, _, err := DecodeAndParsePageToken("thisIsNotBase64Encoded")
			require.Error(err)
			require.EqualError(err, "Incorrectly formatted pagination token.")
		})
		t.Run("Incorrectly formatted", func(t *testing.T) {
			paginationToken := base64.StdEncoding.EncodeToString([]byte("incorrectlyFormattedToken"))
			_, _, err := DecodeAndParsePageToken(paginationToken)
			require.Error(err)
			require.EqualError(err, "Incorrectly formatted pagination token.")
		})
	})
}

func TestEncodeAndSerializePageToken(t *testing.T) {
	require := require.New(t)

	t.Run("Encodes and serializes page token", func(t *testing.T) {
		key := "key"
		value := "value"
		expectedPageToken := base64.StdEncoding.EncodeToString([]byte(key + ":" + value))

		pageToken, err := EncodeAndSerializePageToken(key, value)
		require.NoError(err)
		require.Equal(expectedPageToken, pageToken)
	})

	t.Run("Returns empty string if key or value are empty string", func(t *testing.T) {
		pageToken, err := EncodeAndSerializePageToken("", "")
		require.NoError(err)
		require.Equal("", pageToken)
	})
}
