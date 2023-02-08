// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GoogleRPCStatus The `Status` type defines a logical error model that is suitable for
// different programming environments, including REST APIs and RPC APIs. It is
// used by [gRPC](https://github.com/grpc). Each `Status` message contains
// three pieces of data: error code, error message, and error details.
//
// You can find out more about this error model and how to work with it in the
// [API Design Guide](https://cloud.google.com/apis/design/errors).
//
// swagger:model google.rpc.Status
type GoogleRPCStatus struct {

	// The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
	Code int32 `json:"code,omitempty"`

	// A list of messages that carry the error details.  There is a common set of
	// message types for APIs to use.
	Details []*GoogleProtobufAny `json:"details"`

	// A developer-facing error message, which should be in English. Any
	// user-facing error message should be localized and sent in the
	// [google.rpc.Status.details][google.rpc.Status.details] field, or localized by the client.
	Message string `json:"message,omitempty"`
}

// Validate validates this google rpc status
func (m *GoogleRPCStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDetails(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GoogleRPCStatus) validateDetails(formats strfmt.Registry) error {
	if swag.IsZero(m.Details) { // not required
		return nil
	}

	for i := 0; i < len(m.Details); i++ {
		if swag.IsZero(m.Details[i]) { // not required
			continue
		}

		if m.Details[i] != nil {
			if err := m.Details[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("details" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("details" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this google rpc status based on the context it is used
func (m *GoogleRPCStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDetails(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GoogleRPCStatus) contextValidateDetails(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Details); i++ {

		if m.Details[i] != nil {
			if err := m.Details[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("details" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("details" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *GoogleRPCStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GoogleRPCStatus) UnmarshalBinary(b []byte) error {
	var res GoogleRPCStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
