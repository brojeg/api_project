package shared

import "go.uber.org/zap"

func Init() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger.Sugar()
}
