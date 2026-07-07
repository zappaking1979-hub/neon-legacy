package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         string
	DatabaseURL  string
	RedisURL     string
	SessionKey   string
	CorsOrigin   string
	CentrifugoURL string
	CentrifugoKey string
}

type ConfigResult struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	Auth  AuthConfig
	Cors  CorsConfig
}

type AppConfig struct {
	Port        string
	Environment string
}

type DBConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type AuthConfig struct {
	SessionKey     string
	SessionTTL     time.Duration
	BcryptCost     int
	RateLimit      int
	RateLimitWindow time.Duration
}

type CorsConfig struct {
	AllowedOrigin string
}

func Load() *ConfigResult {
	return &ConfigResult{
		App: AppConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		DB: DBConfig{
			URL:             getEnv("DATABASE_URL", "postgres://neon:neon@localhost:5432/neonlegacy?sslmode=disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Auth: AuthConfig{
			SessionKey:      getEnv("SESSION_KEY", "change-me-in-production"),
			SessionTTL:      getEnvDuration("SESSION_TTL", 24*time.Hour),
			BcryptCost:      getEnvInt("BCRYPT_COST", 12),
			RateLimit:       getEnvInt("RATE_LIMIT", 30),
			RateLimitWindow: getEnvDuration("RATE_LIMIT_WINDOW", 15*time.Second),
		},
		Cors: CorsConfig{
			AllowedOrigin: getEnv("CORS_ORIGIN", "http://localhost:8080"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
