package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost            string `mapstructure:"DB_HOST" validate:"required"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	RedisAddr         string `mapstructure:"REDIS_ADDR"`
	JWTSecretKey      string `mapstructure:"JWT_SECRET_KEY" validate:"required"`
	JWTExpirationDays int    `mapstructure:"JWT_EXPIRATION_DAYS" validate:"required,gt=0"`
	ServerPort        string `mapstructure:"SERVER_PORT" validate:"required"`
	LogLevel          string `mapstructure:"LOG_LEVEL"`
	LogFilePath       string `mapstructure:"LOG_FILE_PATH"`
	AppName           string `mapstructure:"APP_NAME"`
	MongoDB           string `mapstructure:"MONGO_DB"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	if err := setupViper(path); err != nil {
		return config, fmt.Errorf("error setting up viper: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("error unmarshalling config: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return config, fmt.Errorf("error validating config: %w", err)
	}

	return config, nil
}

func setupViper(path string) error {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.SetConfigFile(".env")
	}

	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("JWT_EXPIRATION_DAYS", 7)
	viper.SetDefault("SERVER_PORT", 8080)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	viper.AutomaticEnv()

	return nil
}

func validateConfig(config Config) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}
