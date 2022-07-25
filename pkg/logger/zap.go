package logger

import (
	"go.uber.org/zap"
)

func NewZapLogger(isDev bool) (logger *zap.Logger, err error) {
	if isDev {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
