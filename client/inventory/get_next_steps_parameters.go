// Code generated by go-swagger; DO NOT EDIT.

package inventory

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetNextStepsParams creates a new GetNextStepsParams object
// with the default values initialized.
func NewGetNextStepsParams() *GetNextStepsParams {
	var ()
	return &GetNextStepsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetNextStepsParamsWithTimeout creates a new GetNextStepsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetNextStepsParamsWithTimeout(timeout time.Duration) *GetNextStepsParams {
	var ()
	return &GetNextStepsParams{

		timeout: timeout,
	}
}

// NewGetNextStepsParamsWithContext creates a new GetNextStepsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetNextStepsParamsWithContext(ctx context.Context) *GetNextStepsParams {
	var ()
	return &GetNextStepsParams{

		Context: ctx,
	}
}

// NewGetNextStepsParamsWithHTTPClient creates a new GetNextStepsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetNextStepsParamsWithHTTPClient(client *http.Client) *GetNextStepsParams {
	var ()
	return &GetNextStepsParams{
		HTTPClient: client,
	}
}

/*GetNextStepsParams contains all the parameters to send to the API endpoint
for the get next steps operation typically these are written to a http.Request
*/
type GetNextStepsParams struct {

	/*NodeID
	  ID of node

	*/
	NodeID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get next steps params
func (o *GetNextStepsParams) WithTimeout(timeout time.Duration) *GetNextStepsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get next steps params
func (o *GetNextStepsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get next steps params
func (o *GetNextStepsParams) WithContext(ctx context.Context) *GetNextStepsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get next steps params
func (o *GetNextStepsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get next steps params
func (o *GetNextStepsParams) WithHTTPClient(client *http.Client) *GetNextStepsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get next steps params
func (o *GetNextStepsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNodeID adds the nodeID to the get next steps params
func (o *GetNextStepsParams) WithNodeID(nodeID string) *GetNextStepsParams {
	o.SetNodeID(nodeID)
	return o
}

// SetNodeID adds the nodeId to the get next steps params
func (o *GetNextStepsParams) SetNodeID(nodeID string) {
	o.NodeID = nodeID
}

// WriteToRequest writes these params to a swagger request
func (o *GetNextStepsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param node_id
	if err := r.SetPathParam("node_id", o.NodeID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}