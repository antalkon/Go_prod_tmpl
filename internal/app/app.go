package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/antalkon/Go_prod_tmpl/internal/transport/rest/handlers"
	"github.com/antalkon/Go_prod_tmpl/internal/transport/rest/middlewares"
	"github.com/antalkon/Go_prod_tmpl/internal/transport/rest/router"
	"github.com/antalkon/Go_prod_tmpl/pkg/config"
	"github.com/antalkon/Go_prod_tmpl/pkg/httpserver"
	applog "github.com/antalkon/Go_prod_tmpl/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type App struct {
	cfg *config.Config
	log *zap.Logger
	srv *httpserver.Server
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	log, err := applog.New(string(cfg.Env))
	if err != nil {
		return nil, err
	}

	srv, err := httpserver.New(cfg, func(e *echo.Echo) {
		mw := middlewares.New(log)
		mw.Use(e)

		ping := handlers.NewPingHandler(log)
		router.Register(e, ping)
	})
	if err != nil {
		return nil, err
	}

	return &App{cfg: cfg, log: log, srv: srv}, nil
}

func (a *App) Run() error {
	a.log.Info("app starting",
		zap.String("env", string(a.cfg.Env)),
		zap.String("addr", a.srv.Addr()),
		zap.Duration("timeout", a.cfg.HTTPTimeout),
	)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if err := a.srv.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.log.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.srv.Srv.Shutdown(shutdownCtx); err != nil {
		a.log.Error("server shutdown error", zap.Error(err))
		return err
	}

	a.log.Info("server stopped gracefully")
	_ = a.log.Sync()
	return nil
}
