package httpapi

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleTrigger(t *testing.T) {
	//ctx := context.Background()

	// Get our gRPC server
	impl := &execImpl{}
	addr := testServer(t, impl)

	// Start up our test HTTP server
	httpServer := httptest.NewServer(HandleTrigger(addr, false))
	defer httpServer.Close()

	t.Run("a request with a non-valid trigger", func(t *testing.T) {
		require := require.New(t)

		require.Equal(1, 1)
	})
}
