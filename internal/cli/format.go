// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type idFormat struct {
	long bool
}

func initIdFormat(f *flag.Set, id *idFormat) {
	f.BoolVar(&flag.BoolVar{
		Name:   "long-ids",
		Target: &id.long,
		Usage:  "Show long identifiers rather than sequence numbers.",
	})
}

func (i *idFormat) FormatId(seq uint64, long string) string {
	if i.long {
		return fmt.Sprintf("%d (%s)", seq, long)
	}

	return strconv.FormatUint(seq, 10)
}
