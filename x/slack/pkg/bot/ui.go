// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/olekukonko/tablewriter"
	"github.com/slack-go/slack"
)

type UI struct {
	ctx     context.Context
	L       hclog.Logger
	s       *slack.Client
	channel string
	ts      string

	primary string
}

// Input asks the user for input. This will immediately return an error
// if the UI doesn't support interaction. You can test for interaction
// ahead of time with Interactive().
func (u *UI) Input(_ *terminal.Input) (string, error) {
	return "", terminal.ErrNonInteractive
}

// Interactive returns true if this prompt supports user interaction.
// If this is false, Input will always error.
func (u *UI) Interactive() bool {
	return false
}

// Output outputs a message directly to the terminal. The remaining
// arguments should be interpolations for the format string. After the
// interpolations you may add Options.
func (u *UI) Output(str string, _ ...interface{}) {
	options := []slack.MsgOption{
		slack.MsgOptionText(str, true),
		slack.MsgOptionAsUser(true),
	}

	_, ts, err := u.s.PostMessage(u.channel, options...)
	if err != nil {
		u.L.Error("error posting output", "error", err)
		u.ts = ""
	} else {
		u.ts = ts
		u.primary = str
	}
}

// Output data as a table of data. Each entry is a row which will be output
// with the columns lined up nicely.
func (u *UI) NamedValues(nvs []terminal.NamedValue, _ ...terminal.Option) {
	var fields []slack.AttachmentField

	for _, nv := range nvs {
		fields = append(fields, slack.AttachmentField{
			Title: nv.Name,
			Value: fmt.Sprintf("%s", nv.Value),
		})
	}

	_, _, err := u.s.PostMessage(u.channel, slack.MsgOptionAttachments(
		slack.Attachment{
			Fields: fields,
		},
	))

	if err != nil {
		u.L.Error("error posting named values", "error", err)
	}
}

// OutputWriters returns stdout and stderr writers. These are usually
// but not always TTYs. This is useful for subprocesses, network requests,
// etc. Note that writing to these is not thread-safe by default so
// you must take care that there is only ever one writer.
func (u *UI) OutputWriters() (stdout io.Writer, stderr io.Writer, err error) {
	return nil, nil, terminal.ErrNonInteractive
}

// Status returns a live-updating status that can be used for single-line
// status updates that typically have a spinner or some similar style.
// While a Status is live (Close isn't called), other methods on UI should
// NOT be called.
func (u *UI) Status() terminal.Status {
	return &uiStatus{ui: u}
}

func (u *UI) Table(tbl *terminal.Table, opts ...terminal.Option) {
	var buf bytes.Buffer

	table := tablewriter.NewWriter(&buf)
	table.SetHeader(tbl.Headers)
	table.SetBorder(false)

	for _, row := range tbl.Rows {
		colors := make([]tablewriter.Colors, len(row))
		entries := make([]string, len(row))

		for i, ent := range row {
			entries[i] = ent.Value
		}

		table.Rich(entries, colors)
	}

	table.Render()

	u.s.PostMessage(
		u.channel,
		slack.MsgOptionText(fmt.Sprintf("```\n%s\n```", buf.String()), false),
		slack.MsgOptionAsUser(true),
	)
}

// Table outputs the information formatted into a Table structure.
func (u *UI) FieldTable(tbl *terminal.Table, _ ...terminal.Option) {
	var attachments []slack.Attachment

	for _, row := range tbl.Rows {
		var fields []slack.AttachmentField

		for i, col := range row {
			fields = append(fields, slack.AttachmentField{
				Title: tbl.Headers[i],
				Value: col.Value,
			})
		}

		attachments = append(attachments, slack.Attachment{
			Fields: fields,
		})
	}

	_, _, err := u.s.PostMessage(u.channel, slack.MsgOptionAttachments(attachments...))
	if err != nil {
		u.L.Error("error posting named values", "error", err)
	}
}

// StepGroup returns a value that can be used to output individual (possibly
// parallel) steps that have their own message, status indicator, spinner, and
// body. No other output mechanism (Output, Input, Status, etc.) may be
// called until the StepGroup is complete.
func (u *UI) StepGroup() terminal.StepGroup {
	return &uiStepGroup{
		u:    u,
		done: make(chan struct{}),
	}
}

type uiStatus struct {
	ui *UI
	ts string
}

// Update writes a new status. This should be a single line.
func (u *uiStatus) Update(msg string) {
	options := []slack.MsgOption{
		slack.MsgOptionText(msg, true),
		slack.MsgOptionAsUser(true),
	}

	if u.ts != "" {
		options = append(options, slack.MsgOptionUpdate(u.ts))
	}

	_, ts, err := u.ui.s.PostMessage(u.ui.channel, options...)
	if err != nil {
		u.ui.L.Error("error posting output", "error", err)
		u.ts = ""
	} else {
		u.ts = ts
	}
}

var statusEmoji = map[string]string{
	terminal.StatusOK:      "green_check_mark",
	terminal.StatusWarn:    "warning",
	terminal.StatusError:   "red",
	terminal.StatusTimeout: "hourglass",
}

