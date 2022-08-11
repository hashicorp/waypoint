package serverclient

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint/internal/clicontext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
	"github.com/hashicorp/waypoint/pkg/tokenutil"
)

func TestConnect(t *testing.T) {
	t.Run("supports authenticating with oauth2 tokens", func(t *testing.T) {
		var (
			token pb.Token
			tt    pb.TokenTransport
		)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		l, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)

		defer l.Close()

		caKey, caKeyPEM := generateKey(t)
		_, caCertPEM := generateRootCert(t, caKey)

		serverCert, err := tls.X509KeyPair(caCertPEM, caKeyPEM)
		require.NoError(t, err)

		gl, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
			Certificates: []tls.Certificate{serverCert},
		})
		require.NoError(t, err)

		defer gl.Close()

		token.Kind = &pb.Token_Login_{
			Login: &pb.Token_Login{
				UserId: "xxyyzz",
			},
		}

		tt.ExternalCreds = &pb.TokenTransport_OauthCreds{
			OauthCreds: &pb.TokenTransport_OAuthCredentials{
				Url:          "http://" + l.Addr().String(),
				ClientId:     "xxyyzz",
				ClientSecret: "aabbcc",
			},
		}

		tt.Body, err = proto.Marshal(&token)
		require.NoError(t, err)

		tt.Signature = []byte{0, 1, 2, 3, 4}

		ttData, err := proto.Marshal(&tt)
		require.NoError(t, err)

		data := append([]byte(tokenutil.TokenMagic), ttData...)

		strToken := base58.Encode(data)

		minToken, err := tokenutil.StripCreds(&tt)
		require.NoError(t, err)

		serv := http.Server{
			Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Set("Content-Type", "application/json")

				fmt.Printf("=> %s\n", r.URL.String())
				fmt.Fprintln(
					rw,
					`{ "access_token": "this-is-an-oauth-token", "token_type": "debug" }`,
				)
			}),
		}

		go serv.Serve(l)

		var (
			seenToken   []string
			seenWPToken []string
		)

		gs := grpc.NewServer(
			grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				md, ok := metadata.FromIncomingContext(ctx)
				if !ok {
					return nil, fmt.Errorf("nope 1")
				}

				seenToken = md["authorization"]
				seenWPToken = md[tokenutil.MetadataKey]

				return nil, fmt.Errorf("nope 2")
			}))

		pb.RegisterWaypointServer(gs, &pb.UnimplementedWaypointServer{})

		go gs.Serve(gl)

		cfg := &clicontext.Config{
			Server: serverconfig.Client{
				Address:       gl.Addr().String(),
				Tls:           true,
				TlsSkipVerify: true,
				RequireAuth:   true,
				AuthToken:     strToken,
			},
		}

		cc, err := Connect(ctx, FromContextConfig(cfg))
		require.NoError(t, err)

		client := pb.NewWaypointClient(cc)

		_, err = client.GetVersionInfo(ctx, &emptypb.Empty{})
		require.Error(t, err)

		assert.Equal(t, "debug this-is-an-oauth-token", seenToken[0])
		assert.Equal(t, minToken, seenWPToken[0])
	})
}

func generateRootCert(t *testing.T, key crypto.Signer) (*x509.Certificate, []byte) {
	subjectKeyIdentifier := calculateSubjectKeyIdentifier(t, key.Public())

	template := &x509.Certificate{
		SerialNumber: generateSerial(t),
		Subject: pkix.Name{
			Organization: []string{"Awesomeness, Inc."},
			Country:      []string{"US"},
			Locality:     []string{"San Francisco"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          subjectKeyIdentifier,
		AuthorityKeyId:        subjectKeyIdentifier,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	require.NoError(t, err)

	rootCert, err := x509.ParseCertificate(der)
	require.NoError(t, err)

	rootCertPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})

	return rootCert, rootCertPEM
}

// generateSerial generates a serial number using the maximum number of octets (20) allowed by RFC 5280 4.1.2.2
// (Adapted from https://github.com/cloudflare/cfssl/blob/828c23c22cbca1f7632b9ba85174aaa26e745340/signer/local/local.go#L407-L418)
func generateSerial(t *testing.T) *big.Int {
	serialNumber := make([]byte, 20)
	_, err := io.ReadFull(rand.Reader, serialNumber)
	require.NoError(t, err)

	return new(big.Int).SetBytes(serialNumber)
}

// calculateSubjectKeyIdentifier implements a common method to generate a key identifier
// from a public key, namely, by composing it from the 160-bit SHA-1 hash of the bit string
// of the public key (cf. https://tools.ietf.org/html/rfc5280#section-4.2.1.2).
// (Adapted from https://github.com/jsha/minica/blob/master/main.go).
func calculateSubjectKeyIdentifier(t *testing.T, pubKey crypto.PublicKey) []byte {
	spkiASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	require.NoError(t, err)

	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}
	_, err = asn1.Unmarshal(spkiASN1, &spki)
	require.NoError(t, err)

	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)
	return skid[:]
}

// generateKey generates a 1024-bit RSA private key
func generateKey(t *testing.T) (crypto.Signer, []byte) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	require.NoError(t, err)

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return key, keyPEM
}
