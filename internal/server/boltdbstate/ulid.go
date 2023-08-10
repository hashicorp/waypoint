// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"crypto/rand"

	ulidpkg "github.com/oklog/ulid"
)

// ulid returns a unique ULID.
func ulid() (string, error) {
	id, err := ulidpkg.New(ulidpkg.Now(), rand.Reader)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
