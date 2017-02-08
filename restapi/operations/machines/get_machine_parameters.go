package machines

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetMachineParams creates a new GetMachineParams object
// with the default values initialized.
func NewGetMachineParams() GetMachineParams {
	var ()
	return GetMachineParams{}
}

// GetMachineParams contains all the bound params for the get machine operation
// typically these are obtained from a http.Request
//
// swagger:parameters get-machine
type GetMachineParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*
	  Required: true
	  In: path
	*/
	UUID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetMachineParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rUUID, rhkUUID, _ := route.Params.GetOK("uuid")
	if err := o.bindUUID(rUUID, rhkUUID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetMachineParams) bindUUID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.UUID = raw

	return nil
}