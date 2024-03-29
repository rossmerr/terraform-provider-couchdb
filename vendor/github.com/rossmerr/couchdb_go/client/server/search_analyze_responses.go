// Code generated by go-swagger; DO NOT EDIT.

package server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/rossmerr/couchdb_go/models"
)

// SearchAnalyzeReader is a Reader for the SearchAnalyze structure.
type SearchAnalyzeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SearchAnalyzeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSearchAnalyzeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewSearchAnalyzeBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewSearchAnalyzeInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewSearchAnalyzeOK creates a SearchAnalyzeOK with default headers values
func NewSearchAnalyzeOK() *SearchAnalyzeOK {
	return &SearchAnalyzeOK{}
}

/* SearchAnalyzeOK describes a response with status code 200, with default header values.

Request completed successfully
*/
type SearchAnalyzeOK struct {
	Payload *models.InlineResponse2003
}

func (o *SearchAnalyzeOK) Error() string {
	return fmt.Sprintf("[POST /_search_analyze][%d] searchAnalyzeOK  %+v", 200, o.Payload)
}
func (o *SearchAnalyzeOK) GetPayload() *models.InlineResponse2003 {
	return o.Payload
}

func (o *SearchAnalyzeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.InlineResponse2003)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSearchAnalyzeBadRequest creates a SearchAnalyzeBadRequest with default headers values
func NewSearchAnalyzeBadRequest() *SearchAnalyzeBadRequest {
	return &SearchAnalyzeBadRequest{}
}

/* SearchAnalyzeBadRequest describes a response with status code 400, with default header values.

Request body is wrong (malformed or missing one of the mandatory fields)
*/
type SearchAnalyzeBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *SearchAnalyzeBadRequest) Error() string {
	return fmt.Sprintf("[POST /_search_analyze][%d] searchAnalyzeBadRequest  %+v", 400, o.Payload)
}
func (o *SearchAnalyzeBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *SearchAnalyzeBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSearchAnalyzeInternalServerError creates a SearchAnalyzeInternalServerError with default headers values
func NewSearchAnalyzeInternalServerError() *SearchAnalyzeInternalServerError {
	return &SearchAnalyzeInternalServerError{}
}

/* SearchAnalyzeInternalServerError describes a response with status code 500, with default header values.

A server error (or other kind of error) occurred
*/
type SearchAnalyzeInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *SearchAnalyzeInternalServerError) Error() string {
	return fmt.Sprintf("[POST /_search_analyze][%d] searchAnalyzeInternalServerError  %+v", 500, o.Payload)
}
func (o *SearchAnalyzeInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *SearchAnalyzeInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
