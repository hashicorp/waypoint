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

// HashicorpWaypointRunner hashicorp waypoint runner
//
// swagger:model hashicorp.waypoint.Runner
type HashicorpWaypointRunner struct {

	// The state of whether this runner is adopted or not.
	AdoptionState *HashicorpWaypointRunnerAdoptionState `json:"adoption_state,omitempty"`

	// The runner will only be assigned jobs that directly target this
	// runner by ID. This is used by local runners to prevent external
	// jobs from being assigned to them.
	ByIDOnly bool `json:"by_id_only,omitempty"`

	// Components are the list of components that the runner supports. This
	// is used to match jobs to this runner.
	Components []*HashicorpWaypointComponent `json:"components"`

	// deprecated_is_odr used to be how a runner indicated if it was an ODR type runner.
	// Superseded by the ODR kind (field 5)
	DeprecatedIsOdr bool `json:"deprecated_is_odr,omitempty"`

	// The timestamps store the time this runner was first seen and the time
	// the runner was last seen. These values can be the same if the runner
	// was seen exactly once. The values are updated only when a runner starts
	// up.
	// Format: date-time
	FirstSeen strfmt.DateTime `json:"first_seen,omitempty"`

	// id is a unique ID generated by the runner. This should be a UUID or some
	// other guaranteed unique mechanism. This is not an auth mechanism, just
	// a way to associate an ID to a runner.
	ID string `json:"id,omitempty"`

	// Labels for the runner. These are the same as labels for any other
	// system in Waypoint (see operations such as Deployment). For runners, they
	// can additionally be used as a targeting mechanism.
	Labels map[string]string `json:"labels,omitempty"`

	// last seen
	// Format: date-time
	LastSeen strfmt.DateTime `json:"last_seen,omitempty"`

	// local indicates this runner was created by a cli instantiation
	Local HashicorpWaypointRunnerLocal `json:"local,omitempty"`

	// odr is set if this runner as an on-demand runner. For ODRs, we expect
	// they will accept exactly one job and then exit. This is used by the
	// server to change some other behavior:
	//
	// * The server will give ODRs project/app-scoped config if it exists.
	//   * The server will never assign more than one job to this runner.
	//     This is also enforced in the runner client-side but the server also
	//     does this out of caution.
	Odr *HashicorpWaypointRunnerODR `json:"odr,omitempty"`

	// True if this runner is currently online and connected.
	Online bool `json:"online,omitempty"`

	// remote indicates this is a "static" remote runner
	Remote HashicorpWaypointRunnerRemote `json:"remote,omitempty"`
}

// Validate validates this hashicorp waypoint runner
func (m *HashicorpWaypointRunner) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAdoptionState(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateComponents(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFirstSeen(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastSeen(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOdr(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointRunner) validateAdoptionState(formats strfmt.Registry) error {
	if swag.IsZero(m.AdoptionState) { // not required
		return nil
	}

	if m.AdoptionState != nil {
		if err := m.AdoptionState.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("adoption_state")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("adoption_state")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointRunner) validateComponents(formats strfmt.Registry) error {
	if swag.IsZero(m.Components) { // not required
		return nil
	}

	for i := 0; i < len(m.Components); i++ {
		if swag.IsZero(m.Components[i]) { // not required
			continue
		}

		if m.Components[i] != nil {
			if err := m.Components[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("components" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("components" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointRunner) validateFirstSeen(formats strfmt.Registry) error {
	if swag.IsZero(m.FirstSeen) { // not required
		return nil
	}

	if err := validate.FormatOf("first_seen", "body", "date-time", m.FirstSeen.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *HashicorpWaypointRunner) validateLastSeen(formats strfmt.Registry) error {
	if swag.IsZero(m.LastSeen) { // not required
		return nil
	}

	if err := validate.FormatOf("last_seen", "body", "date-time", m.LastSeen.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *HashicorpWaypointRunner) validateOdr(formats strfmt.Registry) error {
	if swag.IsZero(m.Odr) { // not required
		return nil
	}

	if m.Odr != nil {
		if err := m.Odr.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("odr")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("odr")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint runner based on the context it is used
func (m *HashicorpWaypointRunner) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAdoptionState(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateComponents(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateOdr(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointRunner) contextValidateAdoptionState(ctx context.Context, formats strfmt.Registry) error {

	if m.AdoptionState != nil {
		if err := m.AdoptionState.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("adoption_state")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("adoption_state")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointRunner) contextValidateComponents(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Components); i++ {

		if m.Components[i] != nil {
			if err := m.Components[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("components" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("components" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *HashicorpWaypointRunner) contextValidateOdr(ctx context.Context, formats strfmt.Registry) error {

	if m.Odr != nil {
		if err := m.Odr.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("odr")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("odr")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointRunner) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointRunner) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointRunner
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
