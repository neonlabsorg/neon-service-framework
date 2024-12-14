package collections

type BasicMapCollection[T any] map[string]T

func (c *BasicMapCollection[T]) init() {
	if c == nil {
		*c = make(map[string]T)
	}
}

func (c BasicMapCollection[T]) Set(name string, srv T) {
	c.init()
	c[name] = srv
}

func (c BasicMapCollection[T]) Get(name string) (srv T, ok bool) {
	c.init()
	srv, ok = c[name]
	return srv, ok
}

func (c BasicMapCollection[T]) MustGet(name string) (srv T) {
	c.init()
	srv, ok := c[name]

	if !ok {
		panic("Can't get collection item: " + name)
	}

	return srv
}

func (c BasicMapCollection[T]) Remove(name string) {
	c.init()
	delete(c, name)
}
