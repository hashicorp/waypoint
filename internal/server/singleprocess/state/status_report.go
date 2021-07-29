package state

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var statusReportOp = &appOperation{
	Struct: (*pb.StatusReport)(nil),
	Bucket: []byte("statusreport"),

	// This number is global, not per deployment. So we set this number to a high
	// number instead of trying to store just "one" per deploy/release
	MaximumIndexedRecords: 10000,
}

func init() {
	statusReportOp.register()
}

// get status report by referenced operation
func (s *State) StatusReportGet(ref *pb.Ref_Operation) (*pb.StatusReport, error) {
	result, err := statusReportOp.Get(s, ref)
	if err != nil {
		return nil, err
	}

	return result.(*pb.StatusReport), nil
}

// StatusReportPut creates or updates the latest status report
func (s *State) StatusReportPut(update bool, report *pb.StatusReport) error {
	return statusReportOp.Put(s, update, report)
}

func protoTimestampToTime(timestamp *timestamppb.Timestamp) time.Time {
	return time.Unix(timestamp.Seconds, int64(timestamp.Nanos))
}

func (s *State) StatusReportList(
	ref *pb.Ref_Application,
	opts ...ListOperationOption,
) ([]*pb.StatusReport, error) {
	raw, err := statusReportOp.List(s, buildListOperationsOptions(ref, opts...))
	if err != nil {
		return nil, err
	}

	resMap := map[string]*pb.StatusReport{}

	var result []*pb.StatusReport
	for _, v := range raw {
		sr := v.(*pb.StatusReport)
		theId := ""
		depId := sr.GetDeploymentId()
		resId := sr.GetReleaseId()

		if depId != "" {
			theId = depId
		} else {
			theId = resId
		}

		if sr.Status.CompleteTime == nil {
			continue
		}

		if mVal, ok := resMap[theId]; ok {
			srTime := protoTimestampToTime(sr.Status.CompleteTime)
			mValTime := protoTimestampToTime(mVal.Status.CompleteTime)
			if srTime.After(mValTime) {
				resMap[theId] = sr
			}
		} else {
			resMap[theId] = sr
		}
	}

	for _, v := range resMap {
		result = append(result, v)
	}

	return result, nil
}

// StatusReportLatest gets the latest generated status report
func (s *State) StatusReportLatest(
	ref *pb.Ref_Application,
	ws *pb.Ref_Workspace,
) (*pb.StatusReport, error) {
	result, err := statusReportOp.Latest(s, ref, ws)
	if result == nil || err != nil {
		return nil, err
	}

	return result.(*pb.StatusReport), nil
}
