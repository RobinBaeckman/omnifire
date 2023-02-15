package logger

import (
	"github.com/sirupsen/logrus"
)

type Log struct {
	*logrus.Logger
}

func New() *Log {
	var log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return &Log{log}
}
