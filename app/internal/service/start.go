package service

import (
	dm "ac_bot/internal/domain"
	"ac_bot/internal/helper"
	"log"

	"context"
)

func (s *service) Start(ctx context.Context, telegramID int64) (*dm.SendMsg[dm.StKb], error) {
	// create keyboard
	kb := s.startKb(ctx, telegramID)

	return &dm.SendMsg[dm.StKb]{
		Text:     "Welcome to \"The Academie One\" bot channel. Please press the button below.",
		Keyboard: kb,
	}, nil
}

func (s *service) startKb(ctx context.Context, telegramID int64) dm.StKb {
	// get organization list for menu
	orgs, err := s.getOrgListF(ctx, "%d=%s", ',')
	if err != nil {
		log.Printf("svc: startKb: getOrgListF: %s", err)
	}

	if _, err := s.storage.IsStudent(ctx, telegramID); err == nil {
		return helper.StudentStKb(s.webApp)
	}

	if card, err := s.storage.IsAdmin(ctx, telegramID); err == nil {
		if card.RoleID == dm.SecurityRole {
			return helper.SecurityStKb(s.webApp)
		}

		roles, err := s.getRoleListF(ctx, "%d=%s", ',', true)
		if err != nil {
			log.Printf("svc: startKb: getRoleListF: %s", err)
		}

		kb := helper.GuestStKb(s.webApp, orgs)
		kb = append(kb, helper.AdminStKbButton(s.webApp, orgs, roles))
		return kb
	}

	return helper.GuestStKb(s.webApp, orgs)
}
