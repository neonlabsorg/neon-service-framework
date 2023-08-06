package alerts

import (
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
)

type Dispatcher struct {
	adapter         Adapter
	reservedAdapter Adapter
	registry        *Registry
	log             logger.Logger
}

func NewAlertDispatcher(
	adapter Adapter,
	log logger.Logger,
) *Dispatcher {
	return &Dispatcher{
		adapter: adapter,
		log:     log,
	}
}

func (d *Dispatcher) UseRegistry() {
	if d.registry == nil {
		d.registry = NewRegistry()
	}
}

func (d *Dispatcher) UseReservedAdapter(adapter Adapter) {
	d.reservedAdapter = adapter
}

func (d *Dispatcher) Dispatch(alert Alert) {
	if d.registry != nil {
		d.dispatchWithRegistry(alert)
		return
	}

	go d.send(alert)
}

func (d *Dispatcher) dispatchWithRegistry(alert Alert) {
	if !d.registry.Exists(alert.GetName()) {
		d.log.Error().Msgf("sending not registered alert %s", alert.GetName().String())
		return
	}

	go d.send(alert)
}

func (d *Dispatcher) send(alert Alert) {
	err := d.sendByAdapter(alert)
	if err != nil {
		d.log.Error().Err(err).Msgf("error on send alert %s by main %s adapter", alert.GetName().String(), d.adapter.GetName())
	} else {
		return
	}

	if d.reservedAdapter == nil {
		return
	}

	err = d.sendByReservedAdapter(alert)
	if err != nil {
		d.log.Error().Err(err).Msgf("error on send alert %s by reserved %s adapter", alert.GetName().String(), d.reservedAdapter.GetName())
	}
}

func (d *Dispatcher) sendByAdapter(alert Alert) error {
	if d.adapter == nil {
		return ErrAdapterWasntInstalled
	}
	return d.adapter.Send(alert)
}

func (d *Dispatcher) sendByReservedAdapter(alert Alert) error {
	if d.reservedAdapter == nil {
		return ErrReservedAdapterWasntInstalled
	}
	return d.reservedAdapter.Send(alert)
}
