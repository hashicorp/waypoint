package tfc

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

func TestConfigSourcer(t *testing.T) {
	ctx := context.Background()
	log := hclog.L()
	require := require.New(t)

	var (
		includes string
	)

	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Logf("request: %s", req.URL.String())

		switch req.URL.Path {
		case "/api/v2/organizations/foocorp/workspaces/databases":
			rw.Write([]byte(`{
  "data": {
    "id": "ws-abcdef",
    "type": "workspaces"
	}
}
`))
		case "/api/v2/workspaces/ws-abcdef/current-state-version":
			includes = req.URL.Query().Get("include")

			data, err := ioutil.ReadFile("testdata/state-version.json")
			if err != nil {
				panic(err)
			}

			rw.Write(data)
		}

	}))

	defer s.Close()

	cs := &ConfigSourcer{}
	cs.config.Token = "xxyyzz"
	cs.config.BaseURL = s.URL

	t.Logf("local server: %s", s.URL)

	defer cs.stop()

	// Read
	result, err := cs.read(ctx, log, []*component.ConfigRequest{
		{
			Name: "HOSTNAME",
			Config: map[string]string{
				"organization": "foocorp",
				"workspace":    "databases",
				"output":       "aws_assume_role_arn",
			},
		},
	})
	require.NoError(err)
	require.NotNil(result)
	require.Len(result, 1)

	v := result[0]

	require.Equal("outputs", includes)

	require.Equal("HOSTNAME", v.Name)
	ve, ok := v.Result.(*pb.ConfigSource_Value_Error)
	if ok {
		require.False(ok, "error reading value: %s", ve.Error.Message)
	}

	require.Equal(
		"arn:aws:iam::797645259670:role/waypoint-circleci-release-pipeline",
		v.Result.(*pb.ConfigSource_Value_Value).Value,
	)
}
