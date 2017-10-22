package rails

type Context struct {
	variables map[string]interface{}
	params    map[string]string
}

func newContext() *Context {
	return &Context{
		variables: map[string]interface{}{},
		params:    map[string]string{},
	}
}

func (ctx *Context) Var(name string, value ...interface{}) interface{} {
	if len(value) > 0 {
		ctx.variables[name] = value[0]
	}
	return ctx.variables[name]
}

func (ctx *Context) Param(name string) string {
	return ctx.params[name]
}
