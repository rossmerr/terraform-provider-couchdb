// Code generated by go-swagger; DO NOT EDIT.

package partition

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/rossmerr/couchdb_go/models"
)

// PartitionDesignDocSearchReader is a Reader for the PartitionDesignDocSearch structure.
type PartitionDesignDocSearchReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PartitionDesignDocSearchReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPartitionDesignDocSearchOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPartitionDesignDocSearchBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewPartitionDesignDocSearchUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewPartitionDesignDocSearchNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPartitionDesignDocSearchOK creates a PartitionDesignDocSearchOK with default headers values
func NewPartitionDesignDocSearchOK() *PartitionDesignDocSearchOK {
	return &PartitionDesignDocSearchOK{}
}

/* PartitionDesignDocSearchOK describes a response with status code 200, with default header values.

Request completed successfully
*/
type PartitionDesignDocSearchOK struct {

	/* Response signature
	 */
	ETag string

	/* chunked
	 */
	TransferEncoding string

	Payload *models.Pagination
}

func (o *PartitionDesignDocSearchOK) Error() string {
	return fmt.Sprintf("[GET /{db}/_partition/{partition}/_design/{ddoc}/_search/{index}][%d] partitionDesignDocSearchOK  %+v", 200, o.Payload)
}
func (o *PartitionDesignDocSearchOK) GetPayload() *models.Pagination {
	return o.Payload
}

func (o *PartitionDesignDocSearchOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

	o.Payload = new(models.Pagination)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPartitionDesignDocSearchBadRequest creates a PartitionDesignDocSearchBadRequest with default headers values
func NewPartitionDesignDocSearchBadRequest() *PartitionDesignDocSearchBadRequest {
	return &PartitionDesignDocSearchBadRequest{}
}

/* PartitionDesignDocSearchBadRequest describes a response with status code 400, with default header values.

Invalid request
*/
type PartitionDesignDocSearchBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *PartitionDesignDocSearchBadRequest) Error() string {
	return fmt.Sprintf("[GET /{db}/_partition/{partition}/_design/{ddoc}/_search/{index}][%d] partitionDesignDocSearchBadRequest  %+v", 400, o.Payload)
}
func (o *PartitionDesignDocSearchBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PartitionDesignDocSearchBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPartitionDesignDocSearchUnauthorized creates a PartitionDesignDocSearchUnauthorized with default headers values
func NewPartitionDesignDocSearchUnauthorized() *PartitionDesignDocSearchUnauthorized {
	return &PartitionDesignDocSearchUnauthorized{}
}

/* PartitionDesignDocSearchUnauthorized describes a response with status code 401, with default header values.

Read permission required
*/
type PartitionDesignDocSearchUnauthorized struct {
	Payload *models.ErrorResponse
}

func (o *PartitionDesignDocSearchUnauthorized) Error() string {
	return fmt.Sprintf("[GET /{db}/_partition/{partition}/_design/{ddoc}/_search/{index}][%d] partitionDesignDocSearchUnauthorized  %+v", 401, o.Payload)
}
func (o *PartitionDesignDocSearchUnauthorized) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PartitionDesignDocSearchUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPartitionDesignDocSearchNotFound creates a PartitionDesignDocSearchNotFound with default headers values
func NewPartitionDesignDocSearchNotFound() *PartitionDesignDocSearchNotFound {
	return &PartitionDesignDocSearchNotFound{}
}

/* PartitionDesignDocSearchNotFound describes a response with status code 404, with default header values.

Specified database, design document or view is missed
*/
type PartitionDesignDocSearchNotFound struct {
	Payload *models.ErrorResponse
}

func (o *PartitionDesignDocSearchNotFound) Error() string {
	return fmt.Sprintf("[GET /{db}/_partition/{partition}/_design/{ddoc}/_search/{index}][%d] partitionDesignDocSearchNotFound  %+v", 404, o.Payload)
}
func (o *PartitionDesignDocSearchNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PartitionDesignDocSearchNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
