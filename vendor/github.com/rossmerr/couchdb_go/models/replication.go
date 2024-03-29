// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Replication replication
//
// swagger:model Replication
type Replication struct {
	BasicDoc

	OK
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *Replication) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 BasicDoc
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.BasicDoc = aO0

	// AO1
	var aO1 OK
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.OK = aO1

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m Replication) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.BasicDoc)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	aO1, err := swag.WriteJSON(m.OK)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this replication
func (m *Replication) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with BasicDoc
	if err := m.BasicDoc.Validate(formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with OK
	if err := m.OK.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validate this replication based on the context it is used
func (m *Replication) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with BasicDoc
	if err := m.BasicDoc.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with OK
	if err := m.OK.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Replication) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Replication) UnmarshalBinary(b []byte) error {
	var res Replication
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
