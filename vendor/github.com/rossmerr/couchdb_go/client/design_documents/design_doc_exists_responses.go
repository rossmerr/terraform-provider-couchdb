// Code generated by go-swagger; DO NOT EDIT.

package design_documents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/rossmerr/couchdb_go/models"
)

// DesignDocExistsReader is a Reader for the DesignDocExists structure.
type DesignDocExistsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DesignDocExistsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDesignDocExistsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 304:
		result := NewDesignDocExistsNotModified()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewDesignDocExistsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDesignDocExistsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDesignDocExistsOK creates a DesignDocExistsOK with default headers values
func NewDesignDocExistsOK() *DesignDocExistsOK {
	return &DesignDocExistsOK{}
}

/* DesignDocExistsOK describes a response with status code 200, with default header values.

Document exists
*/
type DesignDocExistsOK struct {

	/* Document size
	 */
	ContentLength int64

	/* Double quoted document’s revision token
	 */
	ETag string
}

func (o *DesignDocExistsOK) Error() string {
	return fmt.Sprintf("[HEAD /{db}/_design/{ddoc}][%d] designDocExistsOK ", 200)
}

func (o *DesignDocExistsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header Content-Length
	hdrContentLength := response.GetHeader("Content-Length")

	if hdrContentLength != "" {
		valcontentLength, err := swag.ConvertInt64(hdrContentLength)
		if err != nil {
			return errors.InvalidType("Content-Length", "header", "int64", hdrContentLength)
		}
		o.ContentLength = valcontentLength
	}

	// hydrates response header ETag
	hdrETag := response.GetHeader("ETag")

	if hdrETag != "" {
		o.ETag = hdrETag
	}

	return nil
}

// NewDesignDocExistsNotModified creates a DesignDocExistsNotModified with default headers values
func NewDesignDocExistsNotModified() *DesignDocExistsNotModified {
	return &DesignDocExistsNotModified{}
}

/* DesignDocExistsNotModified describes a response with status code 304, with default header values.

Document wasn’t modified since specified revision
*/
type DesignDocExistsNotModified struct {
}

func (o *DesignDocExistsNotModified) Error() string {
	return fmt.Sprintf("[HEAD /{db}/_design/{ddoc}][%d] designDocExistsNotModified ", 304)
}

func (o *DesignDocExistsNotModified) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDesignDocExistsUnauthorized creates a DesignDocExistsUnauthorized with default headers values
func NewDesignDocExistsUnauthorized() *DesignDocExistsUnauthorized {
	return &DesignDocExistsUnauthorized{}
}

/* DesignDocExistsUnauthorized describes a response with status code 401, with default header values.

Read privilege required
*/
type DesignDocExistsUnauthorized struct {
	Payload *models.ErrorResponse
}

func (o *DesignDocExistsUnauthorized) Error() string {
	return fmt.Sprintf("[HEAD /{db}/_design/{ddoc}][%d] designDocExistsUnauthorized  %+v", 401, o.Payload)
}
func (o *DesignDocExistsUnauthorized) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DesignDocExistsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDesignDocExistsNotFound creates a DesignDocExistsNotFound with default headers values
func NewDesignDocExistsNotFound() *DesignDocExistsNotFound {
	return &DesignDocExistsNotFound{}
}

/* DesignDocExistsNotFound describes a response with status code 404, with default header values.

Document not found
*/
type DesignDocExistsNotFound struct {
	Payload *models.ErrorResponse
}

func (o *DesignDocExistsNotFound) Error() string {
	return fmt.Sprintf("[HEAD /{db}/_design/{ddoc}][%d] designDocExistsNotFound  %+v", 404, o.Payload)
}
func (o *DesignDocExistsNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DesignDocExistsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
