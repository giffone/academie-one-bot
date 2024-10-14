package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Administrator
	Guest
	Student
	CallbackBuffer
	Lists
	Notification
}

func New(pool *pgxpool.Pool) Storage {
	return &storage{pool: pool}
}

type storage struct {
	pool *pgxpool.Pool
}
