package ping

import (
	"context"

	d "github.com/antalkon/Go_prod_tmpl/internal/domain/ping"
)

type Repository interface {
	Create(ctx context.Context, p *d.Ping) error
}
