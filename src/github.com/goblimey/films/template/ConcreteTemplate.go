package template

import (
	htmlTemplate "html/template"
	"io"
	"log"
)

// ConcreteTemplate satisfies the Template interface and is a proxy for
// html/template.
type ConcreteTemplate struct {
	tplt *htmlTemplate.Template
}

func (ct ConcreteTemplate) HtmlTemplate() *htmlTemplate.Template {
	return ct.tplt
}

// Define the functions.

func (tp *ConcreteTemplate) SetHTMLTemplate(tplt *htmlTemplate.Template) {
	tp.tplt = tplt
}

// Must is a proxy for htmlTemplate.Must.
func Must(tp Template, err error) Template {
	log.SetPrefix("ConcreteTemplate.Must ")
	tplt := htmlTemplate.Must(tp.HtmlTemplate(), err)
	var result Template = &ConcreteTemplate{tplt}
	return result
}

// ParseFiles is a proxy for htmlTemplate.ParseFiles.
func ParseFiles(filenames ...string) (*Template, error) {
	log.SetPrefix("ConcreteTemplate.ParseFiles ")
	tp, err := htmlTemplate.ParseFiles(filenames...)
	if err != nil {
		return nil, err
	}
	var result Template = &ConcreteTemplate{tp}
	return &result, nil
}

// Define the methods.

//Execute is a proxy for htmlTemplate.Execute.
func (ct ConcreteTemplate) Execute(wr io.Writer, data interface{}) error {
	log.SetPrefix("ConcreteTemplate.Execute ")
	return ct.HtmlTemplate().Execute(wr, data)
}
