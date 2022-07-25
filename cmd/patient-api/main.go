package main

import (
	"github.com/synthia-telemed/backend-api/pkg/config"
	"github.com/synthia-telemed/backend-api/pkg/logger"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to parse ENV:", err)
	}

	zapLogger, err := logger.NewZapLogger(cfg.Mode == "development")
	if err != nil {
		log.Fatalln("Failed to initialized Zap:", err)
	}
	defer zapLogger.Sync()
	sugaredLogger := zapLogger.Sugar()

	ginServer := server.NewGinServer(cfg, sugaredLogger)
	ginServer.ListenAndServe()
}
