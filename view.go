package rails

import (
	"fmt"
	"html/template"
	"io"
)

const (
	FormatHTML = "html"
	FormatJSON = "json"
)

type view struct {
	template string
	format   string
	data     interface{}
}

func (v *view) contentType() string {
	switch v.format {
	case FormatHTML:
		return "text/html"
	case FormatJSON:
		return "application/json"
	}
	return "text/plain"
}

func (v *view) render(w io.Writer) {
	if err := template.Must(template.Must(template.
		New("app/views").
		Funcs(helpers).
		ParseFiles(fmt.Sprintf("app/views/%s.%s.tmpl", v.template, v.format))).
		ParseGlob("app/views/layouts/*.tmpl")).
		ExecuteTemplate(w, "layout", v.data); err != nil {
		panic(err)
	}
}
