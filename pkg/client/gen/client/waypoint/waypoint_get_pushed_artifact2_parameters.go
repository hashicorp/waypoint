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

// NewWaypointGetPushedArtifact2Params creates a new WaypointGetPushedArtifact2Params object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWaypointGetPushedArtifact2Params() *WaypointGetPushedArtifact2Params {
	return &WaypointGetPushedArtifact2Params{
		timeout: cr.DefaultTimeout,
	}
}

// NewWaypointGetPushedArtifact2ParamsWithTimeout creates a new WaypointGetPushedArtifact2Params object
// with the ability to set a timeout on a request.
func NewWaypointGetPushedArtifact2ParamsWithTimeout(timeout time.Duration) *WaypointGetPushedArtifact2Params {
	return &WaypointGetPushedArtifact2Params{
		timeout: timeout,
	}
}

// NewWaypointGetPushedArtifact2ParamsWithContext creates a new WaypointGetPushedArtifact2Params object
// with the ability to set a context for a request.
func NewWaypointGetPushedArtifact2ParamsWithContext(ctx context.Context) *WaypointGetPushedArtifact2Params {
	return &WaypointGetPushedArtifact2Params{
		Context: ctx,
	}
}

// NewWaypointGetPushedArtifact2ParamsWithHTTPClient creates a new WaypointGetPushedArtifact2Params object
// with the ability to set a custom HTTPClient for a request.
func NewWaypointGetPushedArtifact2ParamsWithHTTPClient(client *http.Client) *WaypointGetPushedArtifact2Params {
	return &WaypointGetPushedArtifact2Params{
		HTTPClient: client,
	}
}

/*
WaypointGetPushedArtifact2Params contains all the parameters to send to the API endpoint

	for the waypoint get pushed artifact2 operation.

	Typically these are written to a http.Request.
*/
type WaypointGetPushedArtifact2Params struct {

	// RefID.
	RefID *string

	// RefSequenceApplicationApplication.
	RefSequenceApplicationApplication string

	// RefSequenceApplicationProject.
	RefSequenceApplicationProject string

	// RefSequenceNumber.
	//
	// Format: uint64
	RefSequenceNumber string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the waypoint get pushed artifact2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetPushedArtifact2Params) WithDefaults() *WaypointGetPushedArtifact2Params {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the waypoint get pushed artifact2 params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetPushedArtifact2Params) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) WithTimeout(timeout time.Duration) *WaypointGetPushedArtifact2Params {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) WithContext(ctx context.Context) *WaypointGetPushedArtifact2Params {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) WithHTTPClient(client *http.Client) *WaypointGetPushedArtifact2Params {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRefID adds the refID to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) WithRefID(refID *string) *WaypointGetPushedArtifact2Params {
	o.SetRefID(refID)
	return o
}

// SetRefID adds the refId to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) SetRefID(refID *string) {
	o.RefID = refID
}

// WithRefSequenceApplicationApplication adds the refSequenceApplicationApplication to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) WithRefSequenceApplicationApplication(refSequenceApplicationApplication string) *WaypointGetPushedArtifact2Params {
	o.SetRefSequenceApplicationApplication(refSequenceApplicationApplication)
	return o
}

// SetRefSequenceApplicationApplication adds the refSequenceApplicationApplication to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) SetRefSequenceApplicationApplication(refSequenceApplicationApplication string) {
	o.RefSequenceApplicationApplication = refSequenceApplicationApplication
}

// WithRefSequenceApplicationProject adds the refSequenceApplicationProject to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) WithRefSequenceApplicationProject(refSequenceApplicationProject string) *WaypointGetPushedArtifact2Params {
	o.SetRefSequenceApplicationProject(refSequenceApplicationProject)
	return o
}

// SetRefSequenceApplicationProject adds the refSequenceApplicationProject to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) SetRefSequenceApplicationProject(refSequenceApplicationProject string) {
	o.RefSequenceApplicationProject = refSequenceApplicationProject
}

// WithRefSequenceNumber adds the refSequenceNumber to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) WithRefSequenceNumber(refSequenceNumber string) *WaypointGetPushedArtifact2Params {
	o.SetRefSequenceNumber(refSequenceNumber)
	return o
}

// SetRefSequenceNumber adds the refSequenceNumber to the waypoint get pushed artifact2 params
func (o *WaypointGetPushedArtifact2Params) SetRefSequenceNumber(refSequenceNumber string) {
	o.RefSequenceNumber = refSequenceNumber
}

// WriteToRequest writes these params to a swagger request
func (o *WaypointGetPushedArtifact2Params) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.RefID != nil {

		// query param ref.id
		var qrRefID string

		if o.RefID != nil {
			qrRefID = *o.RefID
		}
		qRefID := qrRefID
		if qRefID != "" {

			if err := r.SetQueryParam("ref.id", qRefID); err != nil {
				return err
			}
		}
	}

	// path param ref.sequence.application.application
	if err := r.SetPathParam("ref.sequence.application.application", o.RefSequenceApplicationApplication); err != nil {
		return err
	}

	// path param ref.sequence.application.project
	if err := r.SetPathParam("ref.sequence.application.project", o.RefSequenceApplicationProject); err != nil {
		return err
	}

	// path param ref.sequence.number
	if err := r.SetPathParam("ref.sequence.number", o.RefSequenceNumber); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
