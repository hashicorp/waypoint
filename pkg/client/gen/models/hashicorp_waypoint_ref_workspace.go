// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpWaypointRefWorkspace Workspace references a workspace.
//
// swagger:model hashicorp.waypoint.Ref.Workspace
type HashicorpWaypointRefWorkspace struct {

	// workspace
	Workspace string `json:"workspace,omitempty"`
}

// Validate validates this hashicorp waypoint ref workspace
func (m *HashicorpWaypointRefWorkspace) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hashicorp waypoint ref workspace based on context it is used
func (m *HashicorpWaypointRefWorkspace) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointRefWorkspace) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointRefWorkspace) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointRefWorkspace
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
