// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpWaypointVersionInfoProtocolVersion hashicorp waypoint version info protocol version
//
// swagger:model hashicorp.waypoint.VersionInfo.ProtocolVersion
type HashicorpWaypointVersionInfoProtocolVersion struct {

	// current
	Current int64 `json:"current,omitempty"`

	// minimum
	Minimum int64 `json:"minimum,omitempty"`
}

// Validate validates this hashicorp waypoint version info protocol version
func (m *HashicorpWaypointVersionInfoProtocolVersion) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hashicorp waypoint version info protocol version based on context it is used
func (m *HashicorpWaypointVersionInfoProtocolVersion) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointVersionInfoProtocolVersion) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointVersionInfoProtocolVersion) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointVersionInfoProtocolVersion
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
