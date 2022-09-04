package main

import (
	"github.com/getsentry/sentry-go"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/synthia-telemed/backend-api/cmd/patient-api/docs"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/clock"
	"github.com/synthia-telemed/backend-api/pkg/config"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/logger"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/sms"
	"github.com/synthia-telemed/backend-api/pkg/token"
	"go.uber.org/zap"
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
	assertFatalError(sugaredLogger, err, "Failed to connect to database")

	patientDataStore, err := datastore.NewGormPatientDataStore(db)
	assertFatalError(sugaredLogger, err, "Failed to create patient data store")
	measurementDataStore, err := datastore.NewGormMeasurementDataStore(db)
	assertFatalError(sugaredLogger, err, "Failed to create measurement data store")
	creditCardDataStore, err := datastore.NewGormCreditCardDataStore(db)
	assertFatalError(sugaredLogger, err, "Failed to create credit card data store")
	paymentDataStore, err := datastore.NewGormPaymentDataStore(db)
	assertFatalError(sugaredLogger, err, "Failed to create payment data store")

	hospitalSysClient := hospital.NewGraphQLClient(&cfg.HospitalClient)
	smsClient := sms.NewTwilioClient(&cfg.SMS)
	cacheClient := cache.NewRedisClient(&cfg.Cache)
	tokenService, err := token.NewGRPCTokenService(&cfg.Token)
	assertFatalError(sugaredLogger, err, "Failed to create token service")
	paymentClient, err := payment.NewOmisePaymentClient(&cfg.Payment)
	assertFatalError(sugaredLogger, err, "Failed to create payment client")
	realClock := clock.NewRealClock()

	// Handler
	authHandler := handler.NewAuthHandler(patientDataStore, hospitalSysClient, smsClient, cacheClient, tokenService, realClock, sugaredLogger)
	paymentHandler := handler.NewPaymentHandler(paymentClient, patientDataStore, creditCardDataStore, hospitalSysClient, paymentDataStore, realClock, sugaredLogger)
	measurementHandler := handler.NewMeasurementHandler(measurementDataStore, sugaredLogger)
	appointmentHandler := handler.NewAppointmentHandler(patientDataStore, paymentDataStore, hospitalSysClient, realClock, sugaredLogger)

	ginServer := server.NewGinServer(cfg, sugaredLogger)
	ginServer.RegisterHandlers("/api", authHandler, paymentHandler, measurementHandler, appointmentHandler)
	ginServer.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	ginServer.ListenAndServe()
}

func assertFatalError(logger *zap.SugaredLogger, err error, msg string) {
	if err == nil {
		return
	}
	sentry.CaptureException(err)
	sentry.Flush(time.Second * 2)
	logger.Fatalw(msg, "error", err)
}
