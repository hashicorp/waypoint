// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/hashicorp/go-hclog"
	"golang.org/x/sys/unix"
)

// debugSignalHandler writes a heap profile to the system temporary
// directory when SIGUSR1 is received. The output path for each dump
// is logged.
//
// NOTE(mitchellh): In the future, I expect that we'll dump a lot more
// data and perhaps store a tar file or some other format. Given the
// debug nature of this file, I'm not going to worry too much about tagging
// the file type and so on.
func debugSignalHandler(ctx context.Context, log hclog.Logger) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, unix.SIGUSR1)
	defer signal.Stop(ch)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ch:
		}

		// Dump the profile to a timestamped file
		path := filepath.Join(os.TempDir(), fmt.Sprintf(
			"waypoint_heap_%d", time.Now().Unix()))
		log.Warn("SIGUSR1 received, dumping heap profile", "path", path)

		// Get the profile, this is predefined by Go and should always exist.
		profile := pprof.Lookup("heap")
		if profile == nil {
			log.Error("heap profile not found, this should not be possible")
			continue
		}

		// Write it
		f, err := os.Create(path)
		if err != nil {
			log.Error("error opening file to dump profile", "err", err)
			continue
		}
		err = profile.WriteTo(f, 0)
		f.Close()
		if err != nil {
			log.Error("error writing heap profile", "err", err)
		}
	}
}
