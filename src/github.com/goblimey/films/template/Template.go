package template

import (
	htmlTemplate "html/template"
	"io"
)

// The Template interface mimics some of the the html/Template functionality,
// allowing templates to be mocked.
type Template interface {
	// HtmlTemplate gets the underlying html template
	HtmlTemplate() *htmlTemplate.Template
	// Execute executes the template
	Execute(wr io.Writer, data interface{}) error
}
