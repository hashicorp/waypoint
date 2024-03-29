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

// NewWaypointGetLatestPushedArtifactParams creates a new WaypointGetLatestPushedArtifactParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWaypointGetLatestPushedArtifactParams() *WaypointGetLatestPushedArtifactParams {
	return &WaypointGetLatestPushedArtifactParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewWaypointGetLatestPushedArtifactParamsWithTimeout creates a new WaypointGetLatestPushedArtifactParams object
// with the ability to set a timeout on a request.
func NewWaypointGetLatestPushedArtifactParamsWithTimeout(timeout time.Duration) *WaypointGetLatestPushedArtifactParams {
	return &WaypointGetLatestPushedArtifactParams{
		timeout: timeout,
	}
}

// NewWaypointGetLatestPushedArtifactParamsWithContext creates a new WaypointGetLatestPushedArtifactParams object
// with the ability to set a context for a request.
func NewWaypointGetLatestPushedArtifactParamsWithContext(ctx context.Context) *WaypointGetLatestPushedArtifactParams {
	return &WaypointGetLatestPushedArtifactParams{
		Context: ctx,
	}
}

// NewWaypointGetLatestPushedArtifactParamsWithHTTPClient creates a new WaypointGetLatestPushedArtifactParams object
// with the ability to set a custom HTTPClient for a request.
func NewWaypointGetLatestPushedArtifactParamsWithHTTPClient(client *http.Client) *WaypointGetLatestPushedArtifactParams {
	return &WaypointGetLatestPushedArtifactParams{
		HTTPClient: client,
	}
}

/*
WaypointGetLatestPushedArtifactParams contains all the parameters to send to the API endpoint

	for the waypoint get latest pushed artifact operation.

	Typically these are written to a http.Request.
*/
type WaypointGetLatestPushedArtifactParams struct {

	// ApplicationApplication.
	ApplicationApplication string

	// ApplicationProject.
	ApplicationProject string

	// WorkspaceWorkspace.
	WorkspaceWorkspace *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the waypoint get latest pushed artifact params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetLatestPushedArtifactParams) WithDefaults() *WaypointGetLatestPushedArtifactParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the waypoint get latest pushed artifact params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointGetLatestPushedArtifactParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) WithTimeout(timeout time.Duration) *WaypointGetLatestPushedArtifactParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) WithContext(ctx context.Context) *WaypointGetLatestPushedArtifactParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) WithHTTPClient(client *http.Client) *WaypointGetLatestPushedArtifactParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApplicationApplication adds the applicationApplication to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) WithApplicationApplication(applicationApplication string) *WaypointGetLatestPushedArtifactParams {
	o.SetApplicationApplication(applicationApplication)
	return o
}

// SetApplicationApplication adds the applicationApplication to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) SetApplicationApplication(applicationApplication string) {
	o.ApplicationApplication = applicationApplication
}

// WithApplicationProject adds the applicationProject to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) WithApplicationProject(applicationProject string) *WaypointGetLatestPushedArtifactParams {
	o.SetApplicationProject(applicationProject)
	return o
}

// SetApplicationProject adds the applicationProject to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) SetApplicationProject(applicationProject string) {
	o.ApplicationProject = applicationProject
}

// WithWorkspaceWorkspace adds the workspaceWorkspace to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) WithWorkspaceWorkspace(workspaceWorkspace *string) *WaypointGetLatestPushedArtifactParams {
	o.SetWorkspaceWorkspace(workspaceWorkspace)
	return o
}

// SetWorkspaceWorkspace adds the workspaceWorkspace to the waypoint get latest pushed artifact params
func (o *WaypointGetLatestPushedArtifactParams) SetWorkspaceWorkspace(workspaceWorkspace *string) {
	o.WorkspaceWorkspace = workspaceWorkspace
}

// WriteToRequest writes these params to a swagger request
func (o *WaypointGetLatestPushedArtifactParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param application.application
	if err := r.SetPathParam("application.application", o.ApplicationApplication); err != nil {
		return err
	}

	// path param application.project
	if err := r.SetPathParam("application.project", o.ApplicationProject); err != nil {
		return err
	}

	if o.WorkspaceWorkspace != nil {

		// query param workspace.workspace
		var qrWorkspaceWorkspace string

		if o.WorkspaceWorkspace != nil {
			qrWorkspaceWorkspace = *o.WorkspaceWorkspace
		}
		qWorkspaceWorkspace := qrWorkspaceWorkspace
		if qWorkspaceWorkspace != "" {

			if err := r.SetQueryParam("workspace.workspace", qWorkspaceWorkspace); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
