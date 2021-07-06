package singleprocess

import (
	"bytes"
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/blake2b"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestServiceAuth(t *testing.T) {
	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	s := impl.(*service)
	ctx := context.Background()

	// "Log in" a default user
	ctx = userWithContext(ctx, &pb.User{Id: DefaultUserId})

	// Bootstrap
	var bootstrapToken string
	{
		resp, err := impl.BootstrapToken(ctx, &empty.Empty{})
		require.NoError(t, err)
		require.NotEmpty(t, resp.Token)
		bootstrapToken = resp.Token
	}

	t.Run("authenticate with gibberish", func(t *testing.T) {
		_, err := s.Authenticate(context.Background(), "hello!", "test", nil)
		require.Error(t, err)
	})

	t.Run("authenticate with gibberish to unauthenticated endpoint", func(t *testing.T) {
		_, err := s.Authenticate(context.Background(), "hello!", "GetVersionInfo", nil)
		require.NoError(t, err)
	})

	t.Run("authenticate with bootstrap token", func(t *testing.T) {
		ctx, err := s.Authenticate(context.Background(), bootstrapToken, "test", nil)
		require.NoError(t, err)

		user := s.userFromContext(ctx)
		require.NotNil(t, user)
		require.Equal(t, DefaultUserId, user.Id)
	})

	t.Run("create and validate a new token", func(t *testing.T) {
		require := require.New(t)

		resp, err := s.GenerateLoginToken(ctx, &pb.LoginTokenRequest{})
		require.NoError(err)
		token := resp.Token

		require.True(len(token) > 5)
		t.Logf("token: %s", token)

		// Test some internal state of the token
		_, body, err := s.decodeToken(token)
		require.NoError(err)
		kind, ok := body.Kind.(*pb.Token_Login_)
		assert.True(t, ok)
		assert.Equal(t, DefaultUserId, kind.Login.UserId)

		// Verify authing works
		_, err = s.Authenticate(context.Background(), token, "test", nil)
		require.NoError(err)

		// Now corrupt the token and check that validation fails
		data := []byte(token)
		data[len(data)-2] = data[len(data)-2] + 1
		_, _, err = s.decodeToken(string(data))
		require.Error(err)
	})

	t.Run("exchange an invite token", func(t *testing.T) {
		require := require.New(t)

		// Get our invite token for the entrypoint
		resp, err := s.GenerateInviteToken(ctx, &pb.InviteTokenRequest{
			Duration: "5m",
			Login: &pb.Token_Login{
				UserId: DefaultUserId,
			},
		})
		require.NoError(err)

		// Exchange it
		resp, err = s.ConvertInviteToken(ctx, &pb.ConvertInviteTokenRequest{
			Token: resp.Token,
		})
		require.NoError(err)
		token := resp.Token

		{
			_, err := s.Authenticate(context.Background(), token, "UpsertDeployment", nil)
			require.NoError(err)
		}
	})

	t.Run("exchange an invite token for new user", func(t *testing.T) {
		require := require.New(t)

		// Get our invite token for the entrypoint
		resp, err := s.GenerateInviteToken(ctx, &pb.InviteTokenRequest{
			Duration: "5m",
			Login: &pb.Token_Login{
				UserId: DefaultUserId,
			},
			Signup: &pb.Token_Invite_Signup{
				InitialUsername: "alice",
			},
		})
		require.NoError(err)

		// Exchange it
		resp, err = s.ConvertInviteToken(ctx, &pb.ConvertInviteTokenRequest{
			Token: resp.Token,
		})
		require.NoError(err)
		token := resp.Token

		// Auth
		ctx, err := s.Authenticate(context.Background(), token, "UpsertDeployment", nil)
		require.NoError(err)
		user := s.userFromContext(ctx)
		require.NotNil(user)
		require.NotEqual(DefaultUserId, user.Id)
		require.Equal("alice", user.Username)

		// Generate a login token for that user using the superuser
		{
			resp, err := s.GenerateLoginToken(ctx, &pb.LoginTokenRequest{
				User: &pb.Ref_User{
					Ref: &pb.Ref_User_Username{
						Username: &pb.Ref_UserUsername{
							Username: "alice",
						},
					},
				},
			})
			require.NoError(err)
			token := resp.Token

			// Verify authing works
			ctx, err := s.Authenticate(context.Background(), token, "test", nil)
			require.NoError(err)
			user := s.userFromContext(ctx)
			require.NotNil(t, user)
			require.NotEqual(t, DefaultUserId, user.Id)
			require.Equal("alice", user.Username)
		}
	})

	t.Run("entrypoint token can only access entrypoint APIs", func(t *testing.T) {
		require := require.New(t)

		// Get our invite token for the entrypoint
		resp, err := s.GenerateInviteToken(ctx, &pb.InviteTokenRequest{
			Duration: "5m",
			Login: &pb.Token_Login{
				UserId:     DefaultUserId,
				Entrypoint: &pb.Token_Entrypoint{DeploymentId: "A"},
			},
		})
		require.NoError(err)

		// Exchange it
		resp, err = s.ConvertInviteToken(ctx, &pb.ConvertInviteTokenRequest{
			Token: resp.Token,
		})
		require.NoError(err)
		token := resp.Token

		{
			_, err := s.Authenticate(context.Background(), token, "EntrypointConfig", nil)
			require.NoError(err)
		}

		{
			_, err := s.Authenticate(context.Background(), token, "UpsertDeployment", nil)
			require.Error(err)
		}
	})

	t.Run("entrypoint token can only access entrypoint APIs (LEGACY)", func(t *testing.T) {
		require := require.New(t)

		// Get our invite token for the entrypoint
		resp, err := s.GenerateInviteToken(ctx, &pb.InviteTokenRequest{
			Duration:         "5m",
			UnusedEntrypoint: &pb.Token_Entrypoint{DeploymentId: "A"},
		})
		require.NoError(err)

		// Exchange it
		resp, err = s.ConvertInviteToken(ctx, &pb.ConvertInviteTokenRequest{
			Token: resp.Token,
		})
		require.NoError(err)
		token := resp.Token

		{
			_, err := s.Authenticate(context.Background(), token, "EntrypointConfig", nil)
			require.NoError(err)
		}

		{
			_, err := s.Authenticate(context.Background(), token, "UpsertDeployment", nil)
			require.Error(err)
		}
	})

	t.Run("rejects tokens signed with unknown keys", func(t *testing.T) {
		require := require.New(t)

		resp, err := s.GenerateLoginToken(ctx, &pb.LoginTokenRequest{})
		require.NoError(err)
		token := resp.Token

		require.True(len(token) > 5)
		t.Logf("token: %s", token)

		tt, body, err := s.decodeToken(token)
		require.NoError(err)
		bodyData, err := proto.Marshal(body)
		require.NoError(err)

		h, err := blake2b.New256([]byte("abcdabcdabcdabcadacdacdaaa"))
		require.NoError(err)

		h.Write(bodyData)

		tt.Signature = h.Sum(nil)

		ttData, err := proto.Marshal(tt)
		require.NoError(err)

		var buf bytes.Buffer
		buf.WriteString(tokenMagic)
		buf.Write(ttData)

		rogue := base58.Encode(buf.Bytes())

		_, _, err = s.decodeToken(rogue)
		require.Error(err)
	})
}

