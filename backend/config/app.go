package config

import (
	"backend/database"
	"backend/pkg/jwt"
	"backend/pkg/logger"
	"fmt"
	"time"

	"database/sql"
)

type AppConfig struct {
	Logger logger.Logger
	Config Config
	DB     *sql.DB
	JWT    jwt.JWT
}

func InitializeApp(envPath string) (*AppConfig, error) {
	// Load config
	cfg, err := LoadConfig(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	logger, err := logger.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize database
	db, err := database.Init(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run migrations
	// if err := database.RunMigrations(db, &models.User{}); err != nil {
	// 	return nil, fmt.Errorf("failed to initialize database: %w", err)
	// }

	// Initialize JWT
	jwtConfig := jwt.JWTConfig{
		SecretKey:     []byte(cfg.JWTSecretKey),
		TokenDuration: time.Duration(cfg.JWTExpirationDays) * time.Hour * 24,
	}
	jwt := jwt.New(jwtConfig)

	return &AppConfig{
		Logger: logger,
		Config: cfg,
		DB:     db,
		JWT:    jwt,
	}, nil
}