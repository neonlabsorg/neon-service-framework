package errors

type ErrorContext map[string]string

func (c *ErrorContext) init() {
	c = new(ErrorContext)
}

func (c ErrorContext) Set(key string, value string) {
	if c == nil {
		c.init()
	}
	c[key] = value
}

func (c ErrorContext) Get(key string) (value string) {
	if c == nil {
		c.init()
	}
	if value, ok := c[key]; ok {
		return value
	}

	return ""
}

func (c ErrorContext) Len() int {
	return len(c)
}

func (c ErrorContext) Clear() {
	c.init()
}
