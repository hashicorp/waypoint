package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

		// Authless trigger URLs should be able to make a request
		// without a token.
		token := r.URL.Query().Get("token")
		requireAuth := true
		if token == "" {
			log.Trace("no token provided, will attempt to run authless trigger")
			requireAuth = false
			//http.Error(w, "no token provided", 403)
		}

		// Connect back to our own gRPC service.
		grpcConn, err := serverclient.Connect(ctx,
			serverclient.Logger(log),
			serverclient.FromContextConfig(&clicontext.Config{
				Server: serverconfig.Client{
					Address:     addr,
					RequireAuth: requireAuth,
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

		overrideVarJSONRaw := r.URL.Query().Get("variables")
		var (
			vo                map[string]string
			variableOverrides []*pb.Variable
		)

		if overrideVarJSONRaw != "" {
			if err := json.Unmarshal([]byte(overrideVarJSONRaw), &vo); err != nil {
				http.Error(w, fmt.Sprintf("failed to decode 'variables' json request param into a map: %s", err), 500)
				return
			}

			for name, value := range vo {
				v := &pb.Variable{
					Name:   name,
					Source: &pb.Variable_Cli{Cli: &empty.Empty{}},
				}

				if valBool, err := strconv.ParseBool(value); err == nil {
					v.Value = &pb.Variable_Bool{Bool: valBool}
				} else if valInt, err := strconv.ParseInt(value, 10, 64); err == nil {
					v.Value = &pb.Variable_Num{Num: valInt}
				} else {
					// NOTE: for this case, it can either be a "string" or
					// complex HCL type like an array or map. We can set this value
					// as a Variable_Str here, and when we go to parse the variables
					// later we do the proper string versus HCL check in variables.go
					v.Value = &pb.Variable_Str{Str: value}
				}

				variableOverrides = append(variableOverrides, v)
			}
		}

		var resp *pb.RunTriggerResponse
		runTriggerReq := &pb.RunTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: runTriggerId,
			},
			VariableOverrides: variableOverrides,
		}

		if requireAuth {
			// attempt to make a grpc request to run trigger by id
			resp, err = client.RunTrigger(ctx, runTriggerReq)
		} else {
			// attempt to make a grpc request to run trigger by id
			resp, err = client.NoAuthRunTrigger(ctx, runTriggerReq)
		}
		if err != nil {
			log.Error("server failed to run trigger", "id", runTriggerId, "err", err)

			if status.Code(err) == codes.PermissionDenied {
				http.Error(w, fmt.Sprintf("request not authorized to run trigger: %s", err), 401)
			} else {
				// improve http error code, which is more applicable for general queue failures?
				http.Error(w, fmt.Sprintf("server failed to run trigger: %s", err), 412)
			}

			return
		}
		if resp == nil {
			http.Error(w, fmt.Sprintf("server returned no job ids from run trigger %q", runTriggerId), 500)
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
