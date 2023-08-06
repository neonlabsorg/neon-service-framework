package alerts

import (
	"fmt"

	"github.com/neonlabsorg/neon-service-framework/pkg/tools/collections"
)

type Registry struct {
	list collections.SafeMapCollection[Alert]
}

func NewRegistry() *Registry {
	return &Registry{
		list: collections.NewSafeMapCollection[Alert](),
	}
}

func (r *Registry) Register(alert Alert) {
	_, ok := r.list.Get(alert.GetName().String())
	if ok {
		panic(fmt.Sprintf("alert registry: trying to register alert with same name %s", alert.GetName()))
	}

	r.list.Set(alert.GetName().String(), alert)
}

func (r *Registry) Exists(name Name) bool {
	return r.list.Exists(name.String())
}

func (r *Registry) GetAll() (list []Alert) {
	err := r.list.Iter(func(alert Alert) error {
		list = append(list, alert)
		return nil
	})
	if err != nil {
		panic("error on iter alerts collection")
	}

	return list
}
