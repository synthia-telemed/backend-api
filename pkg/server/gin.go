package server

import (
	"context"
	"errors"
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/synthia-telemed/backend-api/pkg/config"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*gin.Engine
	logger *zap.SugaredLogger
	addr   string
}

type Enum interface {
	IsValid() bool
}

func RegisterValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("enum", func(fl validator.FieldLevel) bool {
			value := fl.Field().Interface().(Enum)
			return value.IsValid()
		})
	}
}

func NewGinServer(cfg *config.Config, logger *zap.SugaredLogger) *Server {
	gin.SetMode(cfg.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/api/healthcheck"}}))
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	router.GET("/api/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"timestamp": time.Now(),
		})
	})
	RegisterValidator()

	return &Server{
		Engine: router,
		logger: logger,
		addr:   fmt.Sprintf(":%d", cfg.Port),
	}
}

func (s Server) ListenAndServe() {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.Engine,
	}

	go func() {
		s.logger.Infow("Starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			s.logger.Infof("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Fatalw("Server forced to shutdown", "error", err)
	}
	s.logger.Info("Server exiting")
}

func (s Server) RegisterHandlers(prefix string, handlers ...Handler) {
	for _, h := range handlers {
		h.Register(s.Group(prefix))
	}
}
