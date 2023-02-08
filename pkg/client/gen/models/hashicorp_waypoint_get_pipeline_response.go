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

// HashicorpWaypointGetPipelineResponse hashicorp waypoint get pipeline response
//
// swagger:model hashicorp.waypoint.GetPipelineResponse
type HashicorpWaypointGetPipelineResponse struct {

	// Graph is the execution graph for the pipeline steps. This can be
	// used to better visualize pipeline execution since the internal data
	// format of pipeline.steps is optimized more for storage rather than usage.
	Graph *HashicorpWaypointGetPipelineResponseGraph `json:"graph,omitempty"`

	// Pipeline is the pipeline that was requested.
	Pipeline *HashicorpWaypointPipeline `json:"pipeline,omitempty"`

	// Root step is the name of the step in pipeline that is the first
	// step executed.
	RootStep string `json:"root_step,omitempty"`
}

// Validate validates this hashicorp waypoint get pipeline response
func (m *HashicorpWaypointGetPipelineResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateGraph(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePipeline(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointGetPipelineResponse) validateGraph(formats strfmt.Registry) error {
	if swag.IsZero(m.Graph) { // not required
		return nil
	}

	if m.Graph != nil {
		if err := m.Graph.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("graph")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("graph")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetPipelineResponse) validatePipeline(formats strfmt.Registry) error {
	if swag.IsZero(m.Pipeline) { // not required
		return nil
	}

	if m.Pipeline != nil {
		if err := m.Pipeline.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pipeline")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("pipeline")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint get pipeline response based on the context it is used
func (m *HashicorpWaypointGetPipelineResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateGraph(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePipeline(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointGetPipelineResponse) contextValidateGraph(ctx context.Context, formats strfmt.Registry) error {

	if m.Graph != nil {
		if err := m.Graph.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("graph")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("graph")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetPipelineResponse) contextValidatePipeline(ctx context.Context, formats strfmt.Registry) error {

	if m.Pipeline != nil {
		if err := m.Pipeline.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pipeline")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("pipeline")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointGetPipelineResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointGetPipelineResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointGetPipelineResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
