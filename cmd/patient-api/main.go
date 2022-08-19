package main

import (
	"github.com/getsentry/sentry-go"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	_ "github.com/synthia-telemed/backend-api/docs"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/config"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/logger"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/sms"
	"github.com/synthia-telemed/backend-api/pkg/token"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

// @title           Synthia Patient Backend API
// @version         1.0.0
// @description     This is a Synthia patient backend API.
// @accept json
// @produce json
// @BasePath  /patient/api

// @securityDefinitions.apikey  UserID
// @in                          header
// @name                        X-USER-ID
// @description					UserID that interacts with the API. Normally this header is set by Heimdall. Development Only!
// @securityDefinitions.apikey  JWSToken
// @in                          header
// @name                        Authorization
// @description					JWS that user possess
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
	measurementDataStore, err := datastore.NewGormMeasurementDataStore(db)
	if err != nil {
		sugaredLogger.Fatalw("Failed to create measurement data store", "error", err)
	}

	hospitalSysClient := hospital.NewGraphQLClient(&cfg.HospitalClient)
	smsClient := sms.NewTwilioClient(&cfg.SMS)
	cacheClient := cache.NewRedisClient(&cfg.Cache)
	tokenService, err := token.NewGRPCTokenService(&cfg.Token)
	if err != nil {
		sentry.CaptureException(err)
		sugaredLogger.Fatalw("Failed to create token service", "error", err)
	}
	paymentClient, err := payment.NewOmisePaymentClient(&cfg.Payment)
	if err != nil {
		sentry.CaptureException(err)
		sugaredLogger.Fatalw("Failed to create payment client", "error", err)
	}

	// Handler
	authHandler := handler.NewAuthHandler(patientDataStore, hospitalSysClient, smsClient, cacheClient, tokenService, sugaredLogger)
	paymentHandler := handler.NewPaymentHandler(paymentClient, patientDataStore, sugaredLogger)
	measurementHandler := handler.NewMeasurementHandler(measurementDataStore, sugaredLogger)

	ginServer := server.NewGinServer(cfg, sugaredLogger)
	apiGroup := ginServer.Group("/api")
	authHandler.Register(apiGroup)
	paymentHandler.Register(apiGroup)
	measurementHandler.Register(apiGroup)
	ginServer.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	ginServer.ListenAndServe()
}
