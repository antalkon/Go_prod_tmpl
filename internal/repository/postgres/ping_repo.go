package postgres

import (
	"context"
	"time"

	d "github.com/antalkon/Go_prod_tmpl/internal/domain/ping"
	lg "github.com/antalkon/Go_prod_tmpl/pkg/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PingRepo struct {
	db  *sqlx.DB
	log *zap.Logger
}

func NewPingRepo(db *sqlx.DB, log *zap.Logger) *PingRepo {
	return &PingRepo{db: db, log: log}
}

const insertPingSQL = `INSERT INTO pings (message) VALUES ($1) RETURNING id, created_at`

func (r *PingRepo) Create(ctx context.Context, p *d.Ping) error {
	start := time.Now()
	reqLog := lg.FromContext(ctx, r.log)
	reqLog.Debug("repo.ping.create.start", zap.String("message", p.Message))

	row := r.db.QueryRowxContext(ctx, insertPingSQL, p.Message)
	if err := row.Scan(&p.ID, &p.CreatedAt); err != nil {
		reqLog.Error("repo.ping.create.error", zap.Error(err))
		return err
	}

	reqLog.Info("repo.ping.create.ok",
		zap.Int64("id", p.ID),
		zap.Duration("latency", time.Since(start)),
	)
	return nil
}
