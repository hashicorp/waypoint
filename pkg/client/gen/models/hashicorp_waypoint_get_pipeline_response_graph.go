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

// HashicorpWaypointGetPipelineResponseGraph Graph represents the execution graph for the pipeline steps. This
// may support multiple formats.
//
// swagger:model hashicorp.waypoint.GetPipelineResponse.Graph
type HashicorpWaypointGetPipelineResponseGraph struct {

	// content
	// Format: byte
	Content strfmt.Base64 `json:"content,omitempty"`

	// format
	Format *HashicorpWaypointGetPipelineResponseGraphFormat `json:"format,omitempty"`
}

// Validate validates this hashicorp waypoint get pipeline response graph
func (m *HashicorpWaypointGetPipelineResponseGraph) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFormat(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointGetPipelineResponseGraph) validateFormat(formats strfmt.Registry) error {
	if swag.IsZero(m.Format) { // not required
		return nil
	}

	if m.Format != nil {
		if err := m.Format.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("format")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("format")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint get pipeline response graph based on the context it is used
func (m *HashicorpWaypointGetPipelineResponseGraph) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateFormat(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointGetPipelineResponseGraph) contextValidateFormat(ctx context.Context, formats strfmt.Registry) error {

	if m.Format != nil {
		if err := m.Format.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("format")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("format")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointGetPipelineResponseGraph) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointGetPipelineResponseGraph) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointGetPipelineResponseGraph
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
