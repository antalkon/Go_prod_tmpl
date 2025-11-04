package ping

import "time"

type Ping struct {
	ID        int64     `db:"id" json:"id"`
	Message   string    `db:"message" json:"message"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
