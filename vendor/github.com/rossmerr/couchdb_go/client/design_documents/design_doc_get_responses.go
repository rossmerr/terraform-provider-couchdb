// Code generated by go-swagger; DO NOT EDIT.

package design_documents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/rossmerr/couchdb_go/models"
)

// DesignDocGetReader is a Reader for the DesignDocGet structure.
type DesignDocGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DesignDocGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDesignDocGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 304:
		result := NewDesignDocGetNotModified()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 400:
		result := NewDesignDocGetBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewDesignDocGetUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDesignDocGetNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDesignDocGetOK creates a DesignDocGetOK with default headers values
func NewDesignDocGetOK() *DesignDocGetOK {
	return &DesignDocGetOK{}
}

/* DesignDocGetOK describes a response with status code 200, with default header values.

Request completed successfully
*/
type DesignDocGetOK struct {

	/* Double quoted document’s revision token
	 */
	ETag string

	/* chunked. Available if requested with query parameter open_revs
	 */
	TransferEncoding string

	Payload *models.DesignDoc
}

func (o *DesignDocGetOK) Error() string {
	return fmt.Sprintf("[GET /{db}/_design/{ddoc}][%d] designDocGetOK  %+v", 200, o.Payload)
}
func (o *DesignDocGetOK) GetPayload() *models.DesignDoc {
	return o.Payload
}

func (o *DesignDocGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header ETag
	hdrETag := response.GetHeader("ETag")

	if hdrETag != "" {
		o.ETag = hdrETag
	}

	// hydrates response header Transfer-Encoding
	hdrTransferEncoding := response.GetHeader("Transfer-Encoding")

	if hdrTransferEncoding != "" {
		o.TransferEncoding = hdrTransferEncoding
	}

	o.Payload = new(models.DesignDoc)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDesignDocGetNotModified creates a DesignDocGetNotModified with default headers values
func NewDesignDocGetNotModified() *DesignDocGetNotModified {
	return &DesignDocGetNotModified{}
}

/* DesignDocGetNotModified describes a response with status code 304, with default header values.

Document wasn’t modified since specified revision
*/
type DesignDocGetNotModified struct {
}

func (o *DesignDocGetNotModified) Error() string {
	return fmt.Sprintf("[GET /{db}/_design/{ddoc}][%d] designDocGetNotModified ", 304)
}

func (o *DesignDocGetNotModified) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDesignDocGetBadRequest creates a DesignDocGetBadRequest with default headers values
func NewDesignDocGetBadRequest() *DesignDocGetBadRequest {
	return &DesignDocGetBadRequest{}
}

/* DesignDocGetBadRequest describes a response with status code 400, with default header values.

The format of the request or revision was invalid
*/
type DesignDocGetBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *DesignDocGetBadRequest) Error() string {
	return fmt.Sprintf("[GET /{db}/_design/{ddoc}][%d] designDocGetBadRequest  %+v", 400, o.Payload)
}
func (o *DesignDocGetBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DesignDocGetBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDesignDocGetUnauthorized creates a DesignDocGetUnauthorized with default headers values
func NewDesignDocGetUnauthorized() *DesignDocGetUnauthorized {
	return &DesignDocGetUnauthorized{}
}

/* DesignDocGetUnauthorized describes a response with status code 401, with default header values.

Read privilege required
*/
type DesignDocGetUnauthorized struct {
	Payload *models.ErrorResponse
}

func (o *DesignDocGetUnauthorized) Error() string {
	return fmt.Sprintf("[GET /{db}/_design/{ddoc}][%d] designDocGetUnauthorized  %+v", 401, o.Payload)
}
func (o *DesignDocGetUnauthorized) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DesignDocGetUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDesignDocGetNotFound creates a DesignDocGetNotFound with default headers values
func NewDesignDocGetNotFound() *DesignDocGetNotFound {
	return &DesignDocGetNotFound{}
}

/* DesignDocGetNotFound describes a response with status code 404, with default header values.

Document not found
*/
type DesignDocGetNotFound struct {
	Payload *models.ErrorResponse
}

func (o *DesignDocGetNotFound) Error() string {
	return fmt.Sprintf("[GET /{db}/_design/{ddoc}][%d] designDocGetNotFound  %+v", 404, o.Payload)
}
func (o *DesignDocGetNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DesignDocGetNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
