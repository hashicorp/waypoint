package httpapi

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/clicontext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

// TODO(briancain): write tests
// HandleTrigger will execute a run trigger, if the requested id exists
// This works by connecting back to our own local gRPC server.
func HandleTrigger(addr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := hclog.FromContext(ctx)
		log.SetLevel(hclog.Debug)

		// TODO(briancain): Authless trigger URLs should be able to make a request
		// without a token.
		// "No token" requests should probably return http 404 on authenticated trigger URLs
		// Get our authorization token
		token := r.URL.Query().Get("token")
		if token == "" {
			// TODO(briancain): Not an error yet, look up trigger by id and if not
			// authenticated, then continue
			http.Error(w, "no token provided", 403)
			return
		}

		// TODO(briancain): handle auth with "token user" tokens for authenticated trigger URLs
		// Connect back to our own gRPC service.
		// TODO(briancain): how do authless requests initiate a grpc connection
		grpcConn, err := serverclient.Connect(ctx,
			serverclient.Logger(log),
			serverclient.FromContextConfig(&clicontext.Config{
				Server: serverconfig.Client{
					Address:     addr,
					RequireAuth: true,
					AuthToken:   token,

					// Our gRPC server should always be listening on TLS.
					// We ignore it because its coming out of our own process.
					Tls:           true,
					TlsSkipVerify: true,
				},
			}),
		)
		if err != nil {
			log.Error("trigger connection back to gRPC failed", "err", err)
			return
		}
		defer grpcConn.Close()

		// Our API client
		client := pb.NewWaypointClient(grpcConn)

		requestVars := mux.Vars(r)
		runTriggerId := requestVars["id"]
		//triggerOverrideVars := requestVars["override_vars"]

		// attempt to make a grpc request to run trigger by id
		resp, err := client.RunTrigger(ctx, &pb.RunTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: runTriggerId,
			},
			VariableOverrides: nil, // TODO fix me
		})
		if err != nil {
			log.Error("server failed to run trigger", "id", runTriggerId, "err", err)
			// improve http error code, which is more applicable for general queue failures?
			http.Error(w, fmt.Sprintf("server failed to run trigger: %s", err), 412)
			return
		}
		jobIds := resp.JobIds

		// TODO(briancain): attempt to stream output back, on request.
		for _, jId := range jobIds {
			_, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{
				JobId: jId,
			})
			if err != nil {
				log.Error("server failed to get job stream output for trigger", "job_id", jId, "err", err)
			}
		}
	}
}
