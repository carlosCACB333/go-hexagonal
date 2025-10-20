package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	API      APIConfig
	DB       DBConfig
	RabbitMQ RabbitMQConfig
	App      AppConfig
}

type APIConfig struct {
	Port string
	Host string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string

	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
}

type AppConfig struct {
	LogLevel    string
	Environment string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.ReadInConfig()

	return &Config{
		API: APIConfig{
			Port: getEnvOrDefault("API_PORT", "8080"),
			Host: getEnvOrDefault("API_HOST", "0.0.0.0"),
		},
		DB: DBConfig{
			Host:            getEnvOrDefault("DB_HOST", "localhost"),
			Port:            getEnvOrDefault("DB_PORT", "5432"),
			User:            getEnvOrDefault("DB_USER", "postgres"),
			Password:        getEnvOrDefault("DB_PASSWORD", "postgres"),
			Name:            getEnvOrDefault("DB_NAME", "backend_hex_cqrs"),
			SSLMode:         getEnvOrDefault("DB_SSL_MODE", "disable"),
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: 10 * time.Minute,
		},
		RabbitMQ: RabbitMQConfig{
			Host:     getEnvOrDefault("RABBITMQ_HOST", "localhost"),
			Port:     getEnvOrDefault("RABBITMQ_PORT", "5672"),
			User:     getEnvOrDefault("RABBITMQ_USER", "guest"),
			Password: getEnvOrDefault("RABBITMQ_PASSWORD", "guest"),
			VHost:    getEnvOrDefault("RABBITMQ_VHOST", "/"),
		},
		App: AppConfig{
			LogLevel:    getEnvOrDefault("LOG_LEVEL", "info"),
			Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return viper.GetString(key)
}
