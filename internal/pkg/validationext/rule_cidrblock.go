// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validationext

import (
	"fmt"
	"net"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// IsPrivateCIDRBlock implements validation.Rule to check if a value is a valid
// IPv4 CIDR block within the ranges that are considered valid.
var IsPrivateCIDRBlock validation.Rule = &isPrivateCIDRBlockRule{}

// isPrivateCIDRBlockRule implements validation.Rule for IsPrivateCIDRBlock.
type isPrivateCIDRBlockRule struct{}

// validRanges contains the set of IP ranges considered valid.
var validRanges = []net.IPNet{
	{
		// 10.*.*.*
		IP:   net.IPv4(10, 0, 0, 0),
		Mask: net.IPv4Mask(255, 0, 0, 0),
	},
	{
		// 192.168.*.*
		IP:   net.IPv4(192, 168, 0, 0),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	},
	{
		// 172.[16-31].*.*
		IP:   net.IPv4(172, 16, 0, 0),
		Mask: net.IPv4Mask(255, 240, 0, 0),
	},
}

// Validate validates if the provided value is a valid IPv4 CIDR
// block contained within the valid ranges.
func (r *isPrivateCIDRBlockRule) Validate(value interface{}) error {
	// assert that value is of type string.
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a valid string")
	}

	// parse the string as CIDR notation IP address and prefix length.
	ip, net, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}

	// validate if the IP address is contained in one of the expected ranges.
	for _, validRange := range validRanges {
		valueSize, _ := net.Mask.Size()
		validRangeSize, _ := validRange.Mask.Size()
		if validRange.Contains(ip) && valueSize >= validRangeSize {
			return nil
		}
	}

	// return an error if the IP address is not contained within expected ranges.
	return fmt.Errorf("must match pattern of 10.*.*.* or 172.[16-31].*.* or " +
		"192.168.*.*; where * is any number from [0-255]")
}
