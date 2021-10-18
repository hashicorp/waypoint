package httpapi

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
	"nhooyr.io/websocket"

	"github.com/hashicorp/waypoint/internal/clicontext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

// HandleExec handles the `waypoint exec` websocket API. This works by
// connecting back to our own local gRPC server.
func HandleExec(addr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := hclog.FromContext(ctx)

		// Accept our websocket connection.
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			// We allow connections from anywhere so that the UI can be on another
			// host and so that we work with proxying. We aren't too worried about
			// CSRF due to the nature of Waypoint servers and the underlying
			// protocols and auth tokens, but one day it'd be nice to be stricter
			// here and accept only known UI hosts (we don't even have this
			// config at the time of writing).
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Info("error accepting websocket connection", "err", err)
			return
		}

		// We can defer close an error because multiple calls to Close are ignored.
		defer c.Close(websocket.StatusInternalError, "early exit")

		// Before we spend a bunch of time connecting back to the gRPC
		// service and all that, let's read the start request from
		// the websocket. This is used to initialize the exec stream. If the
		// client never sends it, then we can just exit out without wasting
		// resources even attempting to connect back to gRPC.
		// TODO

		// Connect back to our own gRPC service.
		grpcConn, err := serverclient.Connect(ctx,
			serverclient.Logger(log),
			serverclient.FromContextConfig(&clicontext.Config{
				Server: serverconfig.Client{
					Address:     addr,
					RequireAuth: true,
					AuthToken:   "TODO",

					// Our gRPC server should always be listening on TLS.
					// We ignore it because its coming out of our own process.
					Tls:           true,
					TlsSkipVerify: true,
				},
			}),
		)
		if err != nil {
			log.Error("exec connection back to gRPC failed", "err", err)
			c.Close(
				websocket.StatusInternalError,
				"failed to connect to internal server",
			)
			return
		}
		defer grpcConn.Close()

		// Our API client
		client := pb.NewWaypointClient(grpcConn)

		// Start our exec stream
		exec, err := client.StartExecStream(ctx)
		if err != nil {
			c.Close(
				websocket.StatusInternalError,
				err.Error(),
			)
			return
		}
		defer exec.CloseSend()
	}
}
