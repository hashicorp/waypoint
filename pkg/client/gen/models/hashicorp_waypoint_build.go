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

// HashicorpWaypointBuild Build represents a process of creating an artifact that can be in any state,
// such as complete. A successful complete build produces an artifact.
//
// swagger:model hashicorp.waypoint.Build
type HashicorpWaypointBuild struct {

	// The application that this build is part of.
	Application *HashicorpWaypointRefApplication `json:"application,omitempty"`

	// artifact is the result of the build if the status.state == SUCCESS
	Artifact *HashicorpWaypointArtifact `json:"artifact,omitempty"`

	// component is the component that was used for this build
	Component *HashicorpWaypointComponent `json:"component,omitempty"`

	// id is the unique ID of the build
	ID string `json:"id,omitempty"`

	// ID of the job that created this build. This may be empty.
	JobID string `json:"job_id,omitempty"`

	// labels are the set of labels that are present on this build.
	Labels map[string]string `json:"labels,omitempty"`

	// Preload is information that is available via further queries but it
	// sometimes pre-populated in the initial load (see the field docs for more
	// info).
	Preload *HashicorpWaypointBuildPreload `json:"preload,omitempty"`

	// The sequence number for this build.
	Sequence string `json:"sequence,omitempty"`

	// status of the build
	Status *HashicorpWaypointStatus `json:"status,omitempty"`

	// template data for HCL variables and template functions, json-encoded
	// Format: byte
	TemplateData strfmt.Base64 `json:"template_data,omitempty"`

	// The workspace that this exists in
	Workspace *HashicorpWaypointRefWorkspace `json:"workspace,omitempty"`
}

// Validate validates this hashicorp waypoint build
func (m *HashicorpWaypointBuild) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateApplication(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateArtifact(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateComponent(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePreload(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
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

func (m *HashicorpWaypointBuild) validateApplication(formats strfmt.Registry) error {
	if swag.IsZero(m.Application) { // not required
		return nil
	}

	if m.Application != nil {
		if err := m.Application.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("application")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("application")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) validateArtifact(formats strfmt.Registry) error {
	if swag.IsZero(m.Artifact) { // not required
		return nil
	}

	if m.Artifact != nil {
		if err := m.Artifact.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("artifact")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("artifact")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) validateComponent(formats strfmt.Registry) error {
	if swag.IsZero(m.Component) { // not required
		return nil
	}

	if m.Component != nil {
		if err := m.Component.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("component")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("component")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) validatePreload(formats strfmt.Registry) error {
	if swag.IsZero(m.Preload) { // not required
		return nil
	}

	if m.Preload != nil {
		if err := m.Preload.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("preload")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("preload")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if m.Status != nil {
		if err := m.Status.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("status")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("status")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) validateWorkspace(formats strfmt.Registry) error {
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

// ContextValidate validate this hashicorp waypoint build based on the context it is used
func (m *HashicorpWaypointBuild) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateApplication(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateArtifact(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateComponent(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePreload(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStatus(ctx, formats); err != nil {
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

func (m *HashicorpWaypointBuild) contextValidateApplication(ctx context.Context, formats strfmt.Registry) error {

	if m.Application != nil {
		if err := m.Application.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("application")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("application")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) contextValidateArtifact(ctx context.Context, formats strfmt.Registry) error {

	if m.Artifact != nil {
		if err := m.Artifact.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("artifact")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("artifact")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) contextValidateComponent(ctx context.Context, formats strfmt.Registry) error {

	if m.Component != nil {
		if err := m.Component.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("component")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("component")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) contextValidatePreload(ctx context.Context, formats strfmt.Registry) error {

	if m.Preload != nil {
		if err := m.Preload.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("preload")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("preload")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) contextValidateStatus(ctx context.Context, formats strfmt.Registry) error {

	if m.Status != nil {
		if err := m.Status.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("status")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("status")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointBuild) contextValidateWorkspace(ctx context.Context, formats strfmt.Registry) error {

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
func (m *HashicorpWaypointBuild) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointBuild) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointBuild
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
