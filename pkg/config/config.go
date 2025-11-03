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
	}, nil
}
