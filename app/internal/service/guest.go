package service

import (
	dm "ac_bot/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

type Guest interface {
	CreateInviteForGuest(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error)
	GuestRegForm(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error)
}

func (s *service) CreateInviteForGuest(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error) {
	// verify that the request is from an administrator
	card, err := s.isAdmin(ctx, data.TelegramUser.ID)
	if err != nil {
		return nil, fmt.Errorf("svc: CreateInviteForGuest: %w", err) // wrap!
	}

	// prepare data for db
	var invite dm.GuestInvite

	if err = json.Unmarshal(data.UserData, &invite); err != nil {
		return nil, fmt.Errorf("svc: CreateInviteForGuest: unmarshal: %s", err.Error())
	}

	// check for legal data
	if invite.Code == "" || invite.Title == "" {
		return nil, &dm.CliMsg{
			Message: "Invite code or org_name is empty.",
			Err:     fmt.Errorf("svc: CreateInviteForGuest: invite code is empty from user: %d", data.TelegramUser.ID),
		}
	}

	// db
	// use admin_id (guest_invite is more personalized)
	expire, err := s.storage.CreateInviteForGuest(ctx, invite, card.TelegramID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// must return expire time
			return nil, fmt.Errorf("svc: CreateInviteForGuest: not returned expite time: %w", err)
		}

		return nil, err
	}

	return []dm.SendMsg[dm.InKb]{
		{
			Text: fmt.Sprintf(dm.InviteForGuestAdmin, invite.Code, invite.Title, expire.Format(time.Stamp)),
		},
	}, nil
}

func (s *service) GuestRegForm(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error) {
	// prepare data for db
	var userData dm.Guest

	if err := json.Unmarshal(data.UserData, &userData); err != nil {
		return nil, fmt.Errorf("svc: GuestRegForm: unmarshal: %s", err.Error())
	}

	// org_id {string -> int}
	oId, err := strconv.ParseInt(userData.OrganizationID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("svc: GuestRegForm: parseInt: organization id: %s", err.Error())
	}

	// check invite
	invite, err := s.storage.GetGuestInviteCard(ctx, userData.InviteCode, int32(oId))
	if err != nil {
		// not found
		if err == pgx.ErrNoRows {
			return []dm.SendMsg[dm.InKb]{
				{
					Text: dm.InviteIDWrong,
				},
			}, nil
		}

		return nil, err
	}

	dto := dm.GuestDTO{
		TelegramID: data.TelegramUser.ID,
		FirstName:  userData.FirstName,
		LastName:   userData.LastName,
		InviteID:   invite.ID,
	}

	// db
	if err = s.storage.RegisterGuest(ctx, dto); err != nil {
		return nil, err
	}

	return []dm.SendMsg[dm.InKb]{
		{
			Text: dm.Created,
		},
		{
			Text: fmt.Sprintf(dm.InviteForGuestClient, invite.Expire.Format(time.Stamp)),
		},
	}, nil
}
