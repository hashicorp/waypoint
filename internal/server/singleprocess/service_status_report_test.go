package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestServiceStatusReport(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	type Req = pb.UpsertStatusReportRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
			StatusReport: serverptypes.TestValidStatusReport(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.StatusReport
		require.NotEmpty(result.Id)

		// Let's write some data
		result.Status = server.NewStatus(pb.Status_RUNNING)
		resp, err = client.UpsertStatusReport(ctx, &Req{
			StatusReport: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.StatusReport
		require.NotNil(result.Status)
		require.Equal(pb.Status_RUNNING, result.Status.State)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertStatusReport(ctx, &Req{
			StatusReport: serverptypes.TestValidStatusReport(t, &pb.StatusReport{
				Id: "nope",
			}),
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}

func TestServiceStatusReport_GetStatusReport(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	statusReportResp, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
		StatusReport: serverptypes.TestValidStatusReport(t, nil),
	})
	require.NoError(t, err)

	type Req = pb.GetStatusReportRequest

	t.Run("get existing", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a status report
		sp, err := client.GetStatusReport(ctx, &Req{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{Id: statusReportResp.StatusReport.Id},
			},
		})
		require.NoError(err)
		require.NotNil(sp)
		require.NotEmpty(sp.Id)
	})

	t.Run("get non-existing", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetStatusReport(ctx, &Req{
			Ref: &pb.Ref_Operation{
				Target: &pb.Ref_Operation_Id{Id: "nope"},
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}

func TestServiceStatusReport_ListStatusReports(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	resp, err := client.UpsertStatusReport(ctx, &pb.UpsertStatusReportRequest{
		StatusReport: serverptypes.TestValidStatusReport(t, nil),
	})
	require.NoError(t, err)

	type Req = pb.ListStatusReportsRequest

	t.Run("list", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a status report
		sr, err := client.ListStatusReports(ctx, &Req{
			Application: resp.StatusReport.Application,
		})
		require.NoError(err)
		require.NotEmpty(sr)
		require.Equal(sr.StatusReports[0].Id, resp.StatusReport.Id)
	})
}
