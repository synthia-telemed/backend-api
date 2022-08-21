package main

import (
	"github.com/getsentry/sentry-go"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/synthia-telemed/backend-api/pkg/config"
	"github.com/synthia-telemed/backend-api/pkg/logger"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"log"
	"time"
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

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		TracesSampleRate: 1.0,
	}); err != nil {
		sugaredLogger.Fatalw("Sentry initialization failed", "error", err)
	}
	defer sentry.Flush(2 * time.Second)

	ginServer := server.NewGinServer(cfg, sugaredLogger)
	ginServer.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	ginServer.ListenAndServe()
}
