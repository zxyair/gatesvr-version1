package log_test

import (
	"gatesvr/log"
	"testing"
)

func TestLog(t *testing.T) {
	logger := log.NewLogger(log.WithFormat(log.JsonFormat))

	logger.Debug("welcome to due-framework")
	logger.Info("welcome to due-framework")
	logger.Warn("welcome to due-framework")
	logger.Error("welcome to due-framework")
}
