package main

import (
	"ac_bot/config"
	"ac_bot/internal/server"
	"context"
)

func main() {
	ctx := context.Background()

	cfg := config.New()
	cfg.Print()

	s := server.NewServer(ctx, cfg)
	s.Start(ctx)
}
