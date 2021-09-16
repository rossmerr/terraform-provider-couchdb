// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// InlineResponse2004 inline response 200 4
//
// swagger:model inline_response_200_4
type InlineResponse2004 struct {

	// status
	Status string `json:"status,omitempty"`
}

// Validate validates this inline response 200 4
func (m *InlineResponse2004) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this inline response 200 4 based on context it is used
func (m *InlineResponse2004) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *InlineResponse2004) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InlineResponse2004) UnmarshalBinary(b []byte) error {
	var res InlineResponse2004
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}