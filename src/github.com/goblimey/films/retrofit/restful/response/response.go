package response

import (
	restful "github.com/emicklei/go-restful"
)

type Response interface {
	// Request gets the embedded restful response.
	Response() *restful.Response
}
