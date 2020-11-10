package vault

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/pointerstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

type ConfigSourcer struct {
	// Client, if set, will be used as the client instead of initializing
	// based on the config. This is only used for tests.
	Client *vaultapi.Client

	resultsLock sync.Mutex
	results     map[string]*pb.ConfigSource_Value
}

// ReadFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) ReadFunc() interface{} {
	return cs.read
}

// StopFunc implements component.ConfigSourcer
func (cs *ConfigSourcer) StopFunc() interface{} {
	return nil
}

func (cs *ConfigSourcer) read(
	ctx context.Context,
	log hclog.Logger,
	reqs []*component.ConfigRequest,
) ([]*pb.ConfigSource_Value, error) {
	client := cs.Client
	if client == nil {
		// Initialize the client
		log.Debug("initializing the Vault client")
		clientConfig := vaultapi.DefaultConfig()
		err := clientConfig.ReadEnvironment()
		if err != nil {
			return nil, err
		}

		client, err = vaultapi.NewClient(clientConfig)
		if err != nil {
			return nil, err
		}
	} else {
		log.Debug("using preconfigured client on struct")
	}

	// We keep track of the secrets by path so that we only request each path once.
	secrets := map[string]*vaultapi.Secret{}

	// Go through each variable and read it
	cs.results = map[string]*pb.ConfigSource_Value{}
	for _, req := range reqs {
		result := &pb.ConfigSource_Value{Name: req.Name}
		cs.results[req.Name] = result

		// Decode our configuration
		var vaultReq reqConfig
		if err := mapstructure.WeakDecode(req.Config, &vaultReq); err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}
		L := log.With("path", vaultReq.Path, "key", vaultReq.Key)

		// Get this secret or read it if we haven't already.
		secret, ok := secrets[vaultReq.Path]
		if !ok {
			var err error
			L.Trace("querying Vault secret")
			secret, err = client.Logical().Read(vaultReq.Path)
			if err != nil {
				result.Result = &pb.ConfigSource_Value_Error{
					Error: status.New(codes.Aborted, err.Error()).Proto(),
				}

				continue
			}

			secrets[vaultReq.Path] = secret
		}

		// Get the value
		if !strings.HasPrefix(vaultReq.Key, "/") {
			vaultReq.Key = "/" + vaultReq.Key
		}
		value, err := pointerstructure.Get(secret.Data, vaultReq.Key)
		if err != nil {
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		// Convert it to a string
		var valueStr string
		if err := mapstructure.WeakDecode(value, &valueStr); err != nil {
			L.Warn("vault secret value couldn't be converted to string")
			result.Result = &pb.ConfigSource_Value_Error{
				Error: status.New(codes.Aborted, err.Error()).Proto(),
			}

			continue
		}

		result.Result = &pb.ConfigSource_Value_Value{
			Value: valueStr,
		}
	}

	list := make([]*pb.ConfigSource_Value, 0, len(cs.results))
	for _, r := range cs.results {
		list = append(list, r)
	}

	return list, nil
}

func (cs *ConfigSourcer) stop() error {
	return nil
}

type reqConfig struct {
	Path string
	Key  string
}
