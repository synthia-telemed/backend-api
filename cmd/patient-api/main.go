package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/config"
	"github.com/synthia-telemed/backend-api/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	gin.SetMode(cfg.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/api/healthcheck"}}))
	router.GET("/api/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"timestamp": time.Now(),
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	go func() {
		sugaredLogger.Infow("Starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			sugaredLogger.Fatalw("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugaredLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugaredLogger.Fatalw("Server forced to shutdown", "error", err)
	}
	sugaredLogger.Info("Server exiting")
}
