package postgres

import (
	"ac_bot/config"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

//	# Example URL
//	postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10

func NewPostgres(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	log.Println("[postgres-pool] init...")

	pool, err := pgxpool.New(ctx, cfg.DBAddr)
	if err != nil {
		log.Fatalf("[postgres-pool] init error: %s", err)
	}

	log.Println("[postgres-pool] check conn")

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("[postgres-pool] check conn error: %s", err)
	}

	defer conn.Release()
	log.Println("[postgres-pool] init done")

	log.Println("[postgres-pool] set time zone: Asia/Almaty")
	if _, err = conn.Exec(ctx, "SET TIME ZONE 'Asia/Almaty'"); err != nil {
		log.Fatalf("[postgres-pool] set time zone error: %s", err)
	}

	return pool
}
