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

// HashicorpWaypointJobQueueProjectOp QueueProjectOp queues a job for all applications in a project. The
// applications queued may not directly align with what can be found in
// ListProjects because the application list will be based on the config
// and not the database.
//
// swagger:model hashicorp.waypoint.Job.QueueProjectOp
type HashicorpWaypointJobQueueProjectOp struct {

	// The template for the job to queue for each application. The "application"
	// field will be overwritten for each application. All other fields are
	// untouched.
	JobTemplate *HashicorpWaypointJob `json:"job_template,omitempty"`
}

// Validate validates this hashicorp waypoint job queue project op
func (m *HashicorpWaypointJobQueueProjectOp) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateJobTemplate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointJobQueueProjectOp) validateJobTemplate(formats strfmt.Registry) error {
	if swag.IsZero(m.JobTemplate) { // not required
		return nil
	}

	if m.JobTemplate != nil {
		if err := m.JobTemplate.SwaggerValidate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("job_template")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("job_template")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint job queue project op based on the context it is used
func (m *HashicorpWaypointJobQueueProjectOp) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateJobTemplate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointJobQueueProjectOp) contextValidateJobTemplate(ctx context.Context, formats strfmt.Registry) error {

	if m.JobTemplate != nil {
		if err := m.JobTemplate.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("job_template")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("job_template")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointJobQueueProjectOp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointJobQueueProjectOp) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointJobQueueProjectOp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
