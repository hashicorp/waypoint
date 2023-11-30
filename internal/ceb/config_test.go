// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// This test has a race condition:
// Running TestCEB starts a CEB in a goroutine, which initalizes the logger
// Reading ceb.logger without explicitly waiting for it to be set could panic
// This could be solved in the future by giving the ceb an explicit Init step
// where we initialize the logger.
//go:build !race

package ceb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessAppEnv_logs(t *testing.T) {
	data := TestCEB(t)
	ceb := data.CEB

	ceb.processAppEnv([]string{
		envLogLevel + "=TRACE",
	})
	require.True(t, ceb.logger.IsTrace())

	ceb.processAppEnv([]string{
		envLogLevel + "=DEBUG",
	})
	require.True(t, ceb.logger.IsDebug())

	// Unset everything should stay the same.
	ceb.processAppEnv([]string{})
	require.True(t, ceb.logger.IsDebug())

	// Send bogus stuff to test that we don't crash
	ceb.processAppEnv([]string{envLogLevel})
	ceb.processAppEnv([]string{envLogLevel + "="})
	ceb.processAppEnv([]string{envLogLevel + "=="})
	ceb.processAppEnv([]string{"=="})
	ceb.processAppEnv([]string{})
	require.True(t, ceb.logger.IsDebug())
}
