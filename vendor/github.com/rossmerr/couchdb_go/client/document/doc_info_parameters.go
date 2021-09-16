// Code generated by go-swagger; DO NOT EDIT.

package document

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

// NewDocInfoParams creates a new DocInfoParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDocInfoParams() *DocInfoParams {
	return &DocInfoParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDocInfoParamsWithTimeout creates a new DocInfoParams object
// with the ability to set a timeout on a request.
func NewDocInfoParamsWithTimeout(timeout time.Duration) *DocInfoParams {
	return &DocInfoParams{
		timeout: timeout,
	}
}

// NewDocInfoParamsWithContext creates a new DocInfoParams object
// with the ability to set a context for a request.
func NewDocInfoParamsWithContext(ctx context.Context) *DocInfoParams {
	return &DocInfoParams{
		Context: ctx,
	}
}

// NewDocInfoParamsWithHTTPClient creates a new DocInfoParams object
// with the ability to set a custom HTTPClient for a request.
func NewDocInfoParamsWithHTTPClient(client *http.Client) *DocInfoParams {
	return &DocInfoParams{
		HTTPClient: client,
	}
}

/* DocInfoParams contains all the parameters to send to the API endpoint
   for the doc info operation.

   Typically these are written to a http.Request.
*/
type DocInfoParams struct {

	/* Db.

	   Database name
	*/
	Db string

	/* Docid.

	   DDocument ID
	*/
	Docid string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the doc info params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DocInfoParams) WithDefaults() *DocInfoParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the doc info params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DocInfoParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the doc info params
func (o *DocInfoParams) WithTimeout(timeout time.Duration) *DocInfoParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the doc info params
func (o *DocInfoParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the doc info params
func (o *DocInfoParams) WithContext(ctx context.Context) *DocInfoParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the doc info params
func (o *DocInfoParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the doc info params
func (o *DocInfoParams) WithHTTPClient(client *http.Client) *DocInfoParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the doc info params
func (o *DocInfoParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDb adds the db to the doc info params
func (o *DocInfoParams) WithDb(db string) *DocInfoParams {
	o.SetDb(db)
	return o
}

// SetDb adds the db to the doc info params
func (o *DocInfoParams) SetDb(db string) {
	o.Db = db
}

// WithDocid adds the docid to the doc info params
func (o *DocInfoParams) WithDocid(docid string) *DocInfoParams {
	o.SetDocid(docid)
	return o
}

// SetDocid adds the docid to the doc info params
func (o *DocInfoParams) SetDocid(docid string) {
	o.Docid = docid
}

// WriteToRequest writes these params to a swagger request
func (o *DocInfoParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param db
	if err := r.SetPathParam("db", o.Db); err != nil {
		return err
	}

	// path param docid
	if err := r.SetPathParam("docid", o.Docid); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
