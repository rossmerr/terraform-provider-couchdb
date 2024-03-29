// Code generated by go-swagger; DO NOT EDIT.

package database

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/rossmerr/couchdb_go/models"
)

// DesignDocAllGetReader is a Reader for the DesignDocAllGet structure.
type DesignDocAllGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DesignDocAllGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDesignDocAllGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewDesignDocAllGetNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDesignDocAllGetOK creates a DesignDocAllGetOK with default headers values
func NewDesignDocAllGetOK() *DesignDocAllGetOK {
	return &DesignDocAllGetOK{}
}

/* DesignDocAllGetOK describes a response with status code 200, with default header values.

Request completed successfully
*/
type DesignDocAllGetOK struct {

	/* Response signature
	 */
	ETag string

	Payload *models.Pagination
}

func (o *DesignDocAllGetOK) Error() string {
	return fmt.Sprintf("[GET /{db}/_design_docs][%d] designDocAllGetOK  %+v", 200, o.Payload)
}
func (o *DesignDocAllGetOK) GetPayload() *models.Pagination {
	return o.Payload
}

func (o *DesignDocAllGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header ETag
	hdrETag := response.GetHeader("ETag")

	if hdrETag != "" {
		o.ETag = hdrETag
	}

	o.Payload = new(models.Pagination)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDesignDocAllGetNotFound creates a DesignDocAllGetNotFound with default headers values
func NewDesignDocAllGetNotFound() *DesignDocAllGetNotFound {
	return &DesignDocAllGetNotFound{}
}

/* DesignDocAllGetNotFound describes a response with status code 404, with default header values.

Requested database not found
*/
type DesignDocAllGetNotFound struct {
	Payload *models.ErrorResponse
}

func (o *DesignDocAllGetNotFound) Error() string {
	return fmt.Sprintf("[GET /{db}/_design_docs][%d] designDocAllGetNotFound  %+v", 404, o.Payload)
}
func (o *DesignDocAllGetNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DesignDocAllGetNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
