// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// SearchIndex search index
//
// swagger:model SearchIndex
type SearchIndex struct {

	// committed seq
	CommittedSeq int64 `json:"committed_seq,omitempty"`

	// disk size
	DiskSize int64 `json:"disk_size,omitempty"`

	// doc count
	DocCount int64 `json:"doc_count,omitempty"`

	// doc del count
	DocDelCount int64 `json:"doc_del_count,omitempty"`

	// pending seq
	PendingSeq int64 `json:"pending_seq,omitempty"`
}

// Validate validates this search index
func (m *SearchIndex) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this search index based on context it is used
func (m *SearchIndex) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SearchIndex) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SearchIndex) UnmarshalBinary(b []byte) error {
	var res SearchIndex
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
