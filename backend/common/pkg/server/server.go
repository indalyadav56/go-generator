package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Router() *gin.Engine
}

type Config struct {
	Port            int           `json:"port" yaml:"port"`
	ReadTimeout     time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout" yaml:"write_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" yaml:"shutdown_timeout"`
	Mode            string        `json:"mode" yaml:"mode"`
	MetricsEnabled  bool          `json:"metrics_enabled" yaml:"metrics_enabled"`
	CorsEnabled     bool          `json:"cors_enabled" yaml:"cors_enabled"`
	RateLimit       RateLimit     `json:"rate_limit" yaml:"rate_limit"`
}

type RateLimit struct {
	Enable   bool          `json:"enable" yaml:"enable"`
	Requests int           `json:"requests" yaml:"requests"`
	Duration time.Duration `json:"duration" yaml:"duration"`
}

func DefaultConfig() Config {
	return Config{
		Port:            8080,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 5 * time.Second,
		Mode:            gin.DebugMode,
		// Mode:            gin.TestMode,
		// Mode:            gin.ReleaseMode,
		MetricsEnabled: true,
		CorsEnabled:    true,
		RateLimit: RateLimit{
			Enable:   true,
			Requests: 100,
			Duration: time.Minute,
		},
	}
}

type server struct {
	cfg    Config
	router *gin.Engine
	srv    *http.Server
}

func New(config ...Config) (Server, error) {
	var cfg Config

	if len(config) == 0 {
		cfg = DefaultConfig()
	} else {
		cfg = config[0]
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	gin.SetMode(cfg.Mode)
	router := gin.Default()

	s := &server{
		cfg:    cfg,
		router: router,
		srv: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s, nil
}

func validateConfig(cfg Config) error {
	if cfg.Port < 1 || cfg.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	if cfg.ReadTimeout <= 0 {
		return errors.New("read timeout must be positive")
	}
	if cfg.WriteTimeout <= 0 {
		return errors.New("write timeout must be positive")
	}
	return nil
}

func (s *server) setupMiddleware() {
	if s.cfg.CorsEnabled {
		s.router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000", "https://example.com"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"X-Request-ID"},
			AllowCredentials: true,
		}))
	}
	// if s.cfg.RateLimit.Enable {
	// 	s.router.Use(s.rateLimitMiddleware())
	// }
	s.router.Use(requestid.New())
	// if s.cfg.MetricsEnabled {
	// 	s.router.Use(s.metricsMiddleware())
	// }
}

func (s *server) setupRoutes() {
	s.setupHealthCheck()
	if s.cfg.MetricsEnabled {
		s.setupMetrics()
	}
}

func (s *server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-sigCh:
		return s.Shutdown(ctx)
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()

	return s.srv.Shutdown(shutdownCtx)
}

func (s *server) Router() *gin.Engine {
	return s.router
}

func (s *server) setupHealthCheck() {
	health := s.router.Group("/health")
	{
		health.GET("", s.handleServerHeath())
		health.GET("/live", s.handleLiveness())
		health.GET("/ready", s.handleReadiness())
	}
}

func (s *server) handleLiveness() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func (s *server) handleReadiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func (s *server) handleServerHeath() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	}
}

func (s *server) setupMetrics() {
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
