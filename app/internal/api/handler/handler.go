package handler

import (
	"ac_bot/config"
	"ac_bot/internal/service"
)

type Handler struct {
	svc       service.Service
	botID     int64
}

func New(cfg *config.Config, svc service.Service) *Handler {
	return &Handler{
		svc:       svc,
		botID:     cfg.BotID,
	}
}
