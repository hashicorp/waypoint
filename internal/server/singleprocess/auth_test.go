package singleprocess

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/blake2b"
)

func TestServiceAuth(t *testing.T) {
	// Create our server
	impl, err := New(testDB(t))
	require.NoError(t, err)

	t.Run("create and validate a token", func(t *testing.T) {
		s := impl.(*service)

		md := map[string]string{
			"addr": "test",
		}

		token, err := s.GenerateLoginToken(DefaultKeyId, md)
		require.NoError(t, err)

		require.True(t, len(token) > 5)
		t.Logf("token: %s", token)

		tt, body, err := s.DecodeToken(token)
		require.NoError(t, err)

		assert.True(t, body.Login)
		assert.False(t, body.Invite)
		assert.Equal(t, DefaultUser, body.User)
		assert.Equal(t, md, tt.Metadata)

		err = s.Authenticate(context.Background(), token, "test", nil)
		require.NoError(t, err)

		// Now corrupt the token and check that validation fails
		data := []byte(token)

		data[len(data)-2] = data[len(data)-2] + 1
		_, _, err = s.DecodeToken(string(data))
		require.Error(t, err)

		// Generate a legit token with an unknown key though
	})

	t.Run("rejects tokens signed with unknown keys", func(t *testing.T) {
		s := impl.(*service)

		md := map[string]string{
			"addr": "test",
		}

		token, err := s.GenerateLoginToken(DefaultKeyId, md)
		require.NoError(t, err)

		require.True(t, len(token) > 5)
		t.Logf("token: %s", token)

		tt, body, err := s.DecodeToken(token)
		require.NoError(t, err)
		bodyData, err := proto.Marshal(body)
		require.NoError(t, err)

		h, err := blake2b.New256([]byte("abcdabcdabcdabcadacdacdaaa"))
		require.NoError(t, err)

		h.Write(bodyData)

		tt.Signature = h.Sum(nil)

		ttData, err := proto.Marshal(tt)
		require.NoError(t, err)

		var buf bytes.Buffer
		buf.WriteString(tokenMagic)
		buf.Write(ttData)

		rogue := base58.Encode(buf.Bytes())

		_, _, err = s.DecodeToken(rogue)
		require.Error(t, err)
	})

	t.Run("exchange an invite token", func(t *testing.T) {
		s := impl.(*service)

		invite, err := s.GenerateInviteToken(2*time.Second, DefaultKeyId, nil)
		require.NoError(t, err)

		lt, err := s.ExchangeInvite(DefaultKeyId, invite)
		require.NoError(t, err)

		_, body, err := s.DecodeToken(lt)
		require.NoError(t, err)

		assert.True(t, body.Login)

		time.Sleep(3 * time.Second)

		_, err = s.ExchangeInvite(DefaultKeyId, invite)
		require.Error(t, err)
	})
}
