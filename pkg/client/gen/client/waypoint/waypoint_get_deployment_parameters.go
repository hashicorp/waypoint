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

// NewWaypointGetDeploymentParams creates a new WaypointGetDeploymentParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWaypointGetDeploymentParams() *WaypointGetDeploymentParams {
	return &WaypointGetDeploymentParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewWaypointGetDeploymentParamsWithTimeout creates a new WaypointGetDeploymentParams object
// with the ability to set a timeout on a request.
func NewWaypointGetDeploymentParamsWithTimeout(timeout time.Duration) *WaypointGetDeploymentParams {
	return &WaypointGetDeploymentParams{
		timeout: timeout,
	}
}

// NewWaypointGetDeploymentParamsWithContext creates a new WaypointGetDeploymentParams object
// with the ability to set a context for a request.
func NewWaypointGetDeploymentParamsWithContext(ctx context.Context) *WaypointGetDeploymentParams {
	return &WaypointGetDeploymentParams{
		Context: ctx,
	}
}

// NewWaypointGetDeploymentParamsWithHTTPClient creates a new WaypointGetDeploymentParams object
// with the ability to set a custom HTTPClient for a request.
func NewWaypointGetDeploymentParamsWithHTTPClient(client *http.Client) *WaypointGetDeploymentParams {
	return &WaypointGetDeploymentParams{
		HTTPClient: client,
	}
}

/*
WaypointGetDeploymentParams contains all the parameters to send to the API endpoint

	for the waypoint get deployment operation.

	Typically these are written to a http.Request.
*/
type WaypointGetDeploymentParams struct {

	/* LoadDetails.

	     Indicate if the fetched deployments should include additional information
	about each deployment.

	     Default: "NONE"
	*/
	LoadDetails *string

	// RefID.
	RefID string

	// RefSequenceApplicationApplication.
	RefSequenceApplicationApplication *string

	// RefSequenceApplicationProject.
	RefSequenceApplicationProject *string

	// RefSequenceNumber.
	//
	// Format: uint64
	RefSequenceNumber *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the waypoint get deployment params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetDeploymentParams) WithDefaults() *WaypointGetDeploymentParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the waypoint get deployment params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetDeploymentParams) SetDefaults() {
	var (
		loadDetailsDefault = string("NONE")
	)

	val := WaypointGetDeploymentParams{
		LoadDetails: &loadDetailsDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) WithTimeout(timeout time.Duration) *WaypointGetDeploymentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) WithContext(ctx context.Context) *WaypointGetDeploymentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) WithHTTPClient(client *http.Client) *WaypointGetDeploymentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithLoadDetails adds the loadDetails to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) WithLoadDetails(loadDetails *string) *WaypointGetDeploymentParams {
	o.SetLoadDetails(loadDetails)
	return o
}

// SetLoadDetails adds the loadDetails to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) SetLoadDetails(loadDetails *string) {
	o.LoadDetails = loadDetails
}

// WithRefID adds the refID to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) WithRefID(refID string) *WaypointGetDeploymentParams {
	o.SetRefID(refID)
	return o
}

// SetRefID adds the refId to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) SetRefID(refID string) {
	o.RefID = refID
}

// WithRefSequenceApplicationApplication adds the refSequenceApplicationApplication to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) WithRefSequenceApplicationApplication(refSequenceApplicationApplication *string) *WaypointGetDeploymentParams {
	o.SetRefSequenceApplicationApplication(refSequenceApplicationApplication)
	return o
}

// SetRefSequenceApplicationApplication adds the refSequenceApplicationApplication to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) SetRefSequenceApplicationApplication(refSequenceApplicationApplication *string) {
	o.RefSequenceApplicationApplication = refSequenceApplicationApplication
}

// WithRefSequenceApplicationProject adds the refSequenceApplicationProject to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) WithRefSequenceApplicationProject(refSequenceApplicationProject *string) *WaypointGetDeploymentParams {
	o.SetRefSequenceApplicationProject(refSequenceApplicationProject)
	return o
}

// SetRefSequenceApplicationProject adds the refSequenceApplicationProject to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) SetRefSequenceApplicationProject(refSequenceApplicationProject *string) {
	o.RefSequenceApplicationProject = refSequenceApplicationProject
}

// WithRefSequenceNumber adds the refSequenceNumber to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) WithRefSequenceNumber(refSequenceNumber *string) *WaypointGetDeploymentParams {
	o.SetRefSequenceNumber(refSequenceNumber)
	return o
}

// SetRefSequenceNumber adds the refSequenceNumber to the waypoint get deployment params
func (o *WaypointGetDeploymentParams) SetRefSequenceNumber(refSequenceNumber *string) {
	o.RefSequenceNumber = refSequenceNumber
}

// WriteToRequest writes these params to a swagger request
func (o *WaypointGetDeploymentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.LoadDetails != nil {

		// query param load_details
		var qrLoadDetails string

		if o.LoadDetails != nil {
			qrLoadDetails = *o.LoadDetails
		}
		qLoadDetails := qrLoadDetails
		if qLoadDetails != "" {

			if err := r.SetQueryParam("load_details", qLoadDetails); err != nil {
				return err
			}
		}
	}

	// path param ref.id
	if err := r.SetPathParam("ref.id", o.RefID); err != nil {
		return err
	}

	if o.RefSequenceApplicationApplication != nil {

		// query param ref.sequence.application.application
		var qrRefSequenceApplicationApplication string

		if o.RefSequenceApplicationApplication != nil {
			qrRefSequenceApplicationApplication = *o.RefSequenceApplicationApplication
		}
		qRefSequenceApplicationApplication := qrRefSequenceApplicationApplication
		if qRefSequenceApplicationApplication != "" {

			if err := r.SetQueryParam("ref.sequence.application.application", qRefSequenceApplicationApplication); err != nil {
				return err
			}
		}
	}

	if o.RefSequenceApplicationProject != nil {

		// query param ref.sequence.application.project
		var qrRefSequenceApplicationProject string

		if o.RefSequenceApplicationProject != nil {
			qrRefSequenceApplicationProject = *o.RefSequenceApplicationProject
		}
		qRefSequenceApplicationProject := qrRefSequenceApplicationProject
		if qRefSequenceApplicationProject != "" {

			if err := r.SetQueryParam("ref.sequence.application.project", qRefSequenceApplicationProject); err != nil {
				return err
			}
		}
	}

	if o.RefSequenceNumber != nil {

		// query param ref.sequence.number
		var qrRefSequenceNumber string

		if o.RefSequenceNumber != nil {
			qrRefSequenceNumber = *o.RefSequenceNumber
		}
		qRefSequenceNumber := qrRefSequenceNumber
		if qRefSequenceNumber != "" {

			if err := r.SetQueryParam("ref.sequence.number", qRefSequenceNumber); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
