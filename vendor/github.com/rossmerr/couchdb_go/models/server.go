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

// Server server
//
// swagger:model Server
type Server struct {

	// couchdb
	Couchdb string `json:"couchdb,omitempty"`

	// uuid
	UUID string `json:"uuid,omitempty"`

	// vendor
	Vendor *ServerVendor `json:"vendor,omitempty"`

	// version
	Version string `json:"version,omitempty"`
}

// Validate validates this server
func (m *Server) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateVendor(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Server) validateVendor(formats strfmt.Registry) error {
	if swag.IsZero(m.Vendor) { // not required
		return nil
	}

	if m.Vendor != nil {
		if err := m.Vendor.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("vendor")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this server based on the context it is used
func (m *Server) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateVendor(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Server) contextValidateVendor(ctx context.Context, formats strfmt.Registry) error {

	if m.Vendor != nil {
		if err := m.Vendor.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("vendor")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Server) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Server) UnmarshalBinary(b []byte) error {
	var res Server
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
