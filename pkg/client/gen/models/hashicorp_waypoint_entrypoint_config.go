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

// HashicorpWaypointEntrypointConfig hashicorp waypoint entrypoint config
//
// swagger:model hashicorp.waypoint.EntrypointConfig
type HashicorpWaypointEntrypointConfig struct {

	// The configuration for any config sources that may be used in the
	// config vars sent down. The server may send down extra configs that
	// aren't used so consumers should filter these based on what env vars
	// are actually in use.
	ConfigSources []*HashicorpWaypointConfigSource `json:"config_sources"`

	// Deployment is the deployment information for this instance. This may
	// be nil if the user is running an old enough server so always nil-check this.
	Deployment *HashicorpWaypointEntrypointConfigDeploymentInfo `json:"deployment,omitempty"`

	// The environment variables to set in the entrypoint.
	EnvVars []*HashicorpWaypointConfigVar `json:"env_vars"`

	// Exec are requested exec sessions for this instance.
	Exec []*HashicorpWaypointEntrypointConfigExec `json:"exec"`

	// The signal to send the application when config files change.
	FileChangeSignal string `json:"file_change_signal,omitempty"`

	// The URL service configuration. This might be nil. If this is nil,
	// then the URL service is disabled.
	URLService *HashicorpWaypointEntrypointConfigURLService `json:"url_service,omitempty"`
}

// Validate validates this hashicorp waypoint entrypoint config
func (m *HashicorpWaypointEntrypointConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateConfigSources(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDeployment(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEnvVars(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExec(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateURLService(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointEntrypointConfig) validateConfigSources(formats strfmt.Registry) error {
	if swag.IsZero(m.ConfigSources) { // not required
		return nil
	}

	for i := 0; i < len(m.ConfigSources); i++ {
		if swag.IsZero(m.ConfigSources[i]) { // not required
			continue
		}

		if m.ConfigSources[i] != nil {
			if err := m.ConfigSources[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("config_sources" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("config_sources" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointEntrypointConfig) validateDeployment(formats strfmt.Registry) error {
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

func (m *HashicorpWaypointEntrypointConfig) validateEnvVars(formats strfmt.Registry) error {
	if swag.IsZero(m.EnvVars) { // not required
		return nil
	}

	for i := 0; i < len(m.EnvVars); i++ {
		if swag.IsZero(m.EnvVars[i]) { // not required
			continue
		}

		if m.EnvVars[i] != nil {
			if err := m.EnvVars[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("env_vars" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("env_vars" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointEntrypointConfig) validateExec(formats strfmt.Registry) error {
	if swag.IsZero(m.Exec) { // not required
		return nil
	}

	for i := 0; i < len(m.Exec); i++ {
		if swag.IsZero(m.Exec[i]) { // not required
			continue
		}

		if m.Exec[i] != nil {
			if err := m.Exec[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("exec" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("exec" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointEntrypointConfig) validateURLService(formats strfmt.Registry) error {
	if swag.IsZero(m.URLService) { // not required
		return nil
	}

	if m.URLService != nil {
		if err := m.URLService.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("url_service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("url_service")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint entrypoint config based on the context it is used
func (m *HashicorpWaypointEntrypointConfig) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateConfigSources(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDeployment(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateEnvVars(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateExec(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateURLService(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointEntrypointConfig) contextValidateConfigSources(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.ConfigSources); i++ {

		if m.ConfigSources[i] != nil {
			if err := m.ConfigSources[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("config_sources" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("config_sources" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointEntrypointConfig) contextValidateDeployment(ctx context.Context, formats strfmt.Registry) error {

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

func (m *HashicorpWaypointEntrypointConfig) contextValidateEnvVars(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.EnvVars); i++ {

		if m.EnvVars[i] != nil {
			if err := m.EnvVars[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("env_vars" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("env_vars" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointEntrypointConfig) contextValidateExec(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Exec); i++ {

		if m.Exec[i] != nil {
			if err := m.Exec[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("exec" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("exec" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointEntrypointConfig) contextValidateURLService(ctx context.Context, formats strfmt.Registry) error {

	if m.URLService != nil {
		if err := m.URLService.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("url_service")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("url_service")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointEntrypointConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointEntrypointConfig) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointEntrypointConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}