package response

import (
	restful "github.com/emicklei/go-restful"
)

// The ConcreteResponse struct implements the Response interface.
type ConcreteResponse struct {
	responseField *restful.Response
}

// MakeResponse creates and returns a new Response object containing the given
// restful response.
func MakeResponse(resp *restful.Response) Response {
	concreteResponse := &ConcreteResponse{}
	concreteResponse.SetResponse(resp)
	var r Response = concreteResponse
	return r
}

func (cr ConcreteResponse) Response() *restful.Response {
	return cr.responseField
}

// SetResponse sets the restful response field.
func (cr *ConcreteResponse) SetResponse(resp *restful.Response) {
	cr.responseField = resp
}
