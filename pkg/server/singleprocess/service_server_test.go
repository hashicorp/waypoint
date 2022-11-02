package singleprocess

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server/boltdbstate"
	"github.com/hashicorp/waypoint/internal/serverconfig"
	"github.com/hashicorp/waypoint/pkg/server"
)

func TestServerConfigWithStartupConfig(t *testing.T) {
	ctx := context.Background()
	cfg := &serverconfig.Config{
		CEBConfig: &serverconfig.CEBConfig{
			Addr:          "myendpoint",
			TLSEnabled:    false,
			TLSSkipVerify: true,
		},
	}

	db := testDB(t)
	// Create our server
	impl, err := New(
		WithDB(db),
		WithConfig(cfg),
	)
	require.NoError(t, err)
	_ = server.TestServer(t, impl)

	st, err := boltdbstate.New(hclog.L(), db)
	require.NoError(t, err)

	t.Run("Check config defaults are set", func(t *testing.T) {
		require := require.New(t)

		retCfg, err := st.ServerConfigGet(ctx)
		require.NoError(err)
		require.NotNil(retCfg)

		addr := retCfg.AdvertiseAddrs[0]
		require.Equal(cfg.CEBConfig.Addr, addr.Addr)
		require.Equal(cfg.CEBConfig.TLSEnabled, addr.Tls)
		require.Equal(cfg.CEBConfig.TLSSkipVerify, addr.TlsSkipVerify)
	})
}

func TestServerConfigWithNoStartupConfig(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	// Create our server
	impl, err := New(
		WithDB(db),
	)
	require.NoError(t, err)
	_ = server.TestServer(t, impl)

	st, err := boltdbstate.New(hclog.L(), db)
	require.NoError(t, err)

	t.Run("Check config defaults are not set", func(t *testing.T) {
		require := require.New(t)

		retCfg, err := st.ServerConfigGet(ctx)
		require.NoError(err)
		require.NotNil(retCfg)
		require.Len(retCfg.AdvertiseAddrs, 0)
	})
}
