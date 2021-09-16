// Code generated by go-swagger; DO NOT EDIT.

package document

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/rossmerr/couchdb_go/models"
)

// DocDeleteReader is a Reader for the DocDelete structure.
type DocDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DocDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDocDeleteOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 202:
		result := NewDocDeleteAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDocDeleteBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewDocDeleteUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDocDeleteNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewDocDeleteConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDocDeleteOK creates a DocDeleteOK with default headers values
func NewDocDeleteOK() *DocDeleteOK {
	return &DocDeleteOK{}
}

/* DocDeleteOK describes a response with status code 200, with default header values.

Document successfully removed
*/
type DocDeleteOK struct {

	/* Double quoted document’s revision token
	 */
	ETag string

	Payload *models.DocumentOK
}

func (o *DocDeleteOK) Error() string {
	return fmt.Sprintf("[DELETE /{db}/{docid}][%d] docDeleteOK  %+v", 200, o.Payload)
}
func (o *DocDeleteOK) GetPayload() *models.DocumentOK {
	return o.Payload
}

func (o *DocDeleteOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header ETag
	hdrETag := response.GetHeader("ETag")

	if hdrETag != "" {
		o.ETag = hdrETag
	}

	o.Payload = new(models.DocumentOK)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDocDeleteAccepted creates a DocDeleteAccepted with default headers values
func NewDocDeleteAccepted() *DocDeleteAccepted {
	return &DocDeleteAccepted{}
}

/* DocDeleteAccepted describes a response with status code 202, with default header values.

Request was accepted, but changes are not yet stored on disk
*/
type DocDeleteAccepted struct {

	/* Double quoted document’s revision token
	 */
	ETag string

	Payload *models.DocumentOK
}

func (o *DocDeleteAccepted) Error() string {
	return fmt.Sprintf("[DELETE /{db}/{docid}][%d] docDeleteAccepted  %+v", 202, o.Payload)
}
func (o *DocDeleteAccepted) GetPayload() *models.DocumentOK {
	return o.Payload
}

func (o *DocDeleteAccepted) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header ETag
	hdrETag := response.GetHeader("ETag")

	if hdrETag != "" {
		o.ETag = hdrETag
	}

	o.Payload = new(models.DocumentOK)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDocDeleteBadRequest creates a DocDeleteBadRequest with default headers values
func NewDocDeleteBadRequest() *DocDeleteBadRequest {
	return &DocDeleteBadRequest{}
}

/* DocDeleteBadRequest describes a response with status code 400, with default header values.

Invalid request body or parameters
*/
type DocDeleteBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *DocDeleteBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /{db}/{docid}][%d] docDeleteBadRequest  %+v", 400, o.Payload)
}
func (o *DocDeleteBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DocDeleteBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDocDeleteUnauthorized creates a DocDeleteUnauthorized with default headers values
func NewDocDeleteUnauthorized() *DocDeleteUnauthorized {
	return &DocDeleteUnauthorized{}
}

/* DocDeleteUnauthorized describes a response with status code 401, with default header values.

Write privileges required
*/
type DocDeleteUnauthorized struct {
	Payload *models.ErrorResponse
}

func (o *DocDeleteUnauthorized) Error() string {
	return fmt.Sprintf("[DELETE /{db}/{docid}][%d] docDeleteUnauthorized  %+v", 401, o.Payload)
}
func (o *DocDeleteUnauthorized) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DocDeleteUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDocDeleteNotFound creates a DocDeleteNotFound with default headers values
func NewDocDeleteNotFound() *DocDeleteNotFound {
	return &DocDeleteNotFound{}
}

/* DocDeleteNotFound describes a response with status code 404, with default header values.

Specified database or document ID doesn’t exists
*/
type DocDeleteNotFound struct {
	Payload *models.ErrorResponse
}

func (o *DocDeleteNotFound) Error() string {
	return fmt.Sprintf("[DELETE /{db}/{docid}][%d] docDeleteNotFound  %+v", 404, o.Payload)
}
func (o *DocDeleteNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DocDeleteNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDocDeleteConflict creates a DocDeleteConflict with default headers values
func NewDocDeleteConflict() *DocDeleteConflict {
	return &DocDeleteConflict{}
}

/* DocDeleteConflict describes a response with status code 409, with default header values.

Specified revision is not the latest for target document
*/
type DocDeleteConflict struct {
	Payload *models.ErrorResponse
}

func (o *DocDeleteConflict) Error() string {
	return fmt.Sprintf("[DELETE /{db}/{docid}][%d] docDeleteConflict  %+v", 409, o.Payload)
}
func (o *DocDeleteConflict) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DocDeleteConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
