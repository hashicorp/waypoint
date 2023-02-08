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

// HashicorpWaypointGetTaskResponse hashicorp waypoint get task response
//
// swagger:model hashicorp.waypoint.GetTaskResponse
type HashicorpWaypointGetTaskResponse struct {

	// start job
	StartJob *HashicorpWaypointJob `json:"start_job,omitempty"`

	// stop job
	StopJob *HashicorpWaypointJob `json:"stop_job,omitempty"`

	// The requested Task
	Task *HashicorpWaypointTask `json:"task,omitempty"`

	// The Job triple that the task is associated with. These jobs are the full
	// message for each based on the Ref_Job id found inside Task
	TaskJob *HashicorpWaypointJob `json:"task_job,omitempty"`

	// watch job
	WatchJob *HashicorpWaypointJob `json:"watch_job,omitempty"`
}

// Validate validates this hashicorp waypoint get task response
func (m *HashicorpWaypointGetTaskResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStartJob(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStopJob(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTask(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTaskJob(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWatchJob(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointGetTaskResponse) validateStartJob(formats strfmt.Registry) error {
	if swag.IsZero(m.StartJob) { // not required
		return nil
	}

	if m.StartJob != nil {
		if err := m.StartJob.SwaggerValidate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("start_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("start_job")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetTaskResponse) validateStopJob(formats strfmt.Registry) error {
	if swag.IsZero(m.StopJob) { // not required
		return nil
	}

	if m.StopJob != nil {
		if err := m.StopJob.SwaggerValidate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stop_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("stop_job")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetTaskResponse) validateTask(formats strfmt.Registry) error {
	if swag.IsZero(m.Task) { // not required
		return nil
	}

	if m.Task != nil {
		if err := m.Task.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("task")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetTaskResponse) validateTaskJob(formats strfmt.Registry) error {
	if swag.IsZero(m.TaskJob) { // not required
		return nil
	}

	if m.TaskJob != nil {
		if err := m.TaskJob.SwaggerValidate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("task_job")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetTaskResponse) validateWatchJob(formats strfmt.Registry) error {
	if swag.IsZero(m.WatchJob) { // not required
		return nil
	}

	if m.WatchJob != nil {
		if err := m.WatchJob.SwaggerValidate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("watch_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("watch_job")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint get task response based on the context it is used
func (m *HashicorpWaypointGetTaskResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateStartJob(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStopJob(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTask(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTaskJob(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateWatchJob(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointGetTaskResponse) contextValidateStartJob(ctx context.Context, formats strfmt.Registry) error {

	if m.StartJob != nil {
		if err := m.StartJob.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("start_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("start_job")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetTaskResponse) contextValidateStopJob(ctx context.Context, formats strfmt.Registry) error {

	if m.StopJob != nil {
		if err := m.StopJob.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stop_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("stop_job")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetTaskResponse) contextValidateTask(ctx context.Context, formats strfmt.Registry) error {

	if m.Task != nil {
		if err := m.Task.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("task")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetTaskResponse) contextValidateTaskJob(ctx context.Context, formats strfmt.Registry) error {

	if m.TaskJob != nil {
		if err := m.TaskJob.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("task_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("task_job")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointGetTaskResponse) contextValidateWatchJob(ctx context.Context, formats strfmt.Registry) error {

	if m.WatchJob != nil {
		if err := m.WatchJob.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("watch_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("watch_job")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointGetTaskResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointGetTaskResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointGetTaskResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
