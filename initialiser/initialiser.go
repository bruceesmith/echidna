package initialiser

import (
	"github.com/bruceesmith/echidna/logger"
	"github.com/urfave/cli/v3"
)

type LogLevelFlag = cli.FlagBase[logger.LogLevel, cli.NoConfig, logLevelValue]

type logLevelValue struct {
	destination *logger.LogLevel
}

func (l logLevelValue) Create(val logger.LogLevel, p *logger.LogLevel, _ cli.NoConfig) cli.Value {
	*p = val
	return &logLevelValue{destination: p}
}

func (l logLevelValue) Get() any {
	return *l.destination
}

func (l logLevelValue) Set(s string) error {
	return l.destination.Set(s)
}

func (l logLevelValue) String() string {
	return l.destination.String()
}

func (l logLevelValue) ToString(lev logger.LogLevel) string {
	return lev.String()
}
