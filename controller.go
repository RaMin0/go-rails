package rails

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
)

const (
	ActionIndex = "index"
	ActionShow  = "show"
)

type controller struct {
	actions map[string]*action
}
type action struct {
	controller *controller
	name       string
	fn         actionFunc
}
type actionFunc func(*Context)

func NewController() *controller {
	return &controller{
		actions: map[string]*action{},
	}
}

func (c *controller) AddAction(name string, fn actionFunc) *controller {
	c.actions[name] = &action{c, name, fn}
	return c
}

func (c *controller) Action(name string) *action {
	return c.actions[name]
}

func (a *action) templateName() string {
	ptr := reflect.ValueOf(a.fn).Pointer()
	file, _ := runtime.FuncForPC(ptr).FileLine(ptr)

	controllerRegexp := regexp.MustCompile("app/controllers/(.+)\\.go$")
	controllerName := controllerRegexp.FindStringSubmatch(file)[1]

	return fmt.Sprintf("%s/%s", controllerName, a.name)
}
