// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpWaypointRunnerJobStreamRequestComplete hashicorp waypoint runner job stream request complete
//
// swagger:model hashicorp.waypoint.RunnerJobStreamRequest.Complete
type HashicorpWaypointRunnerJobStreamRequestComplete struct {

	// result
	Result *HashicorpWaypointJobResult `json:"result,omitempty"`
}

// Validate validates this hashicorp waypoint runner job stream request complete
func (m *HashicorpWaypointRunnerJobStreamRequestComplete) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateResult(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointRunnerJobStreamRequestComplete) validateResult(formats strfmt.Registry) error {
	if swag.IsZero(m.Result) { // not required
		return nil
	}

	if m.Result != nil {
		if err := m.Result.SwaggerValidate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("result")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("result")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint runner job stream request complete based on the context it is used
func (m *HashicorpWaypointRunnerJobStreamRequestComplete) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateResult(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointRunnerJobStreamRequestComplete) contextValidateResult(ctx context.Context, formats strfmt.Registry) error {

	if m.Result != nil {
		if err := m.Result.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("result")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("result")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointRunnerJobStreamRequestComplete) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointRunnerJobStreamRequestComplete) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointRunnerJobStreamRequestComplete
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
