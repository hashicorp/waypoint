// Code generated by go-swagger; DO NOT EDIT.

package waypoint

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/waypoint/pkg/client/gen/models"
)

// NewWaypointConvertInviteTokenParams creates a new WaypointConvertInviteTokenParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWaypointConvertInviteTokenParams() *WaypointConvertInviteTokenParams {
	return &WaypointConvertInviteTokenParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewWaypointConvertInviteTokenParamsWithTimeout creates a new WaypointConvertInviteTokenParams object
// with the ability to set a timeout on a request.
func NewWaypointConvertInviteTokenParamsWithTimeout(timeout time.Duration) *WaypointConvertInviteTokenParams {
	return &WaypointConvertInviteTokenParams{
		timeout: timeout,
	}
}

// NewWaypointConvertInviteTokenParamsWithContext creates a new WaypointConvertInviteTokenParams object
// with the ability to set a context for a request.
func NewWaypointConvertInviteTokenParamsWithContext(ctx context.Context) *WaypointConvertInviteTokenParams {
	return &WaypointConvertInviteTokenParams{
		Context: ctx,
	}
}

// NewWaypointConvertInviteTokenParamsWithHTTPClient creates a new WaypointConvertInviteTokenParams object
// with the ability to set a custom HTTPClient for a request.
func NewWaypointConvertInviteTokenParamsWithHTTPClient(client *http.Client) *WaypointConvertInviteTokenParams {
	return &WaypointConvertInviteTokenParams{
		HTTPClient: client,
	}
}

/*
WaypointConvertInviteTokenParams contains all the parameters to send to the API endpoint

	for the waypoint convert invite token operation.

	Typically these are written to a http.Request.
*/
type WaypointConvertInviteTokenParams struct {

	// Body.
	Body *models.HashicorpWaypointConvertInviteTokenRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the waypoint convert invite token params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointConvertInviteTokenParams) WithDefaults() *WaypointConvertInviteTokenParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the waypoint convert invite token params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointConvertInviteTokenParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the waypoint convert invite token params
func (o *WaypointConvertInviteTokenParams) WithTimeout(timeout time.Duration) *WaypointConvertInviteTokenParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the waypoint convert invite token params
func (o *WaypointConvertInviteTokenParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the waypoint convert invite token params
func (o *WaypointConvertInviteTokenParams) WithContext(ctx context.Context) *WaypointConvertInviteTokenParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the waypoint convert invite token params
func (o *WaypointConvertInviteTokenParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the waypoint convert invite token params
func (o *WaypointConvertInviteTokenParams) WithHTTPClient(client *http.Client) *WaypointConvertInviteTokenParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the waypoint convert invite token params
func (o *WaypointConvertInviteTokenParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the waypoint convert invite token params
func (o *WaypointConvertInviteTokenParams) WithBody(body *models.HashicorpWaypointConvertInviteTokenRequest) *WaypointConvertInviteTokenParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the waypoint convert invite token params
func (o *WaypointConvertInviteTokenParams) SetBody(body *models.HashicorpWaypointConvertInviteTokenRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *WaypointConvertInviteTokenParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
