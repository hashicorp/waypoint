// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpWaypointWaypointHclFmtResponse hashicorp waypoint waypoint hcl fmt response
//
// swagger:model hashicorp.waypoint.WaypointHclFmtResponse
type HashicorpWaypointWaypointHclFmtResponse struct {

	// waypoint hcl
	// Format: byte
	WaypointHcl strfmt.Base64 `json:"waypoint_hcl,omitempty"`
}

// Validate validates this hashicorp waypoint waypoint hcl fmt response
func (m *HashicorpWaypointWaypointHclFmtResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hashicorp waypoint waypoint hcl fmt response based on context it is used
func (m *HashicorpWaypointWaypointHclFmtResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointWaypointHclFmtResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointWaypointHclFmtResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointWaypointHclFmtResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
