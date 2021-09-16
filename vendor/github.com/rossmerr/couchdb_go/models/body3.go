// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Body3 List of documents objects
//
// swagger:model body_3
type Body3 struct {

	// docs
	Docs []Document `json:"docs"`
}

// Validate validates this body 3
func (m *Body3) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this body 3 based on context it is used
func (m *Body3) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Body3) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Body3) UnmarshalBinary(b []byte) error {
	var res Body3
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
