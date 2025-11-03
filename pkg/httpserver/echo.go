package httpserver

import (
	"fmt"
	"net/http"

	"github.com/antalkon/Go_prod_tmpl/pkg/config"
	"github.com/labstack/echo/v4"
)

type RegisterFunc func(e *echo.Echo)

type Server struct {
	E      *echo.Echo
	Srv    *http.Server
	Config *config.Config
}

func New(cfg *config.Config, register RegisterFunc) (*Server, error) {
	e := echo.New()
	e.HideBanner, e.HidePort = true, true

	if register != nil {
		register(e)
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:           e,
		ReadHeaderTimeout: cfg.HTTPTimeout,
		ReadTimeout:       cfg.HTTPTimeout,
		WriteTimeout:      cfg.HTTPTimeout,
		IdleTimeout:       2 * cfg.HTTPTimeout,
	}

	return &Server{E: e, Srv: srv, Config: cfg}, nil
}

func (s *Server) Addr() string { return s.Srv.Addr }
