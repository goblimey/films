package request

import (
	"net/url"
)

type Request interface {
	// URL returns the URL of the request
	URL() *url.URL

	// Method return the http request method ("GET", "POST" etc)
	Method() string

	// FormValue returns the given form value from the http request
	FormValue(name string) string

	// ParseForm parses the form data in the request
	ParseForm() error
}
