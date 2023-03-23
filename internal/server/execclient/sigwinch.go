// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !windows
// +build !windows

package execclient

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func registerSigwinch(winchCh chan os.Signal) {
	signal.Notify(winchCh, unix.SIGWINCH)
}
