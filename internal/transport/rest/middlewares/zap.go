package middlewares

import (
	"time"

	"github.com/antalkon/Go_prod_tmpl/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Set struct{ log *zap.Logger }

func New(log *zap.Logger) *Set { return &Set{log: log} }

func (s *Set) Use(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(s.zapRequestLogger())
}

func (s *Set) zapRequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req, res := c.Request(), c.Response()

			// create request-scoped logger and put into context
			reqLogger := s.log.With(
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("ip", c.RealIP()),
				zap.String("req_id", res.Header().Get(echo.HeaderXRequestID)),
			)
			ctx := logger.WithContext(req.Context(), reqLogger)
			c.SetRequest(c.Request().WithContext(ctx))

			err := next(c)

			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", res.Status),
				zap.Duration("latency", time.Since(start)),
				zap.String("ip", c.RealIP()),
				zap.String("req_id", res.Header().Get(echo.HeaderXRequestID)),
			}
			if err != nil {
				reqLogger.Error("http_request", append(fields, zap.Error(err))...)
				return err
			}
			reqLogger.Info("http_request", fields...)
			_ = ctx
			return nil
		}
	}
}
