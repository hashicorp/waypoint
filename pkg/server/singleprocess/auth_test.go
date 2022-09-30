package singleprocess

import (
	"bytes"
	"context"
	"testing"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/blake2b"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestServiceAuth(t *testing.T) {
	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	s := impl.(*Service)
	ctx := context.Background()

	// "Log in" a default user
	ctx = UserWithContext(ctx, &pb.User{Id: DefaultUserId})

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

		user := s.UserFromContext(ctx)
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
		_, body, err := s.decodeToken(ctx, token)
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
		_, _, err = s.decodeToken(ctx, string(data))
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
		user := s.UserFromContext(ctx)
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
			user := s.UserFromContext(ctx)
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

	t.Run("entrypoint token can only access entrypoint APIs (LEGACY as another user)", func(t *testing.T) {
		require := require.New(t)

		// Log in as some other user
		ctx = UserWithContext(ctx, &pb.User{Id: DefaultUserId + "0"})

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

		tt, body, err := s.decodeToken(ctx, token)
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
		buf.WriteString("wp24")
		buf.Write(ttData)

		rogue := base58.Encode(buf.Bytes())

		_, _, err = s.decodeToken(ctx, rogue)
		require.Error(err)
	})

	t.Run("validate a runner token with no ID set", func(t *testing.T) {
		require := require.New(t)

		token, err := s.newToken(ctx, 0, DefaultKeyId, nil, &pb.Token{
			Kind: &pb.Token_Runner_{
				Runner: &pb.Token_Runner{
					// Being explicit (not setting it at all would be the same)
					// to show what we're testing.
					Id: "",
				},
			},
		})

		// Verify authing works
		_, err = s.Authenticate(context.Background(), token, "test", nil)
		require.NoError(err)
	})

	t.Run("validate a runner token with an ID set and not adopted", func(t *testing.T) {
		require := require.New(t)

		token, err := s.newToken(ctx, 0, DefaultKeyId, nil, &pb.Token{
			Kind: &pb.Token_Runner_{
				Runner: &pb.Token_Runner{
					Id: "i-do-not-exist",
				},
			},
		})

		// Auth should NOT work
		_, err = s.Authenticate(context.Background(), token, "test", nil)
		require.Error(err)
	})

	t.Run("validate a runner token with an ID set and label hash mismatch", func(t *testing.T) {
		require := require.New(t)

		labels := map[string]string{"foo": "bar"}

		// Create a runner and adopt it.
		require.NoError(s.state(ctx).RunnerCreate(&pb.Runner{
			Id:     "A",
			Labels: labels,
			Kind: &pb.Runner_Remote_{
				Remote: &pb.Runner_Remote{},
			},
		}))
		defer s.state(ctx).RunnerDelete("A")
		require.NoError(s.state(ctx).RunnerAdopt("A", false))

		token, err := s.newToken(ctx, 0, DefaultKeyId, nil, &pb.Token{
			Kind: &pb.Token_Runner_{
				Runner: &pb.Token_Runner{
					Id:        "A",
					LabelHash: 42,
				},
			},
		})

		// Auth should NOT work
		_, err = s.Authenticate(context.Background(), token, "test", nil)
		require.Error(err)
	})

	t.Run("validate a runner token with an ID set and label hash good match", func(t *testing.T) {
		require := require.New(t)

		labels := map[string]string{"foo": "bar"}
		hash, err := serverptypes.RunnerLabelHash(labels)
		require.NoError(err)

		// Create a runner and adopt it.
		require.NoError(s.state(ctx).RunnerCreate(&pb.Runner{
			Id:     "A",
			Labels: labels,
			Kind: &pb.Runner_Remote_{
				Remote: &pb.Runner_Remote{},
			},
		}))
		defer s.state(ctx).RunnerDelete("A")
		require.NoError(s.state(ctx).RunnerAdopt("A", false))

		token, err := s.newToken(ctx, 0, DefaultKeyId, nil, &pb.Token{
			Kind: &pb.Token_Runner_{
				Runner: &pb.Token_Runner{
					Id:        "A",
					LabelHash: hash,
				},
			},
		})

		// Auth should work
		_, err = s.Authenticate(context.Background(), token, "test", nil)
		require.NoError(err)
	})
}

func TestServiceAuth_TriggerToken(t *testing.T) {
	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	s := impl.(*Service)
	ctx := context.Background()

	// "Log in" a default user
	ctx = UserWithContext(ctx, &pb.User{Id: DefaultUserId})

	t.Run("create and validate new token cannot authenticate grpc endpoint", func(t *testing.T) {
		require := require.New(t)

		resp, err := s.GenerateLoginToken(ctx, &pb.LoginTokenRequest{
			Trigger: true,
		})
		require.NoError(err)
		token := resp.Token

		require.True(len(token) > 5)
		t.Logf("token: %s", token)

		// Test some internal state of the token
		_, body, err := s.decodeToken(ctx, token)
		require.NoError(err)
		kind, ok := body.Kind.(*pb.Token_Trigger_)
		assert.True(t, ok)
		assert.Equal(t, DefaultUserId, kind.Trigger.FromUserId)

		// Verify authing won't work currently
		_, err = s.Authenticate(context.Background(), token, "test", nil)
		require.Error(err)
		e, _ := status.FromError(err)
		assert.Equal(t, e.Code(), codes.PermissionDenied)
	})

	// TODO(briancain): Add tests for HTTP endpoint when implemeneted
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
	s := impl.(*Service)
	ctx := context.Background()

	user := s.UserFromContext(ctx)
	require.NotNil(t, user)
	require.Equal(t, DefaultUserId, user.Id)
}
