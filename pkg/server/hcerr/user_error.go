package hcerr

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserError is a custom error type designed to package
// detail that will be displayed to a user, alongside
// some internal error message that may be unsafe to surface.
//
// UserError also implements grpcError, allowing it to optionally
// contain a custom grpc status code, which will be returned to the
// user if set.
type UserError struct {
	UserMessage string
	err         error
	statusCode  codes.Code
}

func (m UserError) Error() string {
	return fmt.Sprintf("%s: (user message: %q)", m.err.Error(), m.UserMessage)
}

func (m UserError) Unwrap() error {
	return m.err
}

// GRPCStatus implements grpcError
func (m UserError) GRPCStatus() *status.Status {
	return status.New(m.statusCode, "")
}

// NewUserError returns a new error with a message intended to be seen by server callers.
// The message should not contain anything sensitive or internal (i.e. database addresses,
// details of our internal processing, etc), and should be helpful to users in diagnosing
// why this error occurred and what they can do to avoid it in the future.
func NewUserError(message string) error {
	return UserError{
		UserMessage: message,
	}
}

// NewUserErrorf is the same as NewUserError, with string formatting for convenience
func NewUserErrorf(format string, a ...interface{}) error {
	return UserError{
		UserMessage: fmt.Sprintf(format, a...),
	}
}

// UserErrorf wraps an existing error with a UserError
func UserErrorf(err error, format string, a ...interface{}) error {
	return UserError{
		UserMessage: fmt.Sprintf(format, a...),
		err:         err,
	}
}

// UserErrorWithCodef wraps an existing error with a UserError, and takes
// a grpc status code, which will be used by the final hcerr.Externalize
// call as the status code to present to the user.
func UserErrorWithCodef(c codes.Code, err error, format string, a ...interface{}) error {
	return UserError{
		UserMessage: fmt.Sprintf(format, a...),
		err:         err,
		statusCode:  c,
	}
}
