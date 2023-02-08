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

// HashicorpWaypointAuthMethod AuthMethod is a mechanism for authenticating to the Waypoint server.
// An AuthMethod deals with AuthN only: it provides identity and trades
// that for a Waypoint token.
//
// swagger:model hashicorp.waypoint.AuthMethod
type HashicorpWaypointAuthMethod struct {

	// A selector to determine whether a user who authenticates using this
	// is allowed to authenticate at all. This runs before authentication
	// completes. This can be used to check group membership and so on.
	// Available variables depend on the auth method used.
	//
	// The syntax of this is this:
	// https://github.com/hashicorp/go-bexpr
	// (better docs to follow)
	AccessSelector string `json:"access_selector,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// human friendly name for display and description. This has no impact
	// internally and is only helpful for the UI and API. This is optional.
	DisplayName string `json:"display_name,omitempty"`

	// unique name for this auth method
	Name string `json:"name,omitempty"`

	// OIDC uses OpenID Connect for auth. OIDC is supported by most
	// major identity providers.
	Oidc *HashicorpWaypointAuthMethodOIDC `json:"oidc,omitempty"`
}

// Validate validates this hashicorp waypoint auth method
func (m *HashicorpWaypointAuthMethod) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateOidc(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointAuthMethod) validateOidc(formats strfmt.Registry) error {
	if swag.IsZero(m.Oidc) { // not required
		return nil
	}

	if m.Oidc != nil {
		if err := m.Oidc.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("oidc")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("oidc")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint auth method based on the context it is used
func (m *HashicorpWaypointAuthMethod) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateOidc(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointAuthMethod) contextValidateOidc(ctx context.Context, formats strfmt.Registry) error {

	if m.Oidc != nil {
		if err := m.Oidc.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("oidc")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("oidc")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointAuthMethod) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointAuthMethod) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointAuthMethod
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
