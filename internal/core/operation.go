// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/opaqueany"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/component"

	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/pkg/finalcontext"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// operation is a private interface that we implement for "operations" such
// as build, deploy, push, etc. This lets us share logic around creating
// server metadata, error checking, etc.
type operation interface {

	// Name returns the name of the operation for display purposes
	Name() string

	// Init returns a new metadata message we'll upsert. This is the first
	// function called before any other operation logic is executed and so
	// can be used to initialize state for the other callbacks.
	Init(*App) (proto.Message, error)

	// Upsert performs an upsert operation for some metadata
	Upsert(context.Context, pb.WaypointClient, proto.Message) (proto.Message, error)

	// Do performs the actual operation and returns the result that you
	// want to return from the operation. This result will be marshaled into
	// the ValuePtr if it implements ProtoMarshaler.
	// Do can alter the proto.Message into its final form, as it's the value
	// returned by Init and that will be written back via Upsert after Do
	// has completed.
	Do(context.Context, hclog.Logger, *App, proto.Message) (interface{}, error)

	// StatusPtr and ValuePtr return pointers to the fields in the message
	// for the status and values respectively. For ValuePtr, it returns both
	// a pointer to store the raw message as well as a pointer to store the JSON
	// of the message. The JSON pointer can be nil and it won't be stored.
	StatusPtr(proto.Message) **pb.Status
	ValuePtr(proto.Message) (**opaqueany.Any, *string)

	// Hooks are the hooks to execute as part of this operation keyed by "when"
	Hooks(*App) map[string][]*config.Hook

	// Labels is called to return any labels that should be set for this
	// operation. This should include the component labels. These will be merged
	// with any resulting labels from the operation.
	Labels(*App) map[string]string
}

func (a *App) doOperation(
	ctx context.Context,
	log hclog.Logger,
	op operation,
) (interface{}, proto.Message, error) {
	// Init the metadata
	msg, err := op.Init(a)
	if err != nil {
		return nil, nil, err
	}

	// Get our hooks
	hooks := op.Hooks(a)

	// Initialize our labels
	msgUpdateLabels(a, op.Labels(a), msg, nil)

	// Setup our job id if we have that field.
	if f := msgField(msg, "JobId"); f.IsValid() {
		f.Set(reflect.ValueOf(a.jobInfo.Id))
	}

	// If we have no status pointer, then we just allocate one for this
	// function. We don't send this anywhere but this just lets us follow
	// the remaining logic without a bunch of nil checks.
	statusPtr := op.StatusPtr(msg)
	if statusPtr == nil {
		var status *pb.Status
		statusPtr = &status
	}
	*statusPtr = server.NewStatus(pb.Status_RUNNING)

	// Upsert the metadata for our running state
	log.Debug("creating metadata on server")
	msg, err = op.Upsert(ctx, a.client, msg)
	if err != nil {
		return nil, nil, err
	}
	if id := msgId(msg); id != "" {
		log = log.With("id", id)
	}

	if sequence, err := msgSequence(msg); err == nil {
		sg := a.UI.StepGroup()
		s := sg.Add("Running %s v%d", op.Name(), sequence)
		defer s.Done()
	} else {
		// Not all operation upsert responses have sequence ids (i.e. status report).
		// For now, just ignore in that case.
	}

	// Reset the status pointer because we might have a new message type
	if ptr := op.StatusPtr(msg); ptr != nil {
		statusPtr = ptr
	}

	// Get where we'll set the value. Similar to statusPtr, we set this
	// to a local value if we get nil so that we can avoid nil checks.
	valuePtr, valueJsonPtr := op.ValuePtr(msg)
	if valuePtr == nil {
		var value *opaqueany.Any
		valuePtr = &value
	}
	if valueJsonPtr == nil {
		var valueJson string
		valueJsonPtr = &valueJson
	}

	var doErr error

	// If we have before hooks, run those
	for i, h := range hooks["before"] {
		if err := a.execHook(ctx, log.Named(fmt.Sprintf("hook-before-%d", i)), h); err != nil {
			doErr = fmt.Errorf("Error running before hook index %d: %w", i, err)
			log.Warn("error running before hook", "err", err)

			if h.ContinueOnFailure() {
				log.Info("hook configured to continue on failure, ignoring error")
				doErr = nil
			}
		}
	}

	// Run the actual implementation
	var result interface{}
	if doErr == nil {
		log.Debug("running local operation")
		result, doErr = op.Do(ctx, log, a, msg)

		// If we got an argmapper error, we transfor it into something a
		// bit more readable for the average Waypoint user.
		if err, ok := doErr.(*argmapper.ErrArgumentUnsatisfied); ok && err != nil {
			// Log the full error
			log.Warn("argmapper unsatified error received", "error", doErr)

			// Build our list of missing arguments
			missing := new(bytes.Buffer)
			for _, v := range err.Args {
				s := v.String()

				// If this is an any type with a subtype, then we just use
				// the subtype. The subtype wording confuses people so just
				// note the direct type.
				if v.Type == anyType && v.Subtype != "" {
					s = v.Subtype
				}

				fmt.Fprintf(missing, "    - %s\n", s)
			}

			doErr = fmt.Errorf(
				"There was an error while executing a Waypoint plugin for "+
					"this operation!\n\n"+
					"One or more required arguments for the plugin was not satisfied. "+
					"This is usually due to a missing or incompatible set of plugins. "+
					"For example, only certain build plugins are only compatible with certain "+
					"registries, and so on. Please inspect the missing argument, the set of "+
					"plugins you are using, and the documentation to determine if your "+
					"plugin combination is valid.\n\n"+
					"Plugin function: %s\n\n"+
					"==> Missing arguments:\n\n%s",
				err.Func.Name(),
				missing.String(),
			)
		}

		if doErr == nil {
			// Set our labels if we can
			msgUpdateLabels(a, op.Labels(a), msg, result)

			// Set the deployment URL, if possible
			msgUpdateURL(msg, result)

			// Set our template data. Any errors here are logged but ignored
			// since we don't want to leave dangling physical resources.
			if err := msgUpdateTemplateData(msg, result); err != nil {
				log.Warn("error encoding template data, will not be stored", "err", err)
			}

			// No error, our state is success
			server.StatusSetSuccess(*statusPtr)

			// Set our final value if we have a value pointer
			*valuePtr = nil
			if result != nil {
				*valuePtr, err = component.ProtoAny(result)
				if err != nil {
					doErr = err
				}
			}

			// If we can marshal as JSON, do it
			*valueJsonPtr = ""
			if m, ok := result.(json.Marshaler); ok {
				raw, err := m.MarshalJSON()
				if err != nil {
					doErr = err
				} else {
					*valueJsonPtr = string(raw)
				}
			}
		}
	}

	// Run after hooks
	if doErr == nil {
		for i, h := range hooks["after"] {
			if err := a.execHook(ctx, log.Named(fmt.Sprintf("hook-after-%d", i)), h); err != nil {
				doErr = fmt.Errorf("Error running after hook index %d: %w", i, err)
				log.Warn("error running after hook", "err", err)

				if h.ContinueOnFailure() {
					log.Info("hook configured to continue on failure, ignoring error")
					doErr = nil
				}
			}
		}
	}

	// If we have an error, then we set the error status
	if doErr != nil {
		log.Warn("error during local operation", "err", doErr)
		*valuePtr = nil
		server.StatusSetError(*statusPtr, doErr)
	}

	// If our context ended we need to create a final context so we
	// can attempt to finalize our metadata.
	if ctx.Err() != nil {
		var cancel context.CancelFunc
		ctx, cancel = finalcontext.Context(log)
		defer cancel()
	}

	// Set the final metadata
	msg, err = op.Upsert(ctx, a.client, msg)
	if err != nil {
		log.Warn("error marking server metadata as complete", "err", err)
	} else {
		log.Debug("metadata marked as complete")
	}

	// If we had an original error, return it now that we have saved all metadata
	if doErr != nil {
		return nil, nil, doErr
	}

	return result, msg, nil
}

