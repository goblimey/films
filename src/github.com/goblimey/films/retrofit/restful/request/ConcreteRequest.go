package request

import (
	httprequest "github.com/goblimey/films/retrofit/http/request"

	restful "github.com/emicklei/go-restful"
)

// The ConcreteRequest struct implements the Request interface and contains a
// restful http request and an httprequest constructed from it.
type ConcreteRequest struct {
	requestField *restful.Request
	httpRequest  httprequest.Request
}

// MakeRequest creates and returns a new Request object containing the given
// restful request.
func MakeRequest(req *restful.Request) Request {
	concreteRequest := &ConcreteRequest{}
	concreteRequest.SetRequest(req)
	var r Request = concreteRequest
	return r
}

// Request return the restful request
func (cr ConcreteRequest) Request() *restful.Request {
	return cr.requestField
}

// HTTPRequest return the embedded http request
func (cr ConcreteRequest) HTTPRequest() httprequest.Request {
	return cr.httpRequest
}

// PathParameter gets the named path parameter from the restful request.  It
// returns an empty string if there is no such value.
func (cr ConcreteRequest) PathParameter(name string) string {
	return cr.requestField.PathParameter(name)
}

// SetRequest sets the restful request.
func (cr *ConcreteRequest) SetRequest(req *restful.Request) {
	cr.requestField = req
	cr.httpRequest = httprequest.MakeRequest(req.Request)
}
