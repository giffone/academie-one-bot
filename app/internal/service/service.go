package service

import (
	"ac_bot/config"
	dm "ac_bot/internal/domain"
	pg "ac_bot/internal/repository/postgres"
	"ac_bot/internal/service/tgs"
	"context"
	"log"
)

type Service interface {
	Admin
	Student
	Guest
	QR
	Lists

	Start(ctx context.Context, telegramID int64) (*dm.SendMsg[dm.StKb], error)
	Whoami(ctx context.Context, user dm.TelegramUser) *dm.SendMsg[any]
}

type TelegramService interface {
	Register(tSvc tgs.TelegramService)
}

func New(ctx context.Context, storage pg.Storage, cfg *config.Config) (Service, TelegramService) {
	s := service{
		storage: storage,
		debug:   cfg.Debug,
		webApp: &dm.WebApp{
			QRText: cfg.QRText,
			Addr:   cfg.WebAppAddr,
			Paths:  make(map[string]string),
		},
		ctx: ctx,
		event: event{
			ch: make(chan dm.NewNotification, 1),
		},
	}

	u := cfg.WebAppAddr

	if u[len(u)-1] == '/' {
		u = u[:len(u)-1]
	}

	s.webApp.Paths[dm.WebAppPathAdmin] = u + dm.WebAppPathAdmin
	s.webApp.Paths[dm.WebAppPathRegForm] = u + dm.WebAppPathRegForm
	s.webApp.Paths[dm.WebAppPathIntra] = u + dm.WebAppPathIntra
	s.webApp.Paths[dm.WebAppPathQrScan] = u + dm.WebAppPathQrScan

	for key, value := range s.webApp.Paths {
		log.Printf("registered paths: key: \"%s\", value: \"%s\"", key, value)
	}

	// start listening to events and sending notifications [admin, security]
	go s.eventListener(s.event.ch)

	return &s, &s.event
}

type service struct {
	storage pg.Storage
	debug   bool
	webApp  *dm.WebApp
	ctx     context.Context
	event   event
}

type event struct {
	ch   chan dm.NewNotification
	tSvc tgs.TelegramService
}

func (e *event) Register(tSvc tgs.TelegramService) {
	e.tSvc = tSvc
}
