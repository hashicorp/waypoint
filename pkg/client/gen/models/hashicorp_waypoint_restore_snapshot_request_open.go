// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpWaypointRestoreSnapshotRequestOpen hashicorp waypoint restore snapshot request open
//
// swagger:model hashicorp.waypoint.RestoreSnapshotRequest.Open
type HashicorpWaypointRestoreSnapshotRequestOpen struct {

	// If true, the server will exit after the restore is staged. This will
	// SHUT DOWN the server and some external process you created is expected
	// to bring it back. The Waypoint server on its own WILL NOT automatically
	// restart. You should only set this if you have some operation to
	// automate restart such as running in Nomad or Kubernetes.
	Exit bool `json:"exit,omitempty"`
}

// Validate validates this hashicorp waypoint restore snapshot request open
func (m *HashicorpWaypointRestoreSnapshotRequestOpen) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hashicorp waypoint restore snapshot request open based on context it is used
func (m *HashicorpWaypointRestoreSnapshotRequestOpen) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointRestoreSnapshotRequestOpen) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointRestoreSnapshotRequestOpen) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointRestoreSnapshotRequestOpen
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
