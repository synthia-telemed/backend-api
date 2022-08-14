package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/config"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/logger"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/sms"
	"github.com/synthia-telemed/backend-api/pkg/token"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{})
	if err != nil {
		sugaredLogger.Fatalw("Failed to connect to database", "error", err)
	}
	patientDataStore, err := datastore.NewGormPatientDataStore(db)
	if err != nil {
		sugaredLogger.Fatalw("Failed to create patient data store", "error", err)
	}

	hospitalSysClient := hospital.NewGraphQLClient(&cfg.HospitalClient)
	smsClient := sms.NewTwilioClient(&cfg.SMS)
	cacheClient := cache.NewRedisClient(&cfg.Cache)
	tokenService, err := token.NewGRPCTokenService(&cfg.Token)
	if err != nil {
		sentry.CaptureException(err)
		sugaredLogger.Fatalw("Failed to create token service", "error", err)
	}

	// Handler
	authHandler := handler.NewAuthHandler(patientDataStore, hospitalSysClient, smsClient, cacheClient, tokenService, sugaredLogger)

	ginServer := server.NewGinServer(cfg, sugaredLogger)
	authGroup := ginServer.Group("/api/auth")
	authGroup.POST("/signin", authHandler.Signin)
	authGroup.POST("/verify", authHandler.VerifyOTP)
	ginServer.ListenAndServe()
}
