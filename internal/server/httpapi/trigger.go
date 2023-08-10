// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package httpapi

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/waypoint/internal/clicontext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverclient"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

// Message is the message we return to the requester when streaming job output
type Message struct {
	// The job id that was queued when running the requested trigger
	JobId string `json:"jobId,omitempty"`

	// Value is the job stream event message to stream back to the requester
	Value interface{} `json:"value,omitempty"`
	// ValueType is the kind of job stream event
	ValueType string `json:"valueType,omitempty"`

	// If the job has completed, we return a 0 for success, 1 for failure.
	ExitCode string `json:"exitCode,omitempty"`
	// Error is set when the job failures for any reason
	Error interface{} `json:"error,omitempty"`
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

		variablesJSONRaw := r.URL.Query().Get("variables")
		var (
			variables         map[string]string
			variableOverrides []*pb.Variable
		)

		if variablesJSONRaw != "" {
			if err := json.Unmarshal([]byte(variablesJSONRaw), &variables); err != nil {
				http.Error(w,
					fmt.Sprintf("failed to decode 'variables' json request param into a map: %s", err),
					http.StatusInternalServerError)
				return
			}

			for name, value := range variables {
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
				http.Error(w,
					fmt.Sprintf("request not authorized to run trigger: %s", err),
					http.StatusUnauthorized)
			} else {
				// improve http error code, which is more applicable for general queue failures?
				http.Error(w,
					fmt.Sprintf("server failed to run trigger: %s", err),
					http.StatusPreconditionFailed)
			}

			return
		}
		if resp == nil {
			http.Error(w,
				fmt.Sprintf("server returned no job ids from run trigger %q", html.EscapeString(runTriggerId)),
				http.StatusInternalServerError)
			return
		}
		triggerJobs := resp.JobIds

		streamOutput := r.URL.Query().Get("stream")

		log.Trace("jobs for trigger have been queued")

		if streamOutput == "" {
			// don't stream if not requested
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Attempt to stream output back, on request.
		if !requireAuth {
			// We do not allow streaming job stream info if a no-auth token trigger was requested
			log.Trace("server does not allow for streaming job stream output for no-token trigger URLs")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		log.Debug("attempting to stream back queued job output from running trigger")

		cn, ok := w.(http.CloseNotifier)
		if !ok {
			log.Error("failed to stream job output, could not create http.CloseNotifier")
			http.Error(w, "server failed to create http CloseNotifier", http.StatusInternalServerError)
			return
		}
		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Error("failed to stream job output, could not create http.Flusher")
			http.Error(w, "server failed to create http.Flusher", http.StatusInternalServerError)
			return
		}

		// Send the initial headers saying we're gonna stream the response.
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		enc := json.NewEncoder(w)

		var (
			wg sync.WaitGroup
			mu sync.Mutex
		)

		log.Trace("starting job stream for jobs", "total_jobs", len(triggerJobs))
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
					http.Error(w,
						fmt.Sprintf("server failed to obtain job stream output: %s", err),
						http.StatusInternalServerError)
					return
				}

				log.Trace("reading job stream for job", "job_id", jId)

				// Wait for open confirmation
				resp, err := stream.Recv()
				if err != nil {
					log.Error("server failed to stream job output", "err", err)
					http.Error(w,
						fmt.Sprintf("server failed to receive job stream output: %s", err),
						http.StatusInternalServerError)
					return
				}
				if _, ok := resp.Event.(*pb.GetJobStreamResponse_Open_); !ok {
					log.Error("server failed to open job stream output, got unexpected message", "event", resp.Event)
					http.Error(w,
						fmt.Sprintf("job stream failed to open, got unexpected message: %T", resp.Event),
						http.StatusInternalServerError)
					return
				}

				var (
					jobComplete bool
					exitCode    string
				)

				// read and send the stream
				for {
					select {
					case <-cn.CloseNotify():
						log.Trace("client closed connection to stream")
						return
					default:
						// Get jobstream output and return Message back
						time.Sleep(time.Second)
					}

					resp, err := stream.Recv()
					if err != nil {
						http.Error(w,
							fmt.Sprintf("server failed to receive job stream output: %s", err),
							http.StatusInternalServerError)
						return
					}
					if resp == nil {
						// This shouldn't happen, but if it does, just ignore it.
						log.Warn("nil response received, ignoring")
						continue
					}

					// the message to craft and return
					var (
						value     interface{}
						valueType string
						msgErr    interface{}
					)

					// handle events from job stream
					switch event := resp.Event.(type) {
					case *pb.GetJobStreamResponse_Complete_:
						jobComplete = true
						valueType = "Complete"

						if event.Complete.Error == nil {
							log.Info("job completed successfully")
							exitCode = "0"
						} else {
							exitCode = "1"
							st := status.FromProto(event.Complete.Error)
							log.Warn("job failed", "code", st.Code(), "message", st.Message())
							msgErr = st
						}
					case *pb.GetJobStreamResponse_Error_:
						jobComplete = true
						exitCode = "1"
						valueType = "Error"

						st := status.FromProto(event.Error.Error)
						log.Warn("job stream failure", "code", st.Code(), "message", st.Message())
						msgErr = st
					case *pb.GetJobStreamResponse_Terminal_:
						// We got some job output! Craft a message to be sent back

						for _, ev := range event.Terminal.Events {
							log.Trace("job terminal output", "event", ev)

							switch ev := ev.Event.(type) {
							case *pb.GetJobStreamResponse_Terminal_Event_Line_:
								value = ev.Line.Msg
								valueType = "TerminalEventLine"
							case *pb.GetJobStreamResponse_Terminal_Event_NamedValues_:
								var values []terminal.NamedValue

								for _, tnv := range ev.NamedValues.Values {
									values = append(values, terminal.NamedValue{
										Name:  tnv.Name,
										Value: tnv.Value,
									})
								}

								value = values
								valueType = "TerminalEventNamedValues"
							case *pb.GetJobStreamResponse_Terminal_Event_Status_:
								value = ev.Status.Msg
								valueType = "TerminalEventStatus"
							case *pb.GetJobStreamResponse_Terminal_Event_Raw_:
								value = string(ev.Raw.Data[:])
								valueType = "TerminalEventRaw"
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

								value = tbl
								valueType = "TerminalEventTable"
							case *pb.GetJobStreamResponse_Terminal_Event_Step_:
								m := ev.Step.Msg
								if len(ev.Step.Output) > 0 {
									m = m + "\n" + string(ev.Step.Output[:])
								}
								value = m
								valueType = "TerminalEventStep"
							default:
								log.Error("Unknown terminal event seen", "type", hclog.Fmt("%T", ev))
							}
						}
					default:
						log.Warn("unknown stream event", "event", resp.Event)
					}

					// Send a message job stream back to the client
					if valueType != "" {
						log.Trace("sending job data to client for job", "job_id", jId)

						// Note that all empty values will be omitted
						msg := Message{
							JobId:     jId,
							ExitCode:  exitCode,
							Value:     value,
							ValueType: valueType,
							Error:     msgErr,
						}

						// Lock to ensure multiple routines don't send back a message at the same
						// time to the receiver and mess up the incoming message
						mu.Lock()
						// send the message back
						err := enc.Encode(msg)
						if err != nil {
							log.Error("failed to encode job stream output to send back", "err", err)
							http.Error(w, fmt.Sprintf("server failed to encode job stream output: %s", err), 500)
							mu.Unlock()
							return
						}
						flusher.Flush()
						mu.Unlock()
					}

					if jobComplete {
						log.Trace("job complete, continuing to next job for streaming", "job_id", jId)
						return
					}
				}
			}(jId)
		}

		wg.Wait()
		log.Trace("finished streaming trigger jobs")
	}
}
