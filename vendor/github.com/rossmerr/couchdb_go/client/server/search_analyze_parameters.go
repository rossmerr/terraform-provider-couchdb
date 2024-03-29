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

// NewSearchAnalyzeParams creates a new SearchAnalyzeParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSearchAnalyzeParams() *SearchAnalyzeParams {
	return &SearchAnalyzeParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSearchAnalyzeParamsWithTimeout creates a new SearchAnalyzeParams object
// with the ability to set a timeout on a request.
func NewSearchAnalyzeParamsWithTimeout(timeout time.Duration) *SearchAnalyzeParams {
	return &SearchAnalyzeParams{
		timeout: timeout,
	}
}

// NewSearchAnalyzeParamsWithContext creates a new SearchAnalyzeParams object
// with the ability to set a context for a request.
func NewSearchAnalyzeParamsWithContext(ctx context.Context) *SearchAnalyzeParams {
	return &SearchAnalyzeParams{
		Context: ctx,
	}
}

// NewSearchAnalyzeParamsWithHTTPClient creates a new SearchAnalyzeParams object
// with the ability to set a custom HTTPClient for a request.
func NewSearchAnalyzeParamsWithHTTPClient(client *http.Client) *SearchAnalyzeParams {
	return &SearchAnalyzeParams{
		HTTPClient: client,
	}
}

/* SearchAnalyzeParams contains all the parameters to send to the API endpoint
   for the search analyze operation.

   Typically these are written to a http.Request.
*/
type SearchAnalyzeParams struct {

	// Body.
	Body *models.Body1

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the search analyze params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchAnalyzeParams) WithDefaults() *SearchAnalyzeParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the search analyze params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SearchAnalyzeParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the search analyze params
func (o *SearchAnalyzeParams) WithTimeout(timeout time.Duration) *SearchAnalyzeParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the search analyze params
func (o *SearchAnalyzeParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the search analyze params
func (o *SearchAnalyzeParams) WithContext(ctx context.Context) *SearchAnalyzeParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the search analyze params
func (o *SearchAnalyzeParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the search analyze params
func (o *SearchAnalyzeParams) WithHTTPClient(client *http.Client) *SearchAnalyzeParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the search analyze params
func (o *SearchAnalyzeParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the search analyze params
func (o *SearchAnalyzeParams) WithBody(body *models.Body1) *SearchAnalyzeParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the search analyze params
func (o *SearchAnalyzeParams) SetBody(body *models.Body1) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *SearchAnalyzeParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
