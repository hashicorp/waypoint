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

// HashicorpWaypointUpsertAuthMethodRequest hashicorp waypoint upsert auth method request
//
// swagger:model hashicorp.waypoint.UpsertAuthMethodRequest
type HashicorpWaypointUpsertAuthMethodRequest struct {

	// AuthMethod to upsert. See the message for what fields to set.
	AuthMethod *HashicorpWaypointAuthMethod `json:"auth_method,omitempty"`
}

// Validate validates this hashicorp waypoint upsert auth method request
func (m *HashicorpWaypointUpsertAuthMethodRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAuthMethod(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointUpsertAuthMethodRequest) validateAuthMethod(formats strfmt.Registry) error {
	if swag.IsZero(m.AuthMethod) { // not required
		return nil
	}

	if m.AuthMethod != nil {
		if err := m.AuthMethod.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("auth_method")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("auth_method")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint upsert auth method request based on the context it is used
func (m *HashicorpWaypointUpsertAuthMethodRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAuthMethod(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointUpsertAuthMethodRequest) contextValidateAuthMethod(ctx context.Context, formats strfmt.Registry) error {

	if m.AuthMethod != nil {
		if err := m.AuthMethod.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("auth_method")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("auth_method")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointUpsertAuthMethodRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointUpsertAuthMethodRequest) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointUpsertAuthMethodRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
