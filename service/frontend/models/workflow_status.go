// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// WorkflowStatus workflow status
//
// swagger:model workflowStatus
type WorkflowStatus struct {

	// finished at
	FinishedAt string `json:"FinishedAt,omitempty"`

	// log
	Log string `json:"Log,omitempty"`

	// name
	Name string `json:"Name,omitempty"`

	// nodes
	Nodes []string `json:"Nodes"`

	// on cancel
	OnCancel string `json:"OnCancel,omitempty"`

	// on exit
	OnExit string `json:"OnExit,omitempty"`

	// on failure
	OnFailure string `json:"OnFailure,omitempty"`

	// on success
	OnSuccess string `json:"OnSuccess,omitempty"`

	// params
	Params string `json:"Params,omitempty"`

	// pid
	Pid int64 `json:"Pid,omitempty"`

	// request Id
	RequestID string `json:"RequestId,omitempty"`

	// started at
	StartedAt string `json:"StartedAt,omitempty"`

	// status
	Status int64 `json:"Status,omitempty"`

	// status text
	StatusText string `json:"StatusText,omitempty"`
}

// Validate validates this workflow status
func (m *WorkflowStatus) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this workflow status based on context it is used
func (m *WorkflowStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *WorkflowStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WorkflowStatus) UnmarshalBinary(b []byte) error {
	var res WorkflowStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
