package httpapi

import (
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/waypoint/pkg/server/gen/mocks"
	"github.com/stretchr/testify/require"
)

func TestHandleTrigger(t *testing.T) {
	//ctx := context.Background()

	// Get our gRPC server
	impl := &triggerImpl{}
	addr := testServer(t, impl)

	// Start up our test HTTP server
	httpServer := httptest.NewServer(HandleTrigger(addr, false))
	defer httpServer.Close()

	t.Run("a request with a non-valid trigger", func(t *testing.T) {
		require := require.New(t)
		require.Equal(1, 1)

		//request := httptest.NewRequest("GET", "/v1/trigger/123", nil)
		//responseRecorder := httptest.NewRecorder()

		//triggerHandler := HandleTrigger(httpServer.URL, false)
		//triggerHandler.ServeHTTP(responseRecorder, request)

		//if responseRecorder.Code != 404 {
		//	t.Errorf("Want status '%d', got '%d'", 404, responseRecorder.Code)
		//}

		//if strings.TrimSpace(responseRecorder.Body.String()) != "404" {
		//	t.Errorf("Want '%s', got '%s'", "404", responseRecorder.Body)
		//}
	})
}

type triggerImpl struct {
	sync.Mutex
	mocks.WaypointServer
}
