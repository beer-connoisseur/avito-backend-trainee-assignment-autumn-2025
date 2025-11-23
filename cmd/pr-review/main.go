package main

import (
	"log"
	"os"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/config"
	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg, err := config.New()

	if err != nil {
		log.Fatalf("can not get application config: %s", err)
	}

	logger, err := NewFileLogger()

	if err != nil {
		log.Fatalf("can not initialize logger: %s", err)
	}

	app.Run(logger, cfg)
}

const (
	FilePermissionsExec = 0755
	FilePermissionsWR   = 0644
)

func NewFileLogger() (*zap.Logger, error) {
	const logFile = "/app/logs/pr_review.log"
	_ = os.MkdirAll("/app/logs", FilePermissionsExec)

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, FilePermissionsWR)

	if err != nil {
		return nil, err
	}

	writeSyncer := zapcore.AddSync(file)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)

	return zap.New(core), nil
}