func msgUpdateURL(msg proto.Message, result interface{}) {
	val := msgField(msg, "Url")
	if !val.IsValid() {
		return
	}

	switch t := result.(type) {
	case component.DeploymentWithUrl:
		val.SetString(t.URL())
	}
}

func msgUpdateLabels(
	app *App,
	base map[string]string,
	msg proto.Message,
	result interface{},
) {
	// Get our labels field in our proto message. If we don't have one
	// then we don't bother doing anything else since labels are moot.
	val := msgField(msg, "Labels")
	if !val.IsValid() {
		return
	}

	// Determine any labels we have in our result
	var resultLabels map[string]string
	if labels, ok := result.(interface{ Labels() map[string]string }); ok {
		resultLabels = labels.Labels()
	}

	// Merge them
	val.Set(reflect.ValueOf(app.mergeLabels(base, resultLabels)))
}

func msgUpdateTemplateData(
	msg proto.Message,
	result interface{},
) error {
	// Get our template data field in our proto message. If we don't have one
	// then we don't bother doing anything.
	val := msgField(msg, "TemplateData")
	if !val.IsValid() {
		return nil
	}

	// Determine if we have template data
	tpl, ok := result.(component.Template)
	if !ok {
		return nil
	}

	// Marshal it
	tplData, err := json.Marshal(tpl.TemplateData())
	if err != nil {
		return err
	}

	// Merge them
	val.Set(reflect.ValueOf(tplData))
	return nil
}

// msgId gets the id of the message by looking for the "Id" field. This
// will return empty string if the ID field can't be found for any reason.
func msgId(msg proto.Message) string {
	val := msgField(msg, "Id")
	if !val.IsValid() || val.Kind() != reflect.String {
		return ""
	}

	return val.String()
}

// msgSequence gets the sequence number of the message by looking for the
// "Sequence" field. This will return an error if the Sequence field
// can't be found for any reason.
func msgSequence(msg proto.Message) (uint64, error) {
	val := msgField(msg, "Sequence")
	if !val.IsValid() || val.Kind() != reflect.Uint64 {
		return 0, fmt.Errorf("Sequence field not found")
	}

	return val.Uint(), nil
}

// msgField gets the field from the given message. This will return an
// invalid value if it doesn't exist.
func msgField(msg proto.Message, f string) reflect.Value {
	val := reflect.ValueOf(msg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// If the value is invalid then we don't do anything. This could be because
	// msg is nil.
	if !val.IsValid() {
		return val
	}

	// Get the Id field
	return val.FieldByName(f)
}

// anyType is used to compare types.
var anyType = reflect.TypeOf((*opaqueany.Any)(nil))
