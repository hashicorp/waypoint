package httpapi

import (
	"io"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wspb"

	"github.com/hashicorp/waypoint/internal/clicontext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

// HandleExec handles the `waypoint exec` websocket API. This works by
// connecting back to our own local gRPC server.
func HandleExec(addr string, tls bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := hclog.FromContext(ctx)
		log.SetLevel(hclog.Trace)
		log.Info("Websocket exec start", addr)
		// Get our authorization token
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "no token provided", 403)
			return
		}

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

		// Connect back to our own gRPC service.
		grpcConn, err := serverclient.Connect(ctx,
			serverclient.Logger(log),
			serverclient.FromContextConfig(&clicontext.Config{
				Server: serverconfig.Client{
					Address:     addr,
					RequireAuth: true,
					AuthToken:   token,

					// Our gRPC server should always be listening on TLS.
					// We ignore it because its coming out of our own process.
					Tls:           tls,
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

		// Start a goroutine that'll just read requests from our websocket.
		// This will exit once our connection is closed.
		reqCh := make(chan *pb.ExecStreamRequest)
		go func() {
			for {
				var req pb.ExecStreamRequest
				err := wspb.Read(ctx, c, &req)
				if err != nil {
					if err != io.EOF {
						log.Error("websocket receive error", "err", err)
					}
					return
				}

				select {
				case reqCh <- &req:
				case <-ctx.Done():
					log.Warn("context canceled while waiting to send data")
					return
				}
			}
		}()

		// Start goroutine that'll just read exec stream responses. This
		// will exit once our connection is closed.
		respCh := make(chan *pb.ExecStreamResponse)
		go func() {
			for {
				resp, err := exec.Recv()
				if err != nil {
					if err != io.EOF {
						log.Error("stream receive error", "err", err)
					}

					return
				}

				// Send our received data but exit if our context is canceled.
				select {
				case respCh <- resp:
					log.Info("resp data sent to channel")
				case <-ctx.Done():
					log.Warn("context canceled while waiting to send data")
					return
				}
			}
		}()

		// Sit in an event loop and just shuttle data back and forth.
		for {
			select {
			case <-ctx.Done():
				// Done!
				c.Close(websocket.StatusNormalClosure, "")
				return

			case req := <-reqCh:
				if err := exec.Send(req); err != nil {
					c.Close(websocket.StatusInternalError, err.Error())
					return
				}

			case resp := <-respCh:
				if err := wspb.Write(ctx, c, resp); err != nil {
					c.Close(websocket.StatusInternalError, err.Error())
					return
				}
			}
		}
	}
}
