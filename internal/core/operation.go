package core

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"

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
	Do(context.Context, hclog.Logger, *App) (interface{}, error)

	// StatusPtr and ValuePtr return pointers to the fields in the message
	// for the status and values respectively.
	StatusPtr(proto.Message) **pb.Status
	ValuePtr(proto.Message) **any.Any
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

	// Update the status
	statusPtr := op.StatusPtr(msg)
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
	statusPtr = op.StatusPtr(msg)
	valuePtr := op.ValuePtr(msg)

	// Run the function
	log.Debug("running local operation")
	result, doErr := op.Do(ctx, log, a)
	if doErr == nil {
		// No error, our state is success
		server.StatusSetSuccess(*statusPtr)

		// Set our final value
		*valuePtr, err = component.ProtoAny(result)
		if err != nil {
			doErr = err
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
