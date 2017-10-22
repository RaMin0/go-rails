package rails

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

var helpers = template.FuncMap{
	"linkTo": func(title string, o interface{}) template.HTML {
		return template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, urlFor(o), title))
	},
}

type paramable interface {
	Param() string
}

func urlFor(i interface{}) string {
	if o, ok := i.(paramable); ok {
		paramName := reflect.ValueOf(o).Elem().Type().Name()
		return fmt.Sprintf("/%vs/%s", strings.ToLower(paramName), o.Param())
	}

	panic(fmt.Sprintf("urlFor failed for: %#v", i))
}
