package request

import (
	"net/http"
	"net/url"
)

type ConcreteRequest struct {
	requestField *http.Request
}

// MakeRequest creates and returns a new ConcreteRequest object containing the given
// http request.
func MakeRequest(request *http.Request) Request {
	var concreteRequest ConcreteRequest
	concreteRequest.requestField = request
	var restRequest Request = &concreteRequest
	return restRequest
}

// URI return the URI of the request
func (cr ConcreteRequest) URL() *url.URL {
	return cr.requestField.URL
}

// Method return the http request method ("GET", "POST" etc)
func (cr ConcreteRequest) Method() string {
	return cr.requestField.Method
}

// FormValue returns the given form value from the http request
func (cr ConcreteRequest) FormValue(name string) string {
	return cr.requestField.FormValue(name)
}

// ParseForm parses the form data in the request
func (cr ConcreteRequest) ParseForm() error {
	return cr.requestField.ParseForm()
}
