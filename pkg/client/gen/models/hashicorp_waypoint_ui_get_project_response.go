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

// HashicorpWaypointUIGetProjectResponse hashicorp waypoint UI get project response
//
// swagger:model hashicorp.waypoint.UI.GetProjectResponse
type HashicorpWaypointUIGetProjectResponse struct {

	// latest init job
	LatestInitJob *HashicorpWaypointJob `json:"latest_init_job,omitempty"`

	// project
	Project *HashicorpWaypointProject `json:"project,omitempty"`
}

// Validate validates this hashicorp waypoint UI get project response
func (m *HashicorpWaypointUIGetProjectResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLatestInitJob(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProject(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointUIGetProjectResponse) validateLatestInitJob(formats strfmt.Registry) error {
	if swag.IsZero(m.LatestInitJob) { // not required
		return nil
	}

	if m.LatestInitJob != nil {
		if err := m.LatestInitJob.SwaggerValidate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("latest_init_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("latest_init_job")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointUIGetProjectResponse) validateProject(formats strfmt.Registry) error {
	if swag.IsZero(m.Project) { // not required
		return nil
	}

	if m.Project != nil {
		if err := m.Project.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("project")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("project")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint UI get project response based on the context it is used
func (m *HashicorpWaypointUIGetProjectResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLatestInitJob(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateProject(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointUIGetProjectResponse) contextValidateLatestInitJob(ctx context.Context, formats strfmt.Registry) error {

	if m.LatestInitJob != nil {
		if err := m.LatestInitJob.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("latest_init_job")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("latest_init_job")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointUIGetProjectResponse) contextValidateProject(ctx context.Context, formats strfmt.Registry) error {

	if m.Project != nil {
		if err := m.Project.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("project")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("project")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointUIGetProjectResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointUIGetProjectResponse) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointUIGetProjectResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
