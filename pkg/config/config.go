package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Env string

const (
	EnvProd  Env = "PROD"
	EnvDev   Env = "DEV"
	EnvDebug Env = "DEBUG"
)

type Config struct {
	ServerPort  int
	HTTPTimeout time.Duration
	Env         Env
	// DB
	DatabaseDSN       string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	AutoMigrate       bool
	// Logging
	AppLogLevel string
	// HTTP timeouts (optional overrides)
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	req := func(key string) (string, error) {
		if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
			return v, nil
		}
		return "", fmt.Errorf("required env var %s is missing", key)
	}

	// SERVER_PORT
	portStr, err := req("SERVER_PORT")
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		return nil, fmt.Errorf("invalid SERVER_PORT=%q", portStr)
	}

	// HTTP_TIMEOUT (Go duration: 500ms, 2s, 1m)
	toStr, err := req("HTTP_TIMEOUT")
	if err != nil {
		return nil, err
	}
	timeout, err := time.ParseDuration(toStr)
	if err != nil || timeout <= 0 {
		return nil, fmt.Errorf("invalid HTTP_TIMEOUT=%q (use Go duration, e.g. 10s, 500ms, 1m)", toStr)
	}

	// ENV
	envStr, err := req("ENV")
	if err != nil {
		return nil, err
	}
	env := Env(strings.ToUpper(strings.TrimSpace(envStr)))
	switch env {
	case EnvProd, EnvDev, EnvDebug:
	default:
		return nil, errors.New(`ENV must be one of: PROD|DEV|DEBUG`)
	}

	return &Config{
		ServerPort:  port,
		HTTPTimeout: timeout,
		Env:         env,
		// defaults for new fields
		DatabaseDSN:       os.Getenv("DATABASE_DSN"),
		DBMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", time.Hour),
		AutoMigrate:       getEnvBool("AUTO_MIGRATE", false),
		AppLogLevel:       getEnvString("APP_LOG_LEVEL", "info"),
		HTTPReadTimeout:   getEnvDuration("HTTP_READ_TIMEOUT", timeout),
		HTTPWriteTimeout:  getEnvDuration("HTTP_WRITE_TIMEOUT", timeout),
		HTTPIdleTimeout:   getEnvDuration("HTTP_IDLE_TIMEOUT", 2*timeout),
	}, nil
}

func getEnvString(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes":
			return true
		case "0", "false", "no":
			return false
		}
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
