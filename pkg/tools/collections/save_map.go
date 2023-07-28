package collections

import "sync"

type SafeMapCollection[T any] struct {
	collection map[string]T
	m          sync.Mutex
}

func NewSafeMapCollection[T any]() SafeMapCollection[T] {
	return SafeMapCollection[T]{
		collection: make(map[string]T),
	}
}

func (c *SafeMapCollection[T]) lock() {
	c.m.Lock()
}

func (c *SafeMapCollection[T]) unlock() {
	c.m.Unlock()
}

func (c *SafeMapCollection[T]) init() {
	if c.collection == nil {
		c.collection = make(map[string]T)
	}
}

func (c *SafeMapCollection[T]) Set(name string, item T) {
	c.lock()
	defer c.unlock()

	c.init()
	c.collection[name] = item
}

func (c *SafeMapCollection[T]) Get(name string) (item T, ok bool) {
	c.lock()
	defer c.unlock()

	c.init()
	item, ok = c.collection[name]

	return item, ok
}

func (c *SafeMapCollection[T]) MustGet(name string) (item T) {
	c.lock()
	defer c.unlock()

	c.init()
	item, ok := c.collection[name]

	if !ok {
		panic("Can't get collection item: " + name)
	}

	return item
}

func (c *SafeMapCollection[T]) Remove(name string) {
	c.lock()
	defer c.unlock()

	c.init()

	delete(c.collection, name)
}

func (c *SafeMapCollection[T]) Iter(cb func(item T) error) (err error) {
	c.lock()
	defer c.unlock()

	c.init()

	for _, item := range c.collection {
		err = cb(item)
		if err != nil {
			return err
		}
	}

	return nil
}
