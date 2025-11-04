package handlers

import (
	"net/http"

	"github.com/antalkon/Go_prod_tmpl/internal/service/ping"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PingHandler struct {
	svc ping.Service
	log *zap.Logger
}

func NewPingHandler(svc ping.Service, log *zap.Logger) *PingHandler {
	return &PingHandler{svc: svc, log: log}
}

// Ping godoc
// @Summary Ping
// @Description persist a ping record and return it
// @Tags health
// @Produce json
// @Param message query string false "custom message"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/ping [get]
func (h *PingHandler) Ping(c echo.Context) error {
	ctx := c.Request().Context()
	h.log.Debug("http.request", zap.String("path", c.Path()), zap.String("ip", c.RealIP()))
	p, err := h.svc.Ping(ctx, c.QueryParam("message"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, p)
}
