package funcspec

import (
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
	pb "github.com/mitchellh/devflow/sdk/proto"
)

// Func takes a FuncSpec and returns a *mapper.Func that can be called
// to invoke this function. The callback can have an argument type of Args
// in order to get access to the required dynamic proto.Any types of the
// FuncSpec.
func Func(s *pb.FuncSpec, cb interface{}, opts ...Option) *mapper.Func {
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	// Build the function
	f, err := mapper.NewFunc(cb,
		mapper.WithName(s.Name),
		mapper.WithType(ArgsType, makeArgsMapperType(s)),
		mapper.WithLogger(cfg.Logger),
		mapper.WithValues(cfg.Values...),
	)
	if err != nil {
		panic(err)
	}

	return f
}
