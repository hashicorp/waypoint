// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package bot

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/google/shlex"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/internal/cli"
	mcli "github.com/mitchellh/cli"
	"github.com/slack-go/slack"
)

type Config struct {
	Token string
}

type Bot struct {
	api *slack.Client

	addressable string
}

func NewBot(cfg Config) (*Bot, error) {
	api := slack.New(cfg.Token)

	bot := &Bot{
		api: api,
	}

	return bot, nil
}

var unsupportedCommands = []string{"exec"}

func (b *Bot) Run(ctx context.Context, L hclog.Logger) error {
	rtm := b.api.NewRTM()

	go rtm.ManageConnection()

	defer rtm.Disconnect()

	var connected bool

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m := <-rtm.IncomingEvents:
			switch v := m.Data.(type) {
			case *slack.ConnectedEvent:
				L.Info("Connected to slack", "username", v.Info.User.Name, "id", v.Info.User.ID)
				b.addressable = "<@" + v.Info.User.ID + ">"
				connected = true
			case *slack.MessageEvent:
				L.Trace("incoming msg", "text", v.Text, "channel", v.Channel, "user", v.Username, "connected", connected)
				if connected {
					go b.handleMessage(ctx, L, v, rtm)
				}
			}
		}
	}
}

func (b *Bot) handleMessage(ctx context.Context, L hclog.Logger, m *slack.MessageEvent, rtm *slack.RTM) {
	parts, err := shlex.Split(m.Text)
	if err != nil {
		L.Trace("unable to parse slack message", "error", err)
		return
	}

	if len(parts) == 0 {
		L.Trace("empty message seen", "message", m.Text)
		return
	}

	if parts[0] != b.addressable {
		L.Trace("not being addressed, ignoring", "0th", parts[0])
		return
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg.Add(1)

	go func() {
		defer wg.Done()

		t := time.NewTicker(10 * time.Second)
		defer t.Stop()

	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case <-t.C:
				rtm.SendMessage(rtm.NewTypingMessage(m.Channel))
			}
		}

		rtm.PostMessage(m.Channel, slack.MsgOptionText("_Operation has completed_", false))
	}()

	var buf bytes.Buffer

	var ui UI
	ui.s = b.api
	ui.ctx = ctx
	ui.channel = m.Channel

	base, commands := cli.Commands(ctx, L, ioutil.Discard, cli.WithUI(&ui))
	defer base.Close()
	for _, cmd := range unsupportedCommands {
		delete(commands, cmd)
	}

	cli := &mcli.CLI{
		Name:       "waypoint-slack",
		Args:       parts[1:],
		Commands:   commands,
		HelpFunc:   cli.GroupedHelpFunc(mcli.BasicHelpFunc("waypoint-slack")),
		HelpWriter: &buf,
	}

	cli.Run()

	if buf.Len() != 0 {
		b.api.PostMessage(
			m.Channel,
			slack.MsgOptionText(fmt.Sprintf("<@%s>: ```\n%s\n```", m.User, buf.String()), false),
			slack.MsgOptionAsUser(true),
		)
	}

	cancel()
	wg.Wait()
}
