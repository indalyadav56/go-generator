package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	ErrInvalidConfig = errors.New("invalid database configuration")
	ErrConnection    = errors.New("database connection error")
	ErrMigration     = errors.New("migration error")
)

type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

func (c *Config) Validate() error {
	if c.Host == "" || c.Port == "" || c.User == "" ||
		c.Name == "" || c.SSLMode == "" {
		return ErrInvalidConfig
	}
	return nil
}

func DefaultConfig() *Config {
	return &Config{
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxLifetime:  5 * time.Minute,
		SSLMode:      "disable",
	}
}

type DB struct {
	*sql.DB
}

func New(cfg *Config) (*DB, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnection, err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%w: ping failed: %v", ErrConnection, err)
	}

	return &DB{db}, nil
}

// ApplyMigrations runs database migrations from the specified directory
func (db *DB) ApplyMigrations(ctx context.Context, migrationDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("%w: failed to set dialect: %v", ErrMigration, err)
	}

	if err := goose.Up(db.DB, migrationDir); err != nil {
		return fmt.Errorf("%w: failed to apply migrations: %v", ErrMigration, err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}
