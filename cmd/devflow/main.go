package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/internal/builtin"
	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/mapper"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log := hclog.New(&hclog.LoggerOptions{
		Name:  "devflow",
		Level: hclog.Trace,
		Color: hclog.AutoColor,
	})

	source := &component.Source{App: "myapp", Path: "."}

	// Builder
	fn := builtin.Builders.Func("pack")
	if fn == nil {
		println(fmt.Sprintf("NO BUILDER"))
		return 1
	}

	builder, err := fn.Call(source, log)
	if err != nil {
		panic(err)
	}

	buildFunc, err := mapper.NewFunc(builder.(component.Builder).BuildFunc())
	if err != nil {
		panic(err)
	}

	_, err = buildFunc.Call(context.Background(), source, log)
	if err != nil {
		panic(err)
	}

	// Registry
	fn = builtin.Registries.Func("docker")
	if fn == nil {
		println(fmt.Sprintf("NO REGISTRY"))
		return 1
	}

	registry, err := fn.Call(source, log)
	if err != nil {
		panic(err)
	}

	pushFunc, err := mapper.NewFunc(registry.(component.Registry).PushFunc())
	if err != nil {
		panic(err)
	}

	chain, err := pushFunc.Chain(builtin.Mappers, context.Background(), source, log)
	if err != nil {
		panic(err)
	}

	artifact, err := chain.Call()
	if err != nil {
		panic(err)
	}

	fmt.Printf("DONE: %#v\n", artifact)
	return 0
}
