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
	"github.com/go-openapi/validate"
)

// HashicorpWaypointWorkspaceProject hashicorp waypoint workspace project
//
// swagger:model hashicorp.waypoint.Workspace.Project
type HashicorpWaypointWorkspaceProject struct {

	// active_time is the last time that this project had activity in
	// this workspace.
	// Format: date-time
	ActiveTime strfmt.DateTime `json:"active_time,omitempty"`

	// The list of applications that have executed at least one operation
	// within the context of this workspace. To determine which operations
	// you must call the respect list API for that operation, such as
	// ListDeployments.
	Applications []*HashicorpWaypointWorkspaceApplication `json:"applications"`

	// The last non-local ref that was used for any operation.
	DataSourceRef *HashicorpWaypointJobDataSourceRef `json:"data_source_ref,omitempty"`

	// Project that this is referencing.
	Project *HashicorpWaypointRefProject `json:"project,omitempty"`

	// Workspace that this project is part of. This will only be set
	// when using the GetProject API. This will ALWAYS BE NIL on workspace
	// list and get APIs.
	Workspace *HashicorpWaypointRefWorkspace `json:"workspace,omitempty"`
}

// Validate validates this hashicorp waypoint workspace project
func (m *HashicorpWaypointWorkspaceProject) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActiveTime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateApplications(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDataSourceRef(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProject(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWorkspace(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointWorkspaceProject) validateActiveTime(formats strfmt.Registry) error {
	if swag.IsZero(m.ActiveTime) { // not required
		return nil
	}

	if err := validate.FormatOf("active_time", "body", "date-time", m.ActiveTime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *HashicorpWaypointWorkspaceProject) validateApplications(formats strfmt.Registry) error {
	if swag.IsZero(m.Applications) { // not required
		return nil
	}

	for i := 0; i < len(m.Applications); i++ {
		if swag.IsZero(m.Applications[i]) { // not required
			continue
		}

		if m.Applications[i] != nil {
			if err := m.Applications[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("applications" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("applications" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointWorkspaceProject) validateDataSourceRef(formats strfmt.Registry) error {
	if swag.IsZero(m.DataSourceRef) { // not required
		return nil
	}

	if m.DataSourceRef != nil {
		if err := m.DataSourceRef.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data_source_ref")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("data_source_ref")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointWorkspaceProject) validateProject(formats strfmt.Registry) error {
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

func (m *HashicorpWaypointWorkspaceProject) validateWorkspace(formats strfmt.Registry) error {
	if swag.IsZero(m.Workspace) { // not required
		return nil
	}

	if m.Workspace != nil {
		if err := m.Workspace.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("workspace")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("workspace")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint workspace project based on the context it is used
func (m *HashicorpWaypointWorkspaceProject) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateApplications(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDataSourceRef(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateProject(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateWorkspace(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointWorkspaceProject) contextValidateApplications(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Applications); i++ {

		if m.Applications[i] != nil {
			if err := m.Applications[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("applications" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("applications" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointWorkspaceProject) contextValidateDataSourceRef(ctx context.Context, formats strfmt.Registry) error {

	if m.DataSourceRef != nil {
		if err := m.DataSourceRef.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data_source_ref")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("data_source_ref")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointWorkspaceProject) contextValidateProject(ctx context.Context, formats strfmt.Registry) error {

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

func (m *HashicorpWaypointWorkspaceProject) contextValidateWorkspace(ctx context.Context, formats strfmt.Registry) error {

	if m.Workspace != nil {
		if err := m.Workspace.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("workspace")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("workspace")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointWorkspaceProject) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointWorkspaceProject) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointWorkspaceProject
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
