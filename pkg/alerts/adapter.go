package alerts

type Adapter interface {
	Send(alert Alert) error
	GetName() string
}
