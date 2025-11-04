package ping

import (
	"context"
	"fmt"

	d "github.com/antalkon/Go_prod_tmpl/internal/domain/ping"
	repo "github.com/antalkon/Go_prod_tmpl/internal/repository/ping"
	lg "github.com/antalkon/Go_prod_tmpl/pkg/logger"
	"go.uber.org/zap"
)

type Service interface {
	Ping(ctx context.Context, message string) (*d.Ping, error)
}

type svc struct {
	repo repo.Repository
	log  *zap.Logger
}

func NewService(r repo.Repository, log *zap.Logger) Service {
	return &svc{repo: r, log: log}
}

func (s *svc) Ping(ctx context.Context, message string) (*d.Ping, error) {
	if message == "" {
		message = "pong"
	}
	p := &d.Ping{Message: message}
	reqLog := lg.FromContext(ctx, s.log)
	if err := s.repo.Create(ctx, p); err != nil {
		reqLog.Error("service.ping.create.error", zap.Error(err))
		return nil, fmt.Errorf("create ping: %w", err)
	}
	reqLog.Info("service.ping.created", zap.Int64("id", p.ID))
	return p, nil
}
