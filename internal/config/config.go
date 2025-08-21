package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config содержит все настройки приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
}

// ServerConfig настройки HTTP сервера
type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

// DatabaseConfig настройки подключения к базе данных
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// CacheConfig настройки кеширования
type CacheConfig struct {
	Enabled    bool
	MaxEntries int
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8081"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 10),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "wb_user"),
			Password: getEnv("DB_PASSWORD", "wb_password"),
			DBName:   getEnv("DB_NAME", "wildberries_orders"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Cache: CacheConfig{
			Enabled:    getEnvAsBool("CACHE_ENABLED", true),
			MaxEntries: getEnvAsInt("CACHE_MAX_ENTRIES", 1000),
		},
	}
}

// GetDSN возвращает строку подключения к PostgreSQL
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// GetConnectionString возвращает строку подключения в формате URL
func (c *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

// Вспомогательные функции для работы с переменными окружения

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
} 