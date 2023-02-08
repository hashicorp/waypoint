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
)

// NewWaypointGetRunnerParams creates a new WaypointGetRunnerParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWaypointGetRunnerParams() *WaypointGetRunnerParams {
	return &WaypointGetRunnerParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewWaypointGetRunnerParamsWithTimeout creates a new WaypointGetRunnerParams object
// with the ability to set a timeout on a request.
func NewWaypointGetRunnerParamsWithTimeout(timeout time.Duration) *WaypointGetRunnerParams {
	return &WaypointGetRunnerParams{
		timeout: timeout,
	}
}

// NewWaypointGetRunnerParamsWithContext creates a new WaypointGetRunnerParams object
// with the ability to set a context for a request.
func NewWaypointGetRunnerParamsWithContext(ctx context.Context) *WaypointGetRunnerParams {
	return &WaypointGetRunnerParams{
		Context: ctx,
	}
}

// NewWaypointGetRunnerParamsWithHTTPClient creates a new WaypointGetRunnerParams object
// with the ability to set a custom HTTPClient for a request.
func NewWaypointGetRunnerParamsWithHTTPClient(client *http.Client) *WaypointGetRunnerParams {
	return &WaypointGetRunnerParams{
		HTTPClient: client,
	}
}

/*
WaypointGetRunnerParams contains all the parameters to send to the API endpoint

	for the waypoint get runner operation.

	Typically these are written to a http.Request.
*/
type WaypointGetRunnerParams struct {

	/* RunnerID.

	   ID of the runner to request.
	*/
	RunnerID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the waypoint get runner params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetRunnerParams) WithDefaults() *WaypointGetRunnerParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the waypoint get runner params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetRunnerParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the waypoint get runner params
func (o *WaypointGetRunnerParams) WithTimeout(timeout time.Duration) *WaypointGetRunnerParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the waypoint get runner params
func (o *WaypointGetRunnerParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the waypoint get runner params
func (o *WaypointGetRunnerParams) WithContext(ctx context.Context) *WaypointGetRunnerParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the waypoint get runner params
func (o *WaypointGetRunnerParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the waypoint get runner params
func (o *WaypointGetRunnerParams) WithHTTPClient(client *http.Client) *WaypointGetRunnerParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the waypoint get runner params
func (o *WaypointGetRunnerParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRunnerID adds the runnerID to the waypoint get runner params
func (o *WaypointGetRunnerParams) WithRunnerID(runnerID string) *WaypointGetRunnerParams {
	o.SetRunnerID(runnerID)
	return o
}

// SetRunnerID adds the runnerId to the waypoint get runner params
func (o *WaypointGetRunnerParams) SetRunnerID(runnerID string) {
	o.RunnerID = runnerID
}

// WriteToRequest writes these params to a swagger request
func (o *WaypointGetRunnerParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param runner_id
	if err := r.SetPathParam("runner_id", o.RunnerID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
