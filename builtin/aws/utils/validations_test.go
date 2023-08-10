// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package utils

import (
	"testing"
)

func TestValidateEcsMemCPUPair(t *testing.T) {
	// test values based off of
	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-cpu-memory-error.html
	// circa July 15, 2021

	cases := map[string]struct {
		mem       int
		cpu       int
		shouldErr bool
	}{
		"zeros": {
			shouldErr: true,
		},
		"512/0": {
			mem: 512,
		},
		"512/256": {
			mem: 512,
			cpu: 256,
		},
		"4096": {
			mem: 4096,
		},
		"4096/512": {
			mem: 4096,
			cpu: 512,
		},
		"4096/256": {
			mem:       4096,
			cpu:       256,
			shouldErr: true,
		},
		"512/512": {
			mem:       512,
			cpu:       512,
			shouldErr: true,
		},
		"nonsense": {
			mem:       7,
			shouldErr: true,
		},
		"bad_pair": {
			mem:       512,
			cpu:       512,
			shouldErr: true,
		},
		"zero_mem": {
			cpu:       7,
			shouldErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := ValidateEcsMemCPUPair(c.mem, c.cpu); err != nil {
				if !c.shouldErr {
					t.Error(err)
				}
			}
		})
	}
}
