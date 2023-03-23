// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package hcerr contains helpers to format and sanitize errors before returning
// them to clients as grpc status errors.
//
// At the service layer, it's not safe to assume that errors produced from lower
// layers can be returned directly to callers. Errors may contain sensitive
// and irrelevant information, like the internal address of a database that
// we've failed to connect to.
//
// Hcerr helps the service layer log context from these errors, and return
// a known safe error to the caller.
//
// If an error producer deeper in the call stack wants to produce an error
// that a caller will see, they can use hcerr.UserErrorf()
package hcerr
