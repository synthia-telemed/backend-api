package main

import (
	"github.com/getsentry/sentry-go"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/synthia-telemed/backend-api/cmd/doctor-api/docs"
	"github.com/synthia-telemed/backend-api/cmd/doctor-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/clock"
	"github.com/synthia-telemed/backend-api/pkg/config"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/id"
	"github.com/synthia-telemed/backend-api/pkg/logger"
	"github.com/synthia-telemed/backend-api/pkg/notification"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/token"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

// @title           Synthia Doctor Backend API
// @version         1.0.0
// @description     This is a Synthia doctor backend API.
// @accept json
// @produce json
// @BasePath  /doctor/api

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

	doctorDataStore, err := datastore.NewGormDoctorDataStore(db)
	server.AssertFatalError(sugaredLogger, err, "Failed to create doctor data store")
	patientDataStore, err := datastore.NewGormPatientDataStore(db)
	server.AssertFatalError(sugaredLogger, err, "Failed to create patient data store")
	appointmentDataStore, err := datastore.NewGormAppointmentDataStore(db)
	server.AssertFatalError(sugaredLogger, err, "Failed to create appointment data store")
	notificationDataStore, err := datastore.NewGormNotificationDataStore(db)
	server.AssertFatalError(sugaredLogger, err, "Failed to create notification data store")

	hospitalSysClient := hospital.NewGraphQLClient(&cfg.HospitalClient)
	cacheClient := cache.NewRedisClient(&cfg.Cache)
	realClock := clock.NewRealClock()
	idGenerator := id.NewNanoID()
	tokenService, err := token.NewGRPCTokenService(&cfg.Token)
	server.AssertFatalError(sugaredLogger, err, "Failed to create token service")
	notificationClient, err := notification.NewRabbitMQNotificationClient(&cfg.Notification)
	server.AssertFatalError(sugaredLogger, err, "Failed to create rabbitmq notification client")

	// Handlers
	authHandler := handler.NewAuthHandler(hospitalSysClient, tokenService, doctorDataStore, sugaredLogger)
	appointmentHandler := handler.NewAppointmentHandler(appointmentDataStore, patientDataStore, doctorDataStore, notificationDataStore, hospitalSysClient, cacheClient, realClock, idGenerator, notificationClient, sugaredLogger)

	ginServer := server.NewGinServer(cfg, sugaredLogger)
	ginServer.RegisterHandlers("/api", authHandler, appointmentHandler)
	ginServer.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	ginServer.ListenAndServe()
	server.AssertFatalError(sugaredLogger, notificationClient.Close(), "Failed to close rabbitmq connection")
}
