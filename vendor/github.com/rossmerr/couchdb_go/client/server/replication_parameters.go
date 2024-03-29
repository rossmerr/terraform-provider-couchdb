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

	"github.com/rossmerr/couchdb_go/models"
)

// NewReplicationParams creates a new ReplicationParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewReplicationParams() *ReplicationParams {
	return &ReplicationParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewReplicationParamsWithTimeout creates a new ReplicationParams object
// with the ability to set a timeout on a request.
func NewReplicationParamsWithTimeout(timeout time.Duration) *ReplicationParams {
	return &ReplicationParams{
		timeout: timeout,
	}
}

// NewReplicationParamsWithContext creates a new ReplicationParams object
// with the ability to set a context for a request.
func NewReplicationParamsWithContext(ctx context.Context) *ReplicationParams {
	return &ReplicationParams{
		Context: ctx,
	}
}

// NewReplicationParamsWithHTTPClient creates a new ReplicationParams object
// with the ability to set a custom HTTPClient for a request.
func NewReplicationParamsWithHTTPClient(client *http.Client) *ReplicationParams {
	return &ReplicationParams{
		HTTPClient: client,
	}
}

/* ReplicationParams contains all the parameters to send to the API endpoint
   for the replication operation.

   Typically these are written to a http.Request.
*/
type ReplicationParams struct {

	// Body.
	Body *models.Replicate

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the replication params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ReplicationParams) WithDefaults() *ReplicationParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the replication params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ReplicationParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the replication params
func (o *ReplicationParams) WithTimeout(timeout time.Duration) *ReplicationParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the replication params
func (o *ReplicationParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the replication params
func (o *ReplicationParams) WithContext(ctx context.Context) *ReplicationParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the replication params
func (o *ReplicationParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the replication params
func (o *ReplicationParams) WithHTTPClient(client *http.Client) *ReplicationParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the replication params
func (o *ReplicationParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the replication params
func (o *ReplicationParams) WithBody(body *models.Replicate) *ReplicationParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the replication params
func (o *ReplicationParams) SetBody(body *models.Replicate) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *ReplicationParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
