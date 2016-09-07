package request

import (
	restful "github.com/emicklei/go-restful"
	httprequest "github.com/goblimey/films/retrofit/http/request"
)

type Request interface {
	// Request gets the restful request.
	Request() *restful.Request

	// HTTPRequest gets the HTTP request embedded in the restful request.
	HTTPRequest() httprequest.Request

	// PathParameter gets the named path parameter from the restful request.  It
	// returns an empty string if there is no such value.
	PathParameter(name string) string
}
