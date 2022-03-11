package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/blake2b"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

const (
	// The username of the initial user created during bootstrapping. This
	// also is the user that Waypoint server versions prior to 0.5 used with
	// their token, so we use this to detect that scenario as well.
	DefaultUser = serverstate.DefaultUser

	// The ID of the initial user created during bootstrapping.
	DefaultUserId = serverstate.DefaultUserId

	// The identifier for the default key to use to generating tokens.
	DefaultKeyId = "k1"

	// Used as a byte sequence prepended to the encoded TokenTransport to identify
	// the token as valid before attempting to decode it. This is mostly a nicity to improve
	// understanding of the token data and error messages.
	tokenMagic = "wp24"

	// The size in bytes that the HMAC keys should be. Each key will contain this number of bytes
	// of data from rand.Reader
	hmacKeySize = 32

	// A prefix added to the key id when looking up the HMAC key from the database
	dbKeyPrefix = "hmacKey:"
)

// newToken is the generic internal function to create and encode a new
// token. The final parameter "body" should be set to the initial value
// of the token body, most importantly the "Kind" field should be set.
func newToken(
	state serverstate.Interface,
	duration time.Duration,
	keyId string,
	metadata map[string]string,
	body *pb.Token,
) (string, error) {
	body.IssuedTime = ptypes.TimestampNow()

	// If this token expires at some point, set an expiry
	if duration > 0 {
		now := time.Now().UTC().Add(duration)
		body.ValidUntil = &timestamp.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		}
	}

	// Set the accessor ID
	body.AccessorId = make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, body.AccessorId)
	if err != nil {
		return "", err
	}

	return encodeToken(state, keyId, metadata, body)
}

// Encode the given token with the given key and metadata.
// keyId controls which key is used to sign the key (key values are generated lazily).
// metadata is attached to the token transport as configuration style information
func encodeToken(state serverstate.Interface, keyId string, metadata map[string]string, body *pb.Token) (string, error) {
	// Get the key material
	key, err := state.HMACKeyCreateIfNotExist(keyId, hmacKeySize)
	if err != nil {
		return "", err
	}

	// Proto encode the body, this is what we sign.
	bodyData, err := proto.Marshal(body)
	if err != nil {
		return "", err
	}

	// Sign it
	h, err := blake2b.New256(key.Key)
	if err != nil {
		return "", err
	}
	h.Write(bodyData)

	// Build our wrapper which is not signed or encrypted.
	var tt pb.TokenTransport
	tt.Body = bodyData
	tt.KeyId = keyId
	tt.Metadata = metadata
	tt.Signature = h.Sum(nil)

	// Marshal the wrapper and base58 encode it.
	ttData, err := proto.Marshal(&tt)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.WriteString(tokenMagic)
	buf.Write(ttData)

	return base58.Encode(buf.Bytes()), nil
}

type tokenKey struct{}

// cookieFromRequest returns the server cookie value provided during the request,
// or blank if none (or a blank cookie) is provided.
func cookieFromRequest(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if c, ok := md["wpcookie"]; ok && len(c) > 0 {
			return c[0]
		}
	}

	return ""
}

// tokenFromContext returns the validated token used with the request.
// The token is guaranteed to be valid, meaning that it successfully
// was signed and decrypted. This will return nil if no token was present
// for the request.
func tokenFromContext(s Service, ctx context.Context) *pb.Token {
	value, ok := ctx.Value(tokenKey{}).(*pb.Token)
	if !ok && s.SuperUser() {
		// We are in implicit superuser mode meaning everything is always
		// allowed. Create a login token for the superuser.
		value = &pb.Token{
			Kind: &pb.Token_Login_{
				Login: &pb.Token_Login{
					UserId: DefaultUserId,
				},
			},
		}
	}

	return value
}
