package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/serverclient"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

// Message is the message we return to the requester when streaming output
type Message struct {
	JobId    string `json:"jobId,omitempty"`
	Message  string `json:"message,omitempty"`
	ExitCode string `json:"exitCode,omitempty"`
}

// HandleTrigger will execute a run trigger, if the requested id exists
// This works by connecting back to our own local gRPC server.
func HandleTrigger(addr string, tls bool) http.HandlerFunc {
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
					Tls:           tls,
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

		// attempt to make a grpc request to run trigger by id
		if requireAuth {
			resp, err = client.RunTrigger(ctx, runTriggerReq)
		} else {
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

		streamOutput := r.URL.Query().Get("stream")

		// Attempt to stream output back, on request.
		if streamOutput != "" {
			if !requireAuth {
				// We do not allow streaming job stream info if a no-auth token trigger was requested
				log.Debug("server does not allow for streaming job stream output for no-token trigger URLs")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			log.Trace("attempting to stream back queued job output from running trigger")

			cn, ok := w.(http.CloseNotifier)
			if !ok {
				log.Error("failed to stream job output, could not create http.CloseNotifier")
				http.NotFound(w, r)
				return
			}
			flusher, ok := w.(http.Flusher)
			if !ok {
				log.Error("failed to stream job output, could not create http.Flusher")
				http.NotFound(w, r)
				return
			}

			// Send the initial headers saying we're gonna stream the response.
			w.Header().Set("Transfer-Encoding", "chunked")
			w.WriteHeader(http.StatusOK)
			flusher.Flush()

			enc := json.NewEncoder(w)

			// NOTE(briancain): We skip every two jobs here because when we call RunTrigger
			// via gRPC, it eventually queues the trigger jobs through on-demand runners, and
			// queueJobMulti returns three jobs: StartTask, the job to be queued, and StopTask. People
			// really only expect output from the job to be queued, so we only stream that back.

			// For example, a trigger that queues 2 Waypoint operation jobs returns 6 total job ids:
			// Job List: [ {0: StartTask, 1: WP Operation 1, 2: StopTask}, {3: StartTask, 4: WP Operation 2, 5: StopTask}, ... ]

			var triggerJobs []string
			for i := 1; i < len(jobIds); i += 3 {
				triggerJobs = append(triggerJobs, jobIds[i])
			}

			var (
				wg sync.WaitGroup
				mu sync.Mutex
			)

			wg.Add(len(triggerJobs))

			// NOTE(briancain): This loop starts N goroutines concurrently for
			// each trigger job to stream back to the requester.
			for _, jId := range triggerJobs {
				go func(jId string) {
					defer wg.Done()

					stream, err := client.GetJobStream(ctx, &pb.GetJobStreamRequest{
						JobId: jId,
					})
					if err != nil {
						log.Error("server failed to get job stream output for trigger", "job_id", jId, "err", err)
						http.Error(w, fmt.Sprintf("server failed to obtain job stream output: %s", err), 500)
						return
					}

					// Wait for open confirmation
					resp, err := stream.Recv()
					if err != nil {
						log.Error("server failed to stream job output", "err", err)
						http.Error(w, fmt.Sprintf("server failed to receive job stream output: %s", err), 500)
						return
					}
					if _, ok := resp.Event.(*pb.GetJobStreamResponse_Open_); !ok {
						log.Error("server failed to open job stream output, got unexpected message", "event", resp.Event)
						http.Error(w, fmt.Sprintf("job stream failed to open, got unexpected message: %T", resp.Event), 500)
						return
					}

					var (
						jobComplete bool
						exitCode    string
					)

					// read and send the stream
					for {
						resp, err := stream.Recv()
						if err != nil {
							http.Error(w, fmt.Sprintf("server failed to receive job stream output: %s", err), 500)
							return
						}
						if resp == nil {
							// This shouldn't happen, but if it does, just ignore it.
							log.Warn("nil response received, ignoring")
							continue
						}

						select {
						case <-cn.CloseNotify():
							log.Trace("client closed connection to stream")
							return
						default:
							// Get jobstream output and return string chunk message back
							time.Sleep(time.Second)
							m := "" // the message to return

							switch event := resp.Event.(type) {
							case *pb.GetJobStreamResponse_Complete_:
								jobComplete = true
								m = "job complete"

								if event.Complete.Error == nil {
									log.Info("job completed successfully")
									exitCode = "0"
								} else {
									exitCode = "1"
									st := status.FromProto(event.Complete.Error)
									log.Warn("job failed", "code", st.Code(), "message", st.Message())
									http.Error(w, fmt.Sprintf("job failed to complete: job code %s: %s", st.Code(), st.Message()), 500)
								}
							case *pb.GetJobStreamResponse_Error_:
								jobComplete = true
								exitCode = "1"

								st := status.FromProto(event.Error.Error)
								log.Warn("job stream failure", "code", st.Code(), "message", st.Message())
								http.Error(w, fmt.Sprintf("job failed to complete: job code %s: %s", st.Code(), st.Message()), 500)
								return
							case *pb.GetJobStreamResponse_Terminal_:
								// We got some job output! Craft a message to be sent back

								for _, ev := range event.Terminal.Events {
									log.Trace("job terminal output", "event", ev)

									switch ev := ev.Event.(type) {
									case *pb.GetJobStreamResponse_Terminal_Event_Line_:
										m = ev.Line.Msg
									case *pb.GetJobStreamResponse_Terminal_Event_NamedValues_:
										var values []terminal.NamedValue

										for _, tnv := range ev.NamedValues.Values {
											values = append(values, terminal.NamedValue{
												Name:  tnv.Name,
												Value: tnv.Value,
											})
										}

										jsonNamedValues, err := json.Marshal(values)
										if err != nil {
											log.Warn("job stream failed to marshal NamedValues to json", "err", err)
											http.Error(w, fmt.Sprintf("job failed to marshal NamedValues to json: %s", err), 500)
											return
										}

										m = string(jsonNamedValues)
									case *pb.GetJobStreamResponse_Terminal_Event_Status_:
										// Since we're not writing to a terminal that can update in place
										// we ignore steps and send the message directly instead
										m = ev.Status.Msg
									case *pb.GetJobStreamResponse_Terminal_Event_Raw_:
										// does message need to preserve stdout/stderr??
										m = string(ev.Raw.Data[:])
									case *pb.GetJobStreamResponse_Terminal_Event_Table_:
										tbl := terminal.NewTable(ev.Table.Headers...)

										for _, row := range ev.Table.Rows {
											var trow []terminal.TableEntry

											for _, ent := range row.Entries {
												trow = append(trow, terminal.TableEntry{
													Value: ent.Value,
													Color: ent.Color,
												})
											}
										}

										jsonTable, err := json.Marshal(tbl)
										if err != nil {
											log.Warn("job stream failed to marshal Table Event to json", "err", err)
											http.Error(w, fmt.Sprintf("job failed to marshal Table Event to json: %s", err), 500)
											return
										}

										m = string(jsonTable)
									case *pb.GetJobStreamResponse_Terminal_Event_Step_:
										m = ev.Step.Msg
										if len(ev.Step.Output) > 0 {
											m = m + "\n" + string(ev.Step.Output[:])
										}
									default:
										log.Error("Unknown terminal event seen", "type", hclog.Fmt("%T", ev))
									}
								}
							default:
								log.Warn("unknown stream event", "event", resp.Event)
							}

							// Send a message job stream back to the client
							if m != "" {
								log.Trace("sending job data to client for job", "job_id", jId)
								msg := Message{
									JobId:    jId,
									Message:  m,
									ExitCode: exitCode, // will only be sent if set to non-empty string
								}

								// Lock to ensure multiple routines don't send back a message at the same
								// time to the receiver and mess up the incoming message
								mu.Lock()
								// send the message back
								err := enc.Encode(msg)
								if err != nil {
									log.Error("failed to encode job stream output to send back", "err", err)
									http.Error(w, fmt.Sprintf("server failed to encode job stream output: %s", err), 500)
									return
								}
								flusher.Flush()
								mu.Unlock()
							}

							if jobComplete {
								log.Trace("job complete, continuing to next job for streaming")
								return
							}
						}
					}
				}(jId)
			}

			wg.Wait()
			log.Trace("finished streaming trigger jobs")
		}
	}
}
