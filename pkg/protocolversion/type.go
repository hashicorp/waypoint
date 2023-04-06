// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protocolversion

//go:generate stringer -type=Type -linecomment

// Type is the enum of protocol version types.
type Type uint8

const (
	Invalid    Type = iota // invalid
	Api                    // api
	Entrypoint             // entrypoint
)
