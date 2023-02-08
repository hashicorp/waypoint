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

// NewWaypointGenerateRunnerTokenParams creates a new WaypointGenerateRunnerTokenParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWaypointGenerateRunnerTokenParams() *WaypointGenerateRunnerTokenParams {
	return &WaypointGenerateRunnerTokenParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewWaypointGenerateRunnerTokenParamsWithTimeout creates a new WaypointGenerateRunnerTokenParams object
// with the ability to set a timeout on a request.
func NewWaypointGenerateRunnerTokenParamsWithTimeout(timeout time.Duration) *WaypointGenerateRunnerTokenParams {
	return &WaypointGenerateRunnerTokenParams{
		timeout: timeout,
	}
}

// NewWaypointGenerateRunnerTokenParamsWithContext creates a new WaypointGenerateRunnerTokenParams object
// with the ability to set a context for a request.
func NewWaypointGenerateRunnerTokenParamsWithContext(ctx context.Context) *WaypointGenerateRunnerTokenParams {
	return &WaypointGenerateRunnerTokenParams{
		Context: ctx,
	}
}

// NewWaypointGenerateRunnerTokenParamsWithHTTPClient creates a new WaypointGenerateRunnerTokenParams object
// with the ability to set a custom HTTPClient for a request.
func NewWaypointGenerateRunnerTokenParamsWithHTTPClient(client *http.Client) *WaypointGenerateRunnerTokenParams {
	return &WaypointGenerateRunnerTokenParams{
		HTTPClient: client,
	}
}

/*
WaypointGenerateRunnerTokenParams contains all the parameters to send to the API endpoint

	for the waypoint generate runner token operation.

	Typically these are written to a http.Request.
*/
type WaypointGenerateRunnerTokenParams struct {

	// Body.
	Body *models.HashicorpWaypointGenerateRunnerTokenRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the waypoint generate runner token params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGenerateRunnerTokenParams) WithDefaults() *WaypointGenerateRunnerTokenParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the waypoint generate runner token params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGenerateRunnerTokenParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the waypoint generate runner token params
func (o *WaypointGenerateRunnerTokenParams) WithTimeout(timeout time.Duration) *WaypointGenerateRunnerTokenParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the waypoint generate runner token params
func (o *WaypointGenerateRunnerTokenParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the waypoint generate runner token params
func (o *WaypointGenerateRunnerTokenParams) WithContext(ctx context.Context) *WaypointGenerateRunnerTokenParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the waypoint generate runner token params
func (o *WaypointGenerateRunnerTokenParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the waypoint generate runner token params
func (o *WaypointGenerateRunnerTokenParams) WithHTTPClient(client *http.Client) *WaypointGenerateRunnerTokenParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the waypoint generate runner token params
func (o *WaypointGenerateRunnerTokenParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the waypoint generate runner token params
func (o *WaypointGenerateRunnerTokenParams) WithBody(body *models.HashicorpWaypointGenerateRunnerTokenRequest) *WaypointGenerateRunnerTokenParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the waypoint generate runner token params
func (o *WaypointGenerateRunnerTokenParams) SetBody(body *models.HashicorpWaypointGenerateRunnerTokenRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *WaypointGenerateRunnerTokenParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
