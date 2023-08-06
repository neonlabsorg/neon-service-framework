package alerts

type Context struct {
	project  string
	service  string
	instance string
}

func NewContext(project string, service string, instance string) *Context {
	return &Context{
		project:  project,
		service:  service,
		instance: instance,
	}
}

func (c *Context) GetProject() string {
	return c.project
}

func (c *Context) GetService() string {
	return c.service
}

func (c *Context) GetInstance() string {
	return c.instance
}
