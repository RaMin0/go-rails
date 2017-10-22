package rails

import (
	"fmt"
	"net/http"
	"regexp"
)

const (
	RouteOptionFormat = "format"
)

type Router struct {
	routes map[string][]*route
}

type route struct {
	methods []string
	action  *action
	format  string
}

type routeOption struct {
	name  string
	value string
}

func NewRouter() *Router {
	return &Router{
		routes: map[string][]*route{},
	}
}

func (r *Router) Resources(name string, c *controller, opts ...*routeOption) *Router {
	r.match(fmt.Sprintf("^/%s$", name), c.Action(ActionIndex), opts...).Via(http.MethodGet)
	r.match(fmt.Sprintf("^/%s/(?P<id>\\d+)$", name), c.Action(ActionShow), opts...).Via(http.MethodGet)
	return r
}

func (r *Router) match(pattern string, action *action, opts ...*routeOption) *route {
	t := route{action: action, format: FormatHTML}

	for _, opt := range opts {
		switch opt.name {
		case RouteOptionFormat:
			t.format = opt.value
		}
	}

	if _, ok := r.routes[pattern]; !ok {
		r.routes[pattern] = []*route{}
	}
	r.routes[pattern] = append(r.routes[pattern], &t)

	return &t
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for pattern, routes := range router.routes {
		for _, t := range routes {
			patternRegexp := regexp.MustCompile(pattern)
			if !patternRegexp.MatchString(r.URL.Path) {
				continue
			}
			patternMatch := patternRegexp.FindStringSubmatch(r.URL.Path)

			middlewares.wrap(func(w http.ResponseWriter, r *http.Request) {
				for _, m := range t.methods {
					if r.Method == m {
						ctx := newContext()

						for i, n := range patternRegexp.SubexpNames() {
							ctx.params[n] = patternMatch[i]
						}

						t.action.fn(ctx)

						v := view{t.action.templateName(), t.format, ctx.variables}
						w.Header().Add("Content-Type", v.contentType())
						v.render(w)
						return
					}
				}
			})(w, r)

			return
		}
	}

	http.NotFound(w, r)
}

func (t *route) Via(methods ...string) *route {
	t.methods = append(t.methods, methods...)
	return t
}

func RouteOption(name, value string) *routeOption {
	return &routeOption{name, value}
}
