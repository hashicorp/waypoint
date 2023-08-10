// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !linux
// +build !linux

package freeport

func getEphemeralPortRange() (int, int, error) {
	return 0, 0, nil
}
