package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PingHandler struct {
	log *zap.Logger
}

func NewPingHandler(log *zap.Logger) *PingHandler { return &PingHandler{log: log} }

func (h *PingHandler) Ping(c echo.Context) error {
	h.log.Debug("ping", zap.String("path", c.Path()), zap.String("ip", c.RealIP()))
	return c.JSON(http.StatusOK, map[string]any{
		"message": "pong",
		"ts":      time.Now().UTC().Format(time.RFC3339Nano),
	})
}
