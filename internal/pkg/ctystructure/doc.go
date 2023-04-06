// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package crystructure helps convert map[string]interface{} values to cty.Values.
//
// This is useful for dynamically creating variables that may be available
// to cty-powered environments such as HCL.
package ctystructure
