package service

import (
	"fmt"

	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/tools/collections"
)

type LoggerManager struct {
	loggers collections.BasicMapCollection[logger.Logger]
}

func NewLoggerManager(log logger.Logger) *LoggerManager {
	loggerManager := &LoggerManager{
		loggers: make(collections.BasicMapCollection[logger.Logger]),
	}
	loggerManager.setDefaultLogger(log)

	return loggerManager
}

func (l *LoggerManager) setDefaultLogger(logger logger.Logger) {
	l.loggers.Set("default", logger)
}

func (l *LoggerManager) SetLogger(name string, logger logger.Logger) {
	l.loggers.Set(name, logger)
}

func (l *LoggerManager) GetLogger() logger.Logger {
	log, ok := l.loggers.Get("default")

	if !ok {
		panic("default logger not found")
	}

	return log
}

func (l *LoggerManager) GetLoggerByName(name string) logger.Logger {
	log, ok := l.loggers.Get(name)

	if !ok {
		panic(fmt.Sprintf("logger with name %s not found", name))
	}

	return log
}
