package template

import (
	"io"
)

// The Template interface mimics some of the the html/Template functionality,
// allowing templates to be mocked.
type Template interface {
	// Execute executes the template
	Execute(wr io.Writer, data interface{}) error
}
