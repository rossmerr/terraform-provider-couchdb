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

// UpReader is a Reader for the Up structure.
type UpReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewUpNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewUpOK creates a UpOK with default headers values
func NewUpOK() *UpOK {
	return &UpOK{}
}

/* UpOK describes a response with status code 200, with default header values.

Request completed successfully
*/
type UpOK struct {
	Payload *models.InlineResponse2004
}

func (o *UpOK) Error() string {
	return fmt.Sprintf("[GET /_up][%d] upOK  %+v", 200, o.Payload)
}
func (o *UpOK) GetPayload() *models.InlineResponse2004 {
	return o.Payload
}

func (o *UpOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.InlineResponse2004)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpNotFound creates a UpNotFound with default headers values
func NewUpNotFound() *UpNotFound {
	return &UpNotFound{}
}

/* UpNotFound describes a response with status code 404, with default header values.

The server is unavailable for requests at this time.
*/
type UpNotFound struct {
	Payload *models.ErrorResponse
}

func (o *UpNotFound) Error() string {
	return fmt.Sprintf("[GET /_up][%d] upNotFound  %+v", 404, o.Payload)
}
func (o *UpNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
