package service

import (
	dm "ac_bot/internal/domain"
	"context"
	"fmt"
	"log"
	"time"
)

func (s *service) asGuestQR(ctx context.Context, chatID int64, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error) {
	log.Println("asGuestQR")
	regs, err := s.storage.GetGuestRegistrations(ctx, data.TelegramUser.ID)
	if err != nil {
		return nil, err
	}

	if len(regs) == 0 {
		return nil, nil
	}

	// found
	if len(regs) == 1 {
		reg := regs[0]

		// check invite expire
		if reg.Invite.Expire.Before(time.Now()) {
			return []dm.SendMsg[dm.InKb]{
				{
					Text: fmt.Sprintf(dm.InviteExpired, reg.Invite.Expire.Format(time.Stamp), reg.Invite.Organization.Title),
				},
			}, nil
		}

		// register guest entrance
		return s.asGuestEntr(ctx, data.TelegramUser.ID, reg.ID)
	}

	// guest is registered with more than one organization
	// he can select one of them using the inline keyboard that we have to send
	// after pressing the keyboard the callback handler will be called
	return s.asGuestCallback(ctx, data.TelegramUser.ID, chatID, regs)
}

func (s *service) asGuestEntr(ctx context.Context, telegramID int64, regID int32) ([]dm.SendMsg[dm.InKb], error) {
	err := s.storage.CreateEntranceGuest(ctx, regID)
	if err != nil {
		return nil, err
	}

	// send info to security (post) chat
	go s.createNotification(telegramID, dm.GuestRole, regID)

	// registered
	return []dm.SendMsg[dm.InKb]{
		{
			Text: fmt.Sprintf(dm.QRSuccessGuest),
		},
	}, nil
}

func (s *service) asGuestCallback(ctx context.Context, telegramID, chatID int64, regs []dm.GuestRegDTO) ([]dm.SendMsg[dm.InKb], error) {
	// prepare keyboard for client
	kb := make(dm.InKb, len(regs))

	for i, v := range regs {
		kb[i] = []dm.InKbButton{
			{
				Text:         fmt.Sprintf("%s [%s]", v.Invite.Organization.Title, v.Invite.Admin.Name), // org_name and employee
				CallbackData: fmt.Sprintf("button=%d", v.ID),                                           // registration_id
			},
		}
	}

	// preparing information for callback and recording to db
	cl := dm.CallbackDataDTO{
		TelegramID:  telegramID,
		ChatID:      chatID,
		ServiceName: dm.InlEntranceGuest,
		RoleID:      dm.GuestRole,
	}

	if err := s.storage.CreateCallbackBuffer(ctx, cl); err != nil {
		return nil, fmt.Errorf("svc: asGuestCallback: %w", err)
	}

	return []dm.SendMsg[dm.InKb]{
		{
			Text:     fmt.Sprintf(dm.MultipleInvitations, len(regs)),
			Keyboard: kb,
		},
	}, nil
}
