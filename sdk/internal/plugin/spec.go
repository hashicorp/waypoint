package plugin

import (
	"github.com/mitchellh/devflow/sdk/pkg/mapper"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

// funcToSpec takes a function pointer and generates a FuncSpec from it.
// The function is expected to only take and return proto.Message values.
func funcToSpec(f interface{}) (*pb.FuncSpec, error) {
	return nil, nil
}

// specToFunc takes a FuncSpec and returns a mapper.Func that can be called
// to invoke this function.
func specToFunc(s *pb.FuncSpec, cb interface{}) *mapper.Func {
	// Build the function
	f, err := mapper.NewFunc(cb, mapper.WithType(dynamicArgsType, makeDynamicArgsMapperType(s)))
	if err != nil {
		panic(err)
	}

	return f
}
