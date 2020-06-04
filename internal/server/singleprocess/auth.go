package singleprocess

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"io"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

const (
	DefaultUser  = "waypoint"
	DefaultKeyId = "k1"

	tokenMagic  = "wp24"
	hmacKeySize = 32
	dbKeyPrefix = "hmacKey:"
)

var (
	ErrInvalidToken = errors.New("invalid authentication token")

	authBucket = []byte("auth")
)

func init() {
	dbBuckets = append(dbBuckets, authBucket)
}

// DecodeToken parses the string and validates it as a valid token. If the token
// has a validity period attached to it, the period is checked here.
func (s *service) DecodeToken(token string) (*pb.TokenTransport, *pb.Token, error) {
	data, err := base58.Decode(token)
	if err != nil {
		return nil, nil, err
	}

	if subtle.ConstantTimeCompare(data[:len(tokenMagic)], []byte(tokenMagic)) != 1 {
		return nil, nil, errors.Wrapf(ErrInvalidToken, "bad magic")
	}

	var tt pb.TokenTransport

	err = proto.Unmarshal(data[len(tokenMagic):], &tt)
	if err != nil {
		return nil, nil, err
	}

	var hmacKey []byte

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(authBucket)
		if b == nil {
			return errors.Wrapf(ErrInvalidToken, "unknown key")
		}

		hmacKey = b.Get([]byte(dbKeyPrefix + tt.KeyId))
		if hmacKey == nil {
			return errors.Wrapf(ErrInvalidToken, "unknown key")
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	h, err := blake2b.New256(hmacKey)
	if err != nil {
		return nil, nil, err
	}

	h.Write(tt.Body)

	sum := h.Sum(nil)

	if subtle.ConstantTimeCompare(sum, tt.Signature) != 1 {
		return nil, nil, errors.Wrapf(ErrInvalidToken, "bad signature")
	}

	var body pb.Token

	err = proto.Unmarshal(tt.Body, &body)
	if err != nil {
		return nil, nil, err
	}

	if body.ValidUntil != nil {
		vt := time.Unix(body.ValidUntil.Seconds, int64(body.ValidUntil.Nanos))

		now := time.Now()
		if vt.Before(now.UTC()) {
			return nil, nil, errors.Wrapf(ErrInvalidToken, "token has expired. %s < %s", vt, now)
		}
	}

	return &tt, &body, nil
}

// Authenticate checks if the given endpoint should be allowed. This is called during a
// gRPC request. Effects is some information about the endpoint, at present these are
// either ["readonly"] or ["mutable"] to indicate if the endpoint will be only reading
// data or also mutating it.
func (s *service) Authenticate(ctx context.Context, token, endpoint string, effects []string) error {
	_, body, err := s.DecodeToken(token)
	if err != nil {
		return err
	}

	if !body.Login {
		return ErrInvalidToken
	}

	// TODO When we have a user model, this is where you'll check for the user.
	if body.User != DefaultUser {
		return ErrInvalidToken
	}

	return nil
}

// Generate a new token by signing the data in body.
// keyId controls which key is used to sign the key (key values are generated lazily).
// metadata is attached to the token transport as configuration style information
func (s *service) GenerateToken(keyId string, metadata map[string]string, body *pb.Token) (string, error) {
	var hmacKey []byte

	dbKey := []byte(dbKeyPrefix + keyId)

	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(authBucket)
		if err != nil {
			return err
		}

		hmacKey = b.Get(dbKey)
		if hmacKey == nil {
			hmacKey = make([]byte, hmacKeySize)

			_, err = io.ReadFull(rand.Reader, hmacKey)
			if err != nil {
				return err
			}

			return b.Put(dbKey, hmacKey)
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	bodyData, err := proto.Marshal(body)
	if err != nil {
		return "", err
	}

	h, err := blake2b.New256(hmacKey)
	if err != nil {
		return "", err
	}

	h.Write(bodyData)

	var tt pb.TokenTransport
	tt.Body = bodyData
	tt.KeyId = keyId
	tt.Metadata = metadata
	tt.Signature = h.Sum(nil)

	ttData, err := proto.Marshal(&tt)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.WriteString(tokenMagic)
	buf.Write(ttData)

	return base58.Encode(buf.Bytes()), nil
}

// Create a new login token.
// keyId controls which key is used to sign the key (key values are generated lazily).
// metadata is attached to the token transport as configuration style information
func (s *service) GenerateLoginToken(keyId string, metadata map[string]string) (string, error) {
	var body pb.Token
	body.Login = true
	body.User = DefaultUser
	body.TokenId = make([]byte, 16)

	_, err := io.ReadFull(rand.Reader, body.TokenId)
	if err != nil {
		return "", err
	}

	return s.GenerateToken(keyId, metadata, &body)
}

// Generate a token in the default way. This token is presented to the user when running
// `waypoint server`.
func (s *service) DefaultToken() (string, error) {
	return s.GenerateLoginToken(DefaultKeyId, nil)
}

// Create a new invite token. The duration controls for how long the invite token is valid.
// keyId controls which key is used to sign the key (key values are generated lazily).
// metadata is attached to the token transport as configuration style information
func (s *service) GenerateInviteToken(duration time.Duration, keyId string, metadata map[string]string) (string, error) {
	var body pb.Token
	body.Invite = true
	body.TokenId = make([]byte, 16)

	now := time.Now().UTC().Add(duration)
	body.ValidUntil = &timestamp.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.Nanosecond()),
	}

	_, err := io.ReadFull(rand.Reader, body.TokenId)
	if err != nil {
		return "", err
	}

	return s.GenerateToken(keyId, metadata, &body)
}

// Given an invite token, validate it and return a login token
func (s *service) ExchangeInvite(keyId, invite string) (string, error) {
	tt, body, err := s.DecodeToken(invite)
	if err != nil {
		return "", err
	}

	if !body.Invite {
		return "", errors.Wrapf(ErrInvalidToken, "not an invite token")
	}

	return s.GenerateLoginToken(keyId, tt.Metadata)
}
