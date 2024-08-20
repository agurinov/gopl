package crontab

import (
	gocron "github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type gocronLogger struct {
	*zap.SugaredLogger
}

func (l gocronLogger) Debug(msg string, args ...any) { l.Debugf(msg, args...) }
func (l gocronLogger) Error(msg string, args ...any) { l.Errorf(msg, args...) }
func (l gocronLogger) Info(msg string, args ...any)  { l.Infof(msg, args...) }
func (l gocronLogger) Warn(msg string, args ...any)  { l.Warnf(msg, args...) }

func loggerAdapter(logger *zap.Logger) gocron.Logger {
	return gocronLogger{
		logger.Sugar(),
	}
}
