// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

func OptionalInt64(v int64) *int64 {
	if v == 0 {
		return nil
	}

	return &v
}
