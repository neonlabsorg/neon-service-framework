package alerts

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
)

type ConsoleAdapter struct {
	log logger.Logger
}

func NewConsoleAdapter(
	log logger.Logger,
) *ConsoleAdapter {

	return &ConsoleAdapter{
		log: log,
	}
}

func (s *ConsoleAdapter) GetName() string {
	return "console"
}

func (s *ConsoleAdapter) Send(alert Alert) error {
	message := spew.Sdump(alert)
	println(message)
	s.log.Info().Msg(message)
	return nil
}
