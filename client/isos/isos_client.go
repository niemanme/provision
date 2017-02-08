package isos

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new isos API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for isos API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
DeleteIso deletes iso
*/
func (a *Client) DeleteIso(params *DeleteIsoParams) (*DeleteIsoNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteIsoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "delete-iso",
		Method:             "DELETE",
		PathPattern:        "/isos/{name}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DeleteIsoReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeleteIsoNoContent), nil

}

/*
GetIso gets iso
*/
func (a *Client) GetIso(params *GetIsoParams, writer io.Writer) (*GetIsoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetIsoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "get-iso",
		Method:             "GET",
		PathPattern:        "/isos/{name}",
		ProducesMediaTypes: []string{"application/octet-stream"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetIsoReader{formats: a.formats, writer: writer},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetIsoOK), nil

}

/*
ListIsos lists isos
*/
func (a *Client) ListIsos(params *ListIsosParams) (*ListIsosOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListIsosParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "list-isos",
		Method:             "GET",
		PathPattern:        "/isos",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListIsosReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ListIsosOK), nil

}

/*
PostIso creates iso
*/
func (a *Client) PostIso(params *PostIsoParams) (*PostIsoCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostIsoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "post-iso",
		Method:             "POST",
		PathPattern:        "/isos/{name}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/octet-stream"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostIsoReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PostIsoCreated), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}