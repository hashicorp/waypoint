// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// HashicorpWaypointJobLogsOp Used to start a platform's log function within a runner. API users
// interested in viewing logs should use the GetLogStream API. This
// is only meant for implementing custom log handling by plugins.
//
// swagger:model hashicorp.waypoint.Job.LogsOp
type HashicorpWaypointJobLogsOp struct {

	// The deployment to create the exec session context. Ie, what
	// application code will be available within the exec session.
	Deployment *HashicorpWaypointDeployment `json:"deployment,omitempty"`

	// Id to assign the virtual instance created
	InstanceID string `json:"instance_id,omitempty"`

	// The maximum of log entries to be output.
	Limit int32 `json:"limit,omitempty"`

	// Indicates the time horizon that log entries must be beyond for them
	// to be emitted.
	// Format: date-time
	StartTime strfmt.DateTime `json:"start_time,omitempty"`
}

// Validate validates this hashicorp waypoint job logs op
func (m *HashicorpWaypointJobLogsOp) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDeployment(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartTime(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointJobLogsOp) validateDeployment(formats strfmt.Registry) error {
	if swag.IsZero(m.Deployment) { // not required
		return nil
	}

	if m.Deployment != nil {
		if err := m.Deployment.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("deployment")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("deployment")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointJobLogsOp) validateStartTime(formats strfmt.Registry) error {
	if swag.IsZero(m.StartTime) { // not required
		return nil
	}

	if err := validate.FormatOf("start_time", "body", "date-time", m.StartTime.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint job logs op based on the context it is used
func (m *HashicorpWaypointJobLogsOp) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDeployment(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointJobLogsOp) contextValidateDeployment(ctx context.Context, formats strfmt.Registry) error {

	if m.Deployment != nil {
		if err := m.Deployment.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("deployment")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("deployment")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointJobLogsOp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointJobLogsOp) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointJobLogsOp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
