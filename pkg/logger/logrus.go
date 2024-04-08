package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type logger struct {
	pimpl *logrus.Logger
}

var currentLogLevel = logrus.DebugLevel

func SetLogLevel(level string) error {
	n, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("failed parse: %s err:%v", level, err)
	}
	currentLogLevel = n
	logrus.SetLevel(currentLogLevel)
	return nil
}

func NewLogger() Logger {
	pimpl := logrus.New()
	pimpl.SetLevel(currentLogLevel)
	return &logger{
		pimpl: pimpl,
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.pimpl.Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.pimpl.Infof(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.pimpl.Warnf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.pimpl.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.pimpl.Fatalf(format, args...)
}
