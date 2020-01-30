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

	fmt.Println("DONE")
	return 0
}
