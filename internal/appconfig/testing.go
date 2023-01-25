package appconfig

import (
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	pb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

// StaticConfigSourcer can be used in tests to provide a config sourcer
// plugin type that simply returns the value of the `value` attribute
// in the `configdynamic` func call.
type StaticConfigSourcer struct{}

func (cs *StaticConfigSourcer) ReadFunc() interface{} {
	return cs.readFunc
}

func (cs *StaticConfigSourcer) StopFunc() interface{} {
	return func() {
		// Nothing.
	}
}

func (cs *StaticConfigSourcer) readFunc(
	reqs []*component.ConfigRequest,
) ([]*pb.ConfigSource_Value, error) {
	var results []*pb.ConfigSource_Value

	for _, req := range reqs {
		if v, ok := req.Config["value"]; ok && v != "" {
			result := &pb.ConfigSource_Value{Name: req.Name}
			result.Result = &pb.ConfigSource_Value_Value{
				Value: req.Config["value"],
			}
			results = append(results, result)
		}

		if v, ok := req.Config["json"]; ok && v != "" {
			result := &pb.ConfigSource_Value{Name: req.Name}
			result.Result = &pb.ConfigSource_Value_Json{
				Json: []byte(req.Config["json"]),
			}
			results = append(results, result)
		}
	}

	return results, nil
}
