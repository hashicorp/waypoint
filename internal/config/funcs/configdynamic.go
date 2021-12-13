package funcs

import (
	"reflect"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

func ConfigDynamicFuncs() map[string]function.Function {
	return map[string]function.Function{
		"configdynamic": function.New(&function.Spec{
			Params: []function.Parameter{
				{
					Name: "from",
					Type: cty.String,
				},

				{
					Name: "config",
					Type: cty.Map(cty.String),
				},
			},
			Type: function.StaticReturnType(TypeDynamicConfig),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				var config map[string]string
				if err := gocty.FromCtyValue(args[1], &config); err != nil {
					return cty.NilVal, err
				}

				return cty.CapsuleVal(TypeDynamicConfig, &pb.ConfigVar_DynamicVal{
					From:   args[0].AsString(),
					Config: config,
				}), nil
			},
		}),
	}
}

// TODO(izaak) do these need to be global like this?
var (
	TypeDynamicConfig = cty.Capsule("configval",
		reflect.TypeOf((*pb.ConfigVar_DynamicVal)(nil)).Elem())
)
