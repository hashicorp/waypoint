package plugin

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/go-argmapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/component/mocks"
	"github.com/hashicorp/waypoint/sdk/internal/testproto"
	pb "github.com/hashicorp/waypoint/sdk/proto"
)

func TestBuilderBuild(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	called := false
	buildFunc := func(ctx context.Context, args *component.Source) *testproto.Data {
		called = true
		assert.NotNil(ctx)
		assert.Equal("foo", args.App)
		return &testproto.Data{Value: "hello"}
	}

	mockB := &mocks.Builder{}
	mockB.On("BuildFunc").Return(buildFunc)

	plugins := Plugins(WithComponents(mockB), WithMappers(testDefaultMappers(t)...))
	client, server := plugin.TestPluginGRPCConn(t, plugins[1])
	defer client.Close()
	defer server.Stop()

	raw, err := client.Dispense("builder")
	require.NoError(err)
	builder := raw.(component.Builder)
	f := builder.BuildFunc().(*argmapper.Func)
	require.NotNil(f)

	result := f.Call(
		argmapper.Typed(context.Background()),
		argmapper.Typed(&pb.Args_Source{App: "foo"}),
		argmapper.ConverterGen(func(v argmapper.Value) (*argmapper.Func, error) {
			anyType := reflect.TypeOf((*any.Any)(nil))
			protoMessageType := reflect.TypeOf((*proto.Message)(nil)).Elem()
			if !v.Type.Implements(protoMessageType) {
				return nil, nil
			}

			// We take this value as our input.
			inputSet, err := argmapper.NewValueSet([]argmapper.Value{v})
			if err != nil {
				return nil, err
			}

			// Generate an int with the subtype of the string value
			outputSet, err := argmapper.NewValueSet([]argmapper.Value{argmapper.Value{
				Name:    v.Name,
				Type:    anyType,
				Subtype: proto.MessageName(reflect.Zero(v.Type).Interface().(proto.Message)),
			}})
			if err != nil {
				return nil, err
			}

			return argmapper.BuildFunc(inputSet, outputSet, func(in, out *argmapper.ValueSet) error {
				anyVal, err := ptypes.MarshalAny(inputSet.Typed(v.Type).Value.Interface().(proto.Message))
				if err != nil {
					return err
				}

				outputSet.Typed(anyType).Value = reflect.ValueOf(anyVal)
				return nil
			})

		}),
	)
	require.NoError(result.Err())

	raw = result.Out(0)
	require.NotNil(raw)
	require.Implements((*component.Artifact)(nil), raw)

	anyVal := raw.(component.ProtoMarshaler).Proto().(*any.Any)
	name, err := ptypes.AnyMessageName(anyVal)
	require.NoError(err)
	require.Equal("testproto.Data", name)

	require.True(called)
}

func TestBuilderConfig(t *testing.T) {
	mockV := &mockBuilderConfigurable{}
	testConfigurable(t, "builder", mockV, &mockV.Configurable)
}

type mockBuilderConfigurable struct {
	mocks.Builder
	mocks.Configurable
}
