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

// HashicorpWaypointGetJobStreamResponseTerminalEventTable hashicorp waypoint get job stream response terminal event table
//
// swagger:model hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table
type HashicorpWaypointGetJobStreamResponseTerminalEventTable struct {

	// headers
	Headers []string `json:"headers"`

	// rows
	Rows []*HashicorpWaypointGetJobStreamResponseTerminalEventTableRow `json:"rows"`
}

// Validate validates this hashicorp waypoint get job stream response terminal event table
func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTable) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRows(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTable) validateRows(formats strfmt.Registry) error {
	if swag.IsZero(m.Rows) { // not required
		return nil
	}

	for i := 0; i < len(m.Rows); i++ {
		if swag.IsZero(m.Rows[i]) { // not required
			continue
		}

		if m.Rows[i] != nil {
			if err := m.Rows[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("rows" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("rows" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this hashicorp waypoint get job stream response terminal event table based on the context it is used
func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTable) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateRows(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTable) contextValidateRows(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Rows); i++ {

		if m.Rows[i] != nil {
			if err := m.Rows[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("rows" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("rows" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTable) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTable) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointGetJobStreamResponseTerminalEventTable
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
