// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/internal/cli"
	"github.com/hashicorp/waypoint/x/slack/pkg/bot"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "plugin" {
		os.Exit(cli.Main(os.Args))
	}

	// TODO if this becomes a real program, use a real way to pass the token in.
	token := os.Getenv("SLACK_TOKEN")

	if token == "" {
		fmt.Fprintf(os.Stderr, "Missing slack token, set SLACK_TOKEN env var\n")
		os.Exit(1)
	}

	cfg := bot.Config{
		Token: token,
	}

	L := hclog.New(&hclog.LoggerOptions{
		Name:  "slack",
		Level: hclog.Trace,
		Color: hclog.AutoColor,
	})

	bot, err := bot.NewBot(cfg)
	if err != nil {
		L.Error("error creating bot", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)

	go func() {
		<-ch
		L.Info("Shutting down bot...")
		cancel()
	}()

	signal.Notify(ch, os.Interrupt)

	L.Info("connecting bot to slack...")

	err = bot.Run(ctx, L)
	if err != nil {
		if err == context.Canceled {
			return
		}

		L.Error("error creating bot", "error", err)
		os.Exit(1)
	}
}
