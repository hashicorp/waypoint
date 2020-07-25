package core

import (
	"context"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
)

// operation is a private interface that we implement for "operations" such
// as build, deploy, push, etc. This lets us share logic around creating
// server metadata, error checking, etc.
type operation interface {
	// Init returns a new metadata message we'll upsert
	Init(*App) (proto.Message, error)

	// Upsert performs an upsert operation for some metadata
	Upsert(context.Context, pb.WaypointClient, proto.Message) (proto.Message, error)

	// Do performs the actual operation and returns the result that you
	// want to return from the operation. This result will be marshaled into
	// the ValuePtr if it implements ProtoMarshaler.
	// Do can alter the proto.Message into it's final form, as it's the value
	// returned by Init and that will be written back via Upsert after Do
	// has completed.
	Do(context.Context, hclog.Logger, *App, proto.Message) (interface{}, error)

	// StatusPtr and ValuePtr return pointers to the fields in the message
	// for the status and values respectively.
	StatusPtr(proto.Message) **pb.Status
	ValuePtr(proto.Message) **any.Any

	// Hooks are the hooks to execute as part of this operation keyed by "when"
	Hooks(*App) map[string][]*config.Hook
}

func (a *App) doOperation(
	ctx context.Context,
	log hclog.Logger,
	op operation,
) (interface{}, proto.Message, error) {
	// Get our hooks
	hooks := op.Hooks(a)

	// Init the metadata
	msg, err := op.Init(a)
	if err != nil {
		return nil, nil, err
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

	// Reset the status pointer because we might have a new message type
	if ptr := op.StatusPtr(msg); ptr != nil {
		statusPtr = ptr
	}

	// Get where we'll set the value. Similar to statusPtr, we set this
	// to a local value if we get nil so that we can avoid nil checks.
	valuePtr := op.ValuePtr(msg)
	if valuePtr == nil {
		var value *any.Any
		valuePtr = &value
	}

	var doErr error

	// If we have before hooks, run those
	for i, h := range hooks["before"] {
		if err := a.execHook(ctx, log.Named(fmt.Sprintf("hook-before-%d", i)), h); err != nil {
			doErr = err
			log.Warn("error running before hook", "err", err)
		}
	}

	// Run the actual implementation
	var result interface{}
	if doErr == nil {
		log.Debug("running local operation")
		result, doErr = op.Do(ctx, log, a, msg)
		if doErr == nil {
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
		}
	}

	// Run after hooks
	if doErr == nil {
		for i, h := range hooks["after"] {
			if err := a.execHook(ctx, log.Named(fmt.Sprintf("hook-after-%d", i)), h); err != nil {
				doErr = err
				log.Warn("error running after hook", "err", err)
			}
		}
	}

	// If we have an error, then we set the error status
	if doErr != nil {
		log.Warn("error during local operation", "err", doErr)
		*valuePtr = nil
		server.StatusSetError(*statusPtr, doErr)
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

// msgId gets the id of the message by looking for the "Id" field. This
// will return empty string if the ID field can't be found for any reason.
func msgId(msg proto.Message) string {
	val := reflect.ValueOf(msg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Get the Id field
	val = val.FieldByName("Id")
	if !val.IsValid() || val.Kind() != reflect.String {
		return ""
	}

	return val.String()
}
