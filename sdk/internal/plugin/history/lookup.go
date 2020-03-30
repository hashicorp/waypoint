package history

import (
	"github.com/mitchellh/devflow/sdk/history"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

// LookupToProto converts a history.Lookup to the proto value.
func LookupToProto(req *history.Lookup) *pb.History_LookupRequest {
	if req == nil {
		req = &history.Lookup{}
	}

	return &pb.History_LookupRequest{
		Limit:        int32(req.Limit),
		FilterStatus: pb.History_LookupRequest_FilterStatus(req.FilterStatus),
	}
}

// LookupFromProto converts a history.Lookup to the proto value.
func LookupFromProto(req *pb.History_LookupRequest) *history.Lookup {
	return &history.Lookup{
		Limit:        int(req.Limit),
		FilterStatus: history.FilterStatus(req.FilterStatus),
	}
}