func TestServiceBootstrapToken(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)

	{
		// Initial bootstrap should return a token
		resp, err := impl.BootstrapToken(ctx, &empty.Empty{})
		require.NoError(err)
		require.NotEmpty(resp.Token)
	}

	{
		// Subs calls should fail
		resp, err := impl.BootstrapToken(ctx, &empty.Empty{})
		require.Error(err)
		require.Equal(codes.PermissionDenied, status.Code(err))
		require.Nil(resp)
	}
}

func TestServiceDecodeToken(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)

	// Grab our bootstrap token, that'll work
	resp, err := impl.BootstrapToken(ctx, &empty.Empty{})
	require.NoError(err)
	require.NotEmpty(resp.Token)

	// Decode it
	decodeResp, err := impl.DecodeToken(ctx, &pb.DecodeTokenRequest{Token: resp.Token})
	require.NoError(err)
	require.NotNil(decodeResp.Token)
	require.NotNil(decodeResp.Transport)
}

func TestServiceAuth_userSuperuserForced(t *testing.T) {
	// Create our server
	impl, err := New(WithDB(testDB(t)), WithSuperuser())
	require.NoError(t, err)
	s := impl.(*service)
	ctx := context.Background()

	user := s.userFromContext(ctx)
	require.NotNil(t, user)
	require.Equal(t, DefaultUserId, user.Id)
}
