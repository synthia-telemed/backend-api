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
	server.AssertFatalError(sugaredLogger, err, "Failed to connect to database")

	patientDataStore, err := datastore.NewGormPatientDataStore(db)
	server.AssertFatalError(sugaredLogger, err, "Failed to create patient data store")
	creditCardDataStore, err := datastore.NewGormCreditCardDataStore(db)
	server.AssertFatalError(sugaredLogger, err, "Failed to create credit card data store")
	paymentDataStore, err := datastore.NewGormPaymentDataStore(db)
	server.AssertFatalError(sugaredLogger, err, "Failed to create payment data store")
	appointmentDataStore, err := datastore.NewGormAppointmentDataStore(db)
	server.AssertFatalError(sugaredLogger, err, "Failed to create appointment data store")
	notificationDataStore, err := datastore.NewGormNotificationDataStore(db)
	server.AssertFatalError(sugaredLogger, err, "Failed to create notification data store")

	hospitalSysClient := hospital.NewGraphQLClient(&cfg.HospitalClient)
	smsClient := sms.NewTwilioClient(&cfg.SMS)
	cacheClient := cache.NewRedisClient(&cfg.Cache)
	tokenService, err := token.NewGRPCTokenService(&cfg.Token)
	server.AssertFatalError(sugaredLogger, err, "Failed to create token service")
	paymentClient, err := payment.NewOmisePaymentClient(&cfg.Payment)
	server.AssertFatalError(sugaredLogger, err, "Failed to create payment client")
	realClock := clock.NewRealClock()

	// Handler
	authHandler := handler.NewAuthHandler(patientDataStore, hospitalSysClient, smsClient, cacheClient, tokenService, realClock, sugaredLogger)
	paymentHandler := handler.NewPaymentHandler(paymentClient, patientDataStore, creditCardDataStore, hospitalSysClient, paymentDataStore, realClock, sugaredLogger)
	appointmentHandler := handler.NewAppointmentHandler(patientDataStore, paymentDataStore, appointmentDataStore, hospitalSysClient, cacheClient, realClock, sugaredLogger)
	infoHandler := handler.NewInfoHandler(patientDataStore, hospitalSysClient, sugaredLogger)
	notificationHandler := handler.NewNotificationHandler(notificationDataStore, patientDataStore, sugaredLogger)

	ginServer := server.NewGinServer(cfg, sugaredLogger)
	ginServer.RegisterHandlers("/api", authHandler, paymentHandler, appointmentHandler, infoHandler, notificationHandler)
	ginServer.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	ginServer.ListenAndServe()
}
