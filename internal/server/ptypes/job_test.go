package ptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestValidateJob(t *testing.T) {
	cases := []struct {
		Name   string
		Modify func(*pb.Job)
		Error  string
	}{
		{
			"valid",
			nil,
			"",
		},

		{
			"id is set",
			func(j *pb.Job) { j.Id = "nope" },
			"id: must be empty",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			job := TestJobNew(t, nil)
			if f := tt.Modify; f != nil {
				f(job)
			}

			err := ValidateJob(job)
			if tt.Error == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Error)
		})
	}
}
