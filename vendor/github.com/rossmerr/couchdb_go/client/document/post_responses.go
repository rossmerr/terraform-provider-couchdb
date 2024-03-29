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

// PostReader is a Reader for the Post structure.
type PostReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPostCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 202:
		result := NewPostAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewPostUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewPostNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewPostConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPostCreated creates a PostCreated with default headers values
func NewPostCreated() *PostCreated {
	return &PostCreated{}
}

/* PostCreated describes a response with status code 201, with default header values.

Document created and stored on disk
*/
type PostCreated struct {
	Payload *models.DocumentOK
}

func (o *PostCreated) Error() string {
	return fmt.Sprintf("[POST /{db}][%d] postCreated  %+v", 201, o.Payload)
}
func (o *PostCreated) GetPayload() *models.DocumentOK {
	return o.Payload
}

func (o *PostCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DocumentOK)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostAccepted creates a PostAccepted with default headers values
func NewPostAccepted() *PostAccepted {
	return &PostAccepted{}
}

/* PostAccepted describes a response with status code 202, with default header values.

Document data accepted, but not yet stored on disk
*/
type PostAccepted struct {
	Payload *models.DocumentOK
}

func (o *PostAccepted) Error() string {
	return fmt.Sprintf("[POST /{db}][%d] postAccepted  %+v", 202, o.Payload)
}
func (o *PostAccepted) GetPayload() *models.DocumentOK {
	return o.Payload
}

func (o *PostAccepted) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DocumentOK)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostBadRequest creates a PostBadRequest with default headers values
func NewPostBadRequest() *PostBadRequest {
	return &PostBadRequest{}
}

/* PostBadRequest describes a response with status code 400, with default header values.

Invalid database name
*/
type PostBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *PostBadRequest) Error() string {
	return fmt.Sprintf("[POST /{db}][%d] postBadRequest  %+v", 400, o.Payload)
}
func (o *PostBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PostBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostUnauthorized creates a PostUnauthorized with default headers values
func NewPostUnauthorized() *PostUnauthorized {
	return &PostUnauthorized{}
}

/* PostUnauthorized describes a response with status code 401, with default header values.

Write privileges required
*/
type PostUnauthorized struct {
	Payload *models.ErrorResponse
}

func (o *PostUnauthorized) Error() string {
	return fmt.Sprintf("[POST /{db}][%d] postUnauthorized  %+v", 401, o.Payload)
}
func (o *PostUnauthorized) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PostUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostNotFound creates a PostNotFound with default headers values
func NewPostNotFound() *PostNotFound {
	return &PostNotFound{}
}

/* PostNotFound describes a response with status code 404, with default header values.

Database doesn’t exist
*/
type PostNotFound struct {
	Payload *models.ErrorResponse
}

func (o *PostNotFound) Error() string {
	return fmt.Sprintf("[POST /{db}][%d] postNotFound  %+v", 404, o.Payload)
}
func (o *PostNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PostNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostConflict creates a PostConflict with default headers values
func NewPostConflict() *PostConflict {
	return &PostConflict{}
}

/* PostConflict describes a response with status code 409, with default header values.

A Conflicting Document with same ID already exists
*/
type PostConflict struct {
	Payload *models.ErrorResponse
}

func (o *PostConflict) Error() string {
	return fmt.Sprintf("[POST /{db}][%d] postConflict  %+v", 409, o.Payload)
}
func (o *PostConflict) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PostConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
