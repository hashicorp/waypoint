package tfc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

// Tests that the TFC configsourcer can return a simple string value
func TestConfigSourcer_StringValue(t *testing.T) {
	ctx := context.Background()
	log := hclog.L()
	require := require.New(t)

	var (
		includes                string
		tfcRespJsonTestdataPath string
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

			data, err := ioutil.ReadFile(tfcRespJsonTestdataPath)
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

	t.Run("can read single string value", func(t *testing.T) {
		tfcRespJsonTestdataPath = "testdata/state-version.json"

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
	})

	t.Run("can read single complex-typed value", func(t *testing.T) {
		tfcRespJsonTestdataPath = "testdata/state-version-complex.json"

		result, err := cs.read(ctx, log, []*component.ConfigRequest{
			{
				Name: "HOSTNAME",
				Config: map[string]string{
					"organization": "foocorp",
					"workspace":    "databases",
					"output":       "ecs_task_subnets",
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

		resultJson, ok := v.Result.(*pb.ConfigSource_Value_Json)
		require.True(ok)

		var resultMarshalled []string
		require.NoError(json.Unmarshal(resultJson.Json, &resultMarshalled))

		require.True(ok)
		require.ElementsMatch([]string{
			"subnet-03afefcc38a919083",
			"subnet-087ea2efd4f009fc5",
			"subnet-09e02553a59746ba1",
			"subnet-05a2d1112fbca071c",
			"subnet-0d5f650b4d0eebc56",
		}, resultMarshalled)
	})

	t.Run("Read all outputs", func(t *testing.T) {

		cases := []struct {
			tfcRespJsonTestdataPath string
			numOutputs              int
			spotChecks              map[string]interface{}
		}{
			{
				"testdata/state-version.json",
				3,
				map[string]interface{}{
					"alb_listener_arn": "arn:aws:elasticloadbalancing:us-east-1:797645259670:listener/app/acmeapp1-dev/0ed92920e20ed1dc/07e51901e3cec498",
				},
			},
			{
				"testdata/state-version-complex.json",
				14,
				map[string]interface{}{
					"region": "us-east-1",
					"ecs_task_subnets": []string{
						"subnet-03afefcc38a919083",
						"subnet-087ea2efd4f009fc5",
						"subnet-09e02553a59746ba1",
						"subnet-05a2d1112fbca071c",
						"subnet-0d5f650b4d0eebc56",
					},
				},
			},
		}

		for _, tt := range cases {
			tfcRespJsonTestdataPath = tt.tfcRespJsonTestdataPath

			result, err := cs.read(ctx, log, []*component.ConfigRequest{
				{
					Name: "HOSTNAME",
					Config: map[string]string{
						"organization": "foocorp",
						"workspace":    "databases",
						"all_outputs":  "true",
					},
				},
			})
			require.NoError(err)
			require.NotNil(result)
			require.Len(result, 1)

			v := result[0]

			ve, ok := v.Result.(*pb.ConfigSource_Value_Error)
			if ok {
				require.False(ok, "error reading value: %s", ve.Error.Message)
			}

			jsonResult, ok := v.Result.(*pb.ConfigSource_Value_Json)
			require.True(ok)

			require.NotEmpty(jsonResult.Json)

			var jsonResults map[string]interface{}
			require.NoError(json.Unmarshal(jsonResult.Json, &jsonResults))
			require.Len(jsonResults, tt.numOutputs)
			for k, v := range tt.spotChecks {
				actualV, ok := jsonResults[k]
				require.True(ok)

				switch actualV.(type) {
				case string:
					require.Equal(v, actualV)
				case []interface{}:
					require.ElementsMatch(v, actualV)
				}
			}
		}

	})
}