// Indicate that a step has finished, confering an ok, error, or warn upon
// it's finishing state. If the status is not StatusOK, StatusError, or StatusWarn
// then the status text is written directly to the output, allowing for custom
// statuses.
func (u *uiStatus) Step(status string, msg string) {
	options := []slack.MsgOption{
		slack.MsgOptionText(msg, true),
		slack.MsgOptionAsUser(true),
	}

	if u.ts != "" {
		options = append(options, slack.MsgOptionUpdate(u.ts))
	}

	_, ts, err := u.ui.s.PostMessage(u.ui.channel, options...)
	if err != nil {
		u.ui.L.Error("error posting output", "error", err)
		u.ts = ""
	} else {
		u.ts = ts
	}

	if emoji, ok := statusEmoji[status]; ok {
		var item slack.ItemRef
		item.Channel = u.ui.channel
		item.Timestamp = ts

		u.ui.s.AddReaction(emoji, item)
	}
}

// Close should be called when the live updating is complete. The
// status will be cleared from the line.
func (u *uiStatus) Close() error {
	return nil
}

type uiStep struct {
	sg *uiStepGroup

	num int

	msg  string
	ts   string
	done bool
}

// The Writer has data written to it as though it was a terminal. This will appear
// as body text under the Step's message and status.
func (u *uiStep) TermOutput() io.Writer {
	return ioutil.Discard
}

// Change the Steps displayed message
func (u *uiStep) Update2(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)

	u.msg = msg

	options := []slack.MsgOption{
		slack.MsgOptionText("> "+msg, true),
		slack.MsgOptionAsUser(true),
		slack.MsgOptionUpdate(u.ts),
	}

	_, ts, err := u.sg.u.s.PostMessage(u.sg.u.channel, options...)
	if err != nil {
		u.sg.u.L.Error("error posting output", "error", err)
	}

	u.ts = ts
}

func (u *uiStep) Update(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)

	u.sg.body[u.num] = msg

	u.msg = msg

	u.sg.flush()
}

func (u *uiStep) Status(status string) {
	if emoji, ok := statusEmoji[status]; ok {
		msg := fmt.Sprintf(":%s: %s", emoji, u.msg)
		u.sg.body[u.num] = msg
		u.sg.flush()
	}
}

// Update the status of the message. Supported values are in status.go.
func (u *uiStep) Status2(status string) {
	if emoji, ok := statusEmoji[status]; ok {
		var item slack.ItemRef
		item.Channel = u.sg.u.channel
		item.Timestamp = u.ts

		msg := fmt.Sprintf("> :%s: %s", emoji, u.msg)

		options := []slack.MsgOption{
			slack.MsgOptionText(msg, true),
			slack.MsgOptionAsUser(true),
			slack.MsgOptionUpdate(u.ts),
		}

		_, ts, err := u.sg.u.s.PostMessage(u.sg.u.channel, options...)
		if err != nil {
			u.sg.u.L.Error("error posting output", "error", err)
		}

		u.ts = ts
	}
}

// Update the status of the message. Supported values are in status.go.
func (u *uiStep) ReactStatus(status string) {
	if emoji, ok := statusEmoji[status]; ok {
		var item slack.ItemRef
		item.Channel = u.sg.u.channel
		item.Timestamp = u.ts

		u.sg.u.s.AddReaction(emoji, item)
	}
}

// Called when the step has finished. This must be done otherwise the StepGroup
// will wait forever for its Steps to finish.
func (u *uiStep) Done() {
	if u.done {
		return
	}

	u.Status(terminal.StatusOK)
	u.signalDone()
}

func (u *uiStep) signalDone() {
	go func() {
		select {
		case u.sg.done <- struct{}{}:
		case <-u.sg.u.ctx.Done():
			return
		}
	}()
}

// Sets the status to Error and finishes the Step if it's not already done.
// This is usually done in a defer so that any return before the Done() shows
// the Step didn't completely properly.
func (u *uiStep) Abort() {
	if u.done {
		return
	}

	u.Status(terminal.StatusError)
	u.signalDone()
}

type uiStepGroup struct {
	u     *UI
	steps int
	done  chan struct{}
	body  []string
	ts    string
}

// Start a step in the output with the arguments making up the initial message
func (u *uiStepGroup) Add2(str string, args ...interface{}) terminal.Step {
	u.steps++
	msg := fmt.Sprintf(str, args...)

	options := []slack.MsgOption{
		slack.MsgOptionText("> "+msg, true),
		slack.MsgOptionAsUser(true),
	}

	u.body = append(u.body, msg)

	_, ts, err := u.u.s.PostMessage(u.u.channel, options...)
	if err != nil {
		u.u.L.Error("error posting output", "error", err)
	}

	return &uiStep{sg: u, ts: ts, msg: msg}
}

func (u *uiStepGroup) Add(str string, args ...interface{}) terminal.Step {
	i := u.steps
	u.steps++
	msg := fmt.Sprintf(str, args...)

	u.body = append(u.body, ":hourglass: "+msg)

	u.flush()

	return &uiStep{sg: u, num: i, msg: msg}
}

func (u *uiStepGroup) flush() {
	msg := strings.Join(u.body, "\n> ")
	options := []slack.MsgOption{
		slack.MsgOptionText("> "+msg, true),
		slack.MsgOptionAsUser(true),
	}

	if u.ts != "" {
		options = append(options, slack.MsgOptionUpdate(u.ts))
	}

	_, ts, err := u.u.s.PostMessage(u.u.channel, options...)
	if err != nil {
		u.u.L.Error("error posting output", "error", err)
	}

	u.ts = ts
}

// Wait for all steps to finish. This allows a StepGroup to be used like
// a sync.WaitGroup with each step being run in a separate goroutine.
// This must be called to properly clean up the step group.
func (u *uiStepGroup) Wait() {
	for u.steps > 0 {
		select {
		case <-u.u.ctx.Done():
			return
		case <-u.done:
			u.steps--
		}
	}
}
