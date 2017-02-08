package bootenvs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteBootenvHandlerFunc turns a function with the right signature into a delete bootenv handler
type DeleteBootenvHandlerFunc func(DeleteBootenvParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteBootenvHandlerFunc) Handle(params DeleteBootenvParams) middleware.Responder {
	return fn(params)
}

// DeleteBootenvHandler interface for that can handle valid delete bootenv params
type DeleteBootenvHandler interface {
	Handle(DeleteBootenvParams) middleware.Responder
}

// NewDeleteBootenv creates a new http.Handler for the delete bootenv operation
func NewDeleteBootenv(ctx *middleware.Context, handler DeleteBootenvHandler) *DeleteBootenv {
	return &DeleteBootenv{Context: ctx, Handler: handler}
}

/*DeleteBootenv swagger:route DELETE /bootenvs/{name} Bootenvs deleteBootenv

Delete Bootenv

*/
type DeleteBootenv struct {
	Context *middleware.Context
	Handler DeleteBootenvHandler
}

func (o *DeleteBootenv) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewDeleteBootenvParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}