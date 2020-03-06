package runner

import (
	"github.com/flynn/noise"
	"github.com/hashicorp/securetunnel"
)

type Tunnel struct {
	Host   string
	Params *securetunnel.TunnelParams

	key noise.DHKey
}

func CreateTunnel(host string) (*Tunnel, error) {
	params, err := securetunnel.CreateTunnel(securetunnel.TunnelOptions{
		Host: host,
	})

	if err != nil {
		return nil, err
	}

	key, err := securetunnel.GenerateKey()
	if err != nil {
		return nil, err
	}

	return &Tunnel{host, params, key}, nil
}

func (t *Tunnel) ServerToken() string {
	return t.Params.DestinationToken
}

func (t *Tunnel) ServerKey() string {
	return securetunnel.PublicKey(t.key)
}

func (t *Tunnel) Connect() (*securetunnel.Tunnel, error) {
	return securetunnel.Open(t.Params.SourceToken, t.key)
}

func (t *Tunnel) Close() error {
	return securetunnel.DeleteTunnel(t.Params.SourceToken)
}
