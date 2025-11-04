package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	pingrepo "github.com/antalkon/Go_prod_tmpl/internal/repository/ping"
	"github.com/antalkon/Go_prod_tmpl/internal/repository/postgres"
	pingsvc "github.com/antalkon/Go_prod_tmpl/internal/service/ping"
	"github.com/antalkon/Go_prod_tmpl/internal/transport/rest/handlers"
	"github.com/antalkon/Go_prod_tmpl/internal/transport/rest/middlewares"
	"github.com/antalkon/Go_prod_tmpl/internal/transport/rest/router"
	"github.com/antalkon/Go_prod_tmpl/pkg/config"
	"github.com/antalkon/Go_prod_tmpl/pkg/db"
	"github.com/antalkon/Go_prod_tmpl/pkg/httpserver"
	applog "github.com/antalkon/Go_prod_tmpl/pkg/logger"
	"github.com/antalkon/Go_prod_tmpl/pkg/migrations"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type App struct {
	cfg *config.Config
	log *zap.Logger
	srv *httpserver.Server
	db  *sqlx.DB
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

	// DB
	dbConn, err := db.New(cfg.DatabaseDSN, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	if err != nil {
		return nil, err
	}

	// AutoMigrate (use pkg/migrations)
	if cfg.AutoMigrate {
		if err := migrations.MigrateUp("./migrations", cfg.DatabaseDSN, log); err != nil {
			return nil, err
		}
	}

	// create component-named loggers
	httpLog := applog.Named(log, "http")
	dbLog := applog.Named(log, "db")
	svcLog := applog.Named(log, "svc")

	// wire repository / service / handlers
	pingRepo := postgres.NewPingRepo(dbConn, dbLog)
	var pingRepoIface pingrepo.Repository = pingRepo
	pingSvc := pingsvc.NewService(pingRepoIface, svcLog)

	srv, err := httpserver.New(cfg, func(e *echo.Echo) {
		mw := middlewares.New(httpLog)
		mw.Use(e)

		ping := handlers.NewPingHandler(pingSvc, httpLog)
		router.Register(e, ping)
	})
	if err != nil {
		return nil, err
	}

	return &App{cfg: cfg, log: log, srv: srv, db: dbConn}, nil
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

	if a.db != nil {
		_ = a.db.Close()
	}

	a.log.Info("server stopped gracefully")
	_ = a.log.Sync()
	return nil
}
