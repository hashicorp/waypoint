package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/internal"
	"github.com/mitchellh/devflow/internal/builtin"
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

	source := &internal.Source{App: "myapp", Path: "."}
	fn := builtin.BuilderM.Mapper("pack", source, log)
	if fn == nil {
		println(fmt.Sprintf("NO BUILDER"))
		return 1
	}

	builder, err := fn()
	if err != nil {
		panic(err)
	}

	_, err = builder.(internal.Builder).Build(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("DONE")
	return 0
}
