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

// HashicorpWaypointListPipelinesResponse hashicorp waypoint list pipelines response
//
// swagger:model hashicorp.waypoint.ListPipelinesResponse
type HashicorpWaypointListPipelinesResponse struct {

	// All of the listed pipelines
	Pipelines []*HashicorpWaypointPipeline `json:"pipelines"`
}

// Validate validates this hashicorp waypoint list pipelines response
func (m *HashicorpWaypointListPipelinesResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePipelines(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointListPipelinesResponse) validatePipelines(formats strfmt.Registry) error {
	if swag.IsZero(m.Pipelines) { // not required
		return nil
	}

	for i := 0; i < len(m.Pipelines); i++ {
		if swag.IsZero(m.Pipelines[i]) { // not required
			continue
		}

		if m.Pipelines[i] != nil {
			if err := m.Pipelines[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("pipelines" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("pipelines" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this hashicorp waypoint list pipelines response based on the context it is used
func (m *HashicorpWaypointListPipelinesResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePipelines(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointListPipelinesResponse) contextValidatePipelines(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Pipelines); i++ {

		if m.Pipelines[i] != nil {
			if err := m.Pipelines[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("pipelines" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("pipelines" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointListPipelinesResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointListPipelinesResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointListPipelinesResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
