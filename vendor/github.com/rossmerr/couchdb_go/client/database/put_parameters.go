// Code generated by go-swagger; DO NOT EDIT.

package database

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
	"github.com/go-openapi/swag"
)

// NewPutParams creates a new PutParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPutParams() *PutParams {
	return &PutParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPutParamsWithTimeout creates a new PutParams object
// with the ability to set a timeout on a request.
func NewPutParamsWithTimeout(timeout time.Duration) *PutParams {
	return &PutParams{
		timeout: timeout,
	}
}

// NewPutParamsWithContext creates a new PutParams object
// with the ability to set a context for a request.
func NewPutParamsWithContext(ctx context.Context) *PutParams {
	return &PutParams{
		Context: ctx,
	}
}

// NewPutParamsWithHTTPClient creates a new PutParams object
// with the ability to set a custom HTTPClient for a request.
func NewPutParamsWithHTTPClient(client *http.Client) *PutParams {
	return &PutParams{
		HTTPClient: client,
	}
}

/* PutParams contains all the parameters to send to the API endpoint
   for the put operation.

   Typically these are written to a http.Request.
*/
type PutParams struct {

	/* Db.

	   Database name
	*/
	Db string

	/* N.

	   Replicas. The number of copies of the database in the cluster. The default is 3, unless overridden in the cluster config .

	   Format: int32
	*/
	N *int32

	/* Partitioned.

	   Whether to create a partitioned database. Default is false.
	*/
	Partitioned *bool

	/* Q.

	   Shards, aka the number of range partitions. Default is 8, unless overridden in the cluster config.

	   Format: int32
	*/
	Q *int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the put params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutParams) WithDefaults() *PutParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the put params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the put params
func (o *PutParams) WithTimeout(timeout time.Duration) *PutParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put params
func (o *PutParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put params
func (o *PutParams) WithContext(ctx context.Context) *PutParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put params
func (o *PutParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put params
func (o *PutParams) WithHTTPClient(client *http.Client) *PutParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put params
func (o *PutParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDb adds the db to the put params
func (o *PutParams) WithDb(db string) *PutParams {
	o.SetDb(db)
	return o
}

// SetDb adds the db to the put params
func (o *PutParams) SetDb(db string) {
	o.Db = db
}

// WithN adds the n to the put params
func (o *PutParams) WithN(n *int32) *PutParams {
	o.SetN(n)
	return o
}

// SetN adds the n to the put params
func (o *PutParams) SetN(n *int32) {
	o.N = n
}

// WithPartitioned adds the partitioned to the put params
func (o *PutParams) WithPartitioned(partitioned *bool) *PutParams {
	o.SetPartitioned(partitioned)
	return o
}

// SetPartitioned adds the partitioned to the put params
func (o *PutParams) SetPartitioned(partitioned *bool) {
	o.Partitioned = partitioned
}

// WithQ adds the q to the put params
func (o *PutParams) WithQ(q *int32) *PutParams {
	o.SetQ(q)
	return o
}

// SetQ adds the q to the put params
func (o *PutParams) SetQ(q *int32) {
	o.Q = q
}

// WriteToRequest writes these params to a swagger request
func (o *PutParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param db
	if err := r.SetPathParam("db", o.Db); err != nil {
		return err
	}

	if o.N != nil {

		// query param n
		var qrN int32

		if o.N != nil {
			qrN = *o.N
		}
		qN := swag.FormatInt32(qrN)
		if qN != "" {

			if err := r.SetQueryParam("n", qN); err != nil {
				return err
			}
		}
	}

	if o.Partitioned != nil {

		// query param partitioned
		var qrPartitioned bool

		if o.Partitioned != nil {
			qrPartitioned = *o.Partitioned
		}
		qPartitioned := swag.FormatBool(qrPartitioned)
		if qPartitioned != "" {

			if err := r.SetQueryParam("partitioned", qPartitioned); err != nil {
				return err
			}
		}
	}

	if o.Q != nil {

		// query param q
		var qrQ int32

		if o.Q != nil {
			qrQ = *o.Q
		}
		qQ := swag.FormatInt32(qrQ)
		if qQ != "" {

			if err := r.SetQueryParam("q", qQ); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
