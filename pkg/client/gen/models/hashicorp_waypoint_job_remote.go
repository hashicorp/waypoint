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

// HashicorpWaypointJobRemote hashicorp waypoint job remote
//
// swagger:model hashicorp.waypoint.Job.Remote
type HashicorpWaypointJobRemote struct {

	// This corresponds with the implicit behavior associated with data source
	// polling, whereby if the polling is successful, we perform an Up operation.
	DeployOnChange bool `json:"deploy_on_change,omitempty"`

	// Description is information about how the Waypoint server
	// acquires the data.
	Description string `json:"description,omitempty"`

	// If remote refers to a git repo, git_remote will be partially populate
	// with information about which information within the git repo to use.
	GitRemote *HashicorpWaypointJobGit `json:"git_remote,omitempty"`
}

// Validate validates this hashicorp waypoint job remote
func (m *HashicorpWaypointJobRemote) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateGitRemote(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointJobRemote) validateGitRemote(formats strfmt.Registry) error {
	if swag.IsZero(m.GitRemote) { // not required
		return nil
	}

	if m.GitRemote != nil {
		if err := m.GitRemote.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("git_remote")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("git_remote")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint job remote based on the context it is used
func (m *HashicorpWaypointJobRemote) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateGitRemote(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointJobRemote) contextValidateGitRemote(ctx context.Context, formats strfmt.Registry) error {

	if m.GitRemote != nil {
		if err := m.GitRemote.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("git_remote")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("git_remote")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointJobRemote) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointJobRemote) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointJobRemote
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
