// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DatabaseSizes database sizes
//
// swagger:model Database_sizes
type DatabaseSizes struct {

	// The size of live data inside the database, in bytes.
	Active int64 `json:"active,omitempty"`

	// The uncompressed size of database contents in bytes. sizes.file (number) – The size of the database file on disk in bytes. Views indexes are not included in the calculation.
	External int64 `json:"external,omitempty"`

	// An opaque string that describes the state of the database. Do not rely on this string for counting the number of updates.
	File int64 `json:"file,omitempty"`
}

// Validate validates this database sizes
func (m *DatabaseSizes) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this database sizes based on context it is used
func (m *DatabaseSizes) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DatabaseSizes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DatabaseSizes) UnmarshalBinary(b []byte) error {
	var res DatabaseSizes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
