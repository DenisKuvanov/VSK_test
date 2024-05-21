package log

import (
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
	Log       = logger.Sugar()
)
