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

// HashicorpWaypointPipelineRun hashicorp waypoint pipeline run
//
// swagger:model hashicorp.waypoint.PipelineRun
type HashicorpWaypointPipelineRun struct {

	// id
	ID string `json:"id,omitempty"`

	// The full list of jobs that are associated with this run.
	Jobs []*HashicorpWaypointRefJob `json:"jobs"`

	// The pipeline associated with this run.
	Pipeline *HashicorpWaypointRefPipeline `json:"pipeline,omitempty"`

	// The sequence number for this pipeline run.
	Sequence string `json:"sequence,omitempty"`

	// The current state of this pipeline run,
	State *HashicorpWaypointPipelineRunState `json:"state,omitempty"`
}

// Validate validates this hashicorp waypoint pipeline run
func (m *HashicorpWaypointPipelineRun) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateJobs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePipeline(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateState(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointPipelineRun) validateJobs(formats strfmt.Registry) error {
	if swag.IsZero(m.Jobs) { // not required
		return nil
	}

	for i := 0; i < len(m.Jobs); i++ {
		if swag.IsZero(m.Jobs[i]) { // not required
			continue
		}

		if m.Jobs[i] != nil {
			if err := m.Jobs[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("jobs" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("jobs" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointPipelineRun) validatePipeline(formats strfmt.Registry) error {
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

func (m *HashicorpWaypointPipelineRun) validateState(formats strfmt.Registry) error {
	if swag.IsZero(m.State) { // not required
		return nil
	}

	if m.State != nil {
		if err := m.State.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("state")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("state")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint pipeline run based on the context it is used
func (m *HashicorpWaypointPipelineRun) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateJobs(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePipeline(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateState(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointPipelineRun) contextValidateJobs(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Jobs); i++ {

		if m.Jobs[i] != nil {
			if err := m.Jobs[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("jobs" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("jobs" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointPipelineRun) contextValidatePipeline(ctx context.Context, formats strfmt.Registry) error {

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

func (m *HashicorpWaypointPipelineRun) contextValidateState(ctx context.Context, formats strfmt.Registry) error {

	if m.State != nil {
		if err := m.State.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("state")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("state")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointPipelineRun) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointPipelineRun) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointPipelineRun
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
