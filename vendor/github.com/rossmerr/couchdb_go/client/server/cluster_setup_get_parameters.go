// Code generated by go-swagger; DO NOT EDIT.

package server

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

// NewClusterSetupGetParams creates a new ClusterSetupGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewClusterSetupGetParams() *ClusterSetupGetParams {
	return &ClusterSetupGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewClusterSetupGetParamsWithTimeout creates a new ClusterSetupGetParams object
// with the ability to set a timeout on a request.
func NewClusterSetupGetParamsWithTimeout(timeout time.Duration) *ClusterSetupGetParams {
	return &ClusterSetupGetParams{
		timeout: timeout,
	}
}

// NewClusterSetupGetParamsWithContext creates a new ClusterSetupGetParams object
// with the ability to set a context for a request.
func NewClusterSetupGetParamsWithContext(ctx context.Context) *ClusterSetupGetParams {
	return &ClusterSetupGetParams{
		Context: ctx,
	}
}

// NewClusterSetupGetParamsWithHTTPClient creates a new ClusterSetupGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewClusterSetupGetParamsWithHTTPClient(client *http.Client) *ClusterSetupGetParams {
	return &ClusterSetupGetParams{
		HTTPClient: client,
	}
}

/* ClusterSetupGetParams contains all the parameters to send to the API endpoint
   for the cluster setup get operation.

   Typically these are written to a http.Request.
*/
type ClusterSetupGetParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the cluster setup get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ClusterSetupGetParams) WithDefaults() *ClusterSetupGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the cluster setup get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ClusterSetupGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the cluster setup get params
func (o *ClusterSetupGetParams) WithTimeout(timeout time.Duration) *ClusterSetupGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the cluster setup get params
func (o *ClusterSetupGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the cluster setup get params
func (o *ClusterSetupGetParams) WithContext(ctx context.Context) *ClusterSetupGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the cluster setup get params
func (o *ClusterSetupGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the cluster setup get params
func (o *ClusterSetupGetParams) WithHTTPClient(client *http.Client) *ClusterSetupGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the cluster setup get params
func (o *ClusterSetupGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *ClusterSetupGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
