// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hcerr

import (
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserError is a custom error type designed to package
// details that will be displayed to a user, alongside
// some internal error message that may be unsafe to surface.
//
// UserError also implements grpcError, allowing it to optionally
// contain a custom grpc status code, which will be returned to the
// user if set.
type UserError struct {
	UserMessage string
	err         error
	statusCode  *codes.Code // Optional - nil represents unset
}

func (m *UserError) Error() string {
	if m.err == nil {
		return m.UserMessage
	}

	return fmt.Sprintf("%s: (user message: %q)", m.err.Error(), m.UserMessage)
}

func (m *UserError) Unwrap() error {
	return m.err
}

// GRPCStatus implements grpcError. If no code has been set,
// and no status error exists further up the error stack,
// will default to Internal with no message.
func (m *UserError) GRPCStatus() *status.Status {

	// status.Status does not support errors.As (https://github.com/grpc/grpc-go/issues/2934)
	var grpcstatus interface{ GRPCStatus() *status.Status }

	// If we don't have a status, return the status of errors further up the chain. Otherwise,
	// anyone calling GRPCStatus() on us will get code zero.
	if m.statusCode == nil && m.err != nil && errors.As(m.err, &grpcstatus) {
		// This error has no code, and there's another grpc status
		// further up the chain, so use that
		return grpcstatus.GRPCStatus()
	}

	code := m.statusCode
	// If no other code is set, default this to an internal type error.
	if code == nil {
		i := codes.Internal
		code = &i
	}

	return status.New(*code, "")
}

// NewUserError returns a new error with a message intended to be seen by server callers.
// The message should not contain anything sensitive or internal (i.e. database addresses,
// details of our internal processing, etc), and should be helpful to users in diagnosing
// why this error occurred and what they can do to avoid it in the future.
func NewUserError(message string) error {
	return &UserError{
		UserMessage: message,
	}
}

// NewUserErrorf is the same as NewUserError, with string formatting for convenience
func NewUserErrorf(format string, a ...interface{}) error {
	return &UserError{
		UserMessage: fmt.Sprintf(format, a...),
	}
}

// UserErrorf wraps an existing error with a UserError
func UserErrorf(err error, format string, a ...interface{}) error {
	return &UserError{
		UserMessage: fmt.Sprintf(format, a...),
		err:         err,
	}
}

// UserErrorWithCodef wraps an existing error with a UserError, and takes
// a grpc status code, which will be used by the final hcerr.Externalize
// call as the status code to present to the user.
func UserErrorWithCodef(c codes.Code, err error, format string, a ...interface{}) error {
	return &UserError{
		UserMessage: fmt.Sprintf(format, a...),
		err:         err,
		statusCode:  &c,
	}
}

// UserConditionWithCodef generates a new error that takes
// a grpc status code, which will be used by the final hcerr.Externalize
// call as the status code to present to the user.
func UserConditionWithCodef(c codes.Code, format string, a ...interface{}) error {
	return &UserError{
		UserMessage: fmt.Sprintf(format, a...),
		statusCode:  &c,
	}
}
