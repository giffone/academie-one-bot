package service

import (
	dm "ac_bot/internal/domain"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *service) asStudentQR(ctx context.Context, chatID int64, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error) {
	log.Println("asStudentQR")
	regs, err := s.storage.GetStudentRegistrations(ctx, data.TelegramUser.ID)
	if err != nil {
		return nil, err
	}

	if len(regs) == 0 {
		return nil, nil
	}

	// found
	if len(regs) == 1 {
		reg := regs[0]

		// check approve invite
		if !reg.Approved {
			return []dm.SendMsg[dm.InKb]{
				{
					Text: fmt.Sprintf(dm.RegNotConfirmed, reg.Invite.Organization.Title),
				},
			}, nil
		}

		// check invite expire
		if reg.Invite.Expire.Before(time.Now()) {
			return []dm.SendMsg[dm.InKb]{
				{
					Text: fmt.Sprintf(dm.InviteExpired, reg.Invite.Expire.Format(time.Stamp), reg.Invite.Organization.Title),
				},
			}, nil
		}

		// register student entrance
		return s.asStudentEntr(ctx, data.TelegramUser.ID, reg.Invite.ID)
	}

	// student is registered with more than one organization
	// he can select one of them using the inline keyboard that we have to send
	// after pressing the keyboard the callback handler will be called
	return s.asStudentCallback(ctx, data.TelegramUser.ID, chatID, regs)
}

func (s *service) asStudentEntr(ctx context.Context, telegramID int64, inviteID int32) ([]dm.SendMsg[dm.InKb], error) {
	expire, err := s.storage.CreateEntranceStudent(ctx, telegramID, inviteID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// not returned expire time
			return nil, fmt.Errorf("svc: asStudentEntr: not returned expite time: %w", err)
		}

		return nil, err
	}

	// send info to security (post) chat
	go s.createNotification(telegramID, dm.StudentRole, inviteID)

	// registered
	return []dm.SendMsg[dm.InKb]{
		{
			Text: fmt.Sprintf(dm.QRSuccessStudent, expire.Format(time.Stamp)),
		},
	}, nil
}

func (s *service) asStudentCallback(ctx context.Context, telegramID, chatID int64, regs []dm.StudentRegDTO) ([]dm.SendMsg[dm.InKb], error) {
	// prepare keyboard for client
	kb := make(dm.InKb, len(regs))

	for i, v := range regs {
		kb[i] = []dm.InKbButton{
			{
				Text:         v.Invite.Organization.Title,           // only org_name
				CallbackData: fmt.Sprintf("button=%d", v.Invite.ID), // invite_id
			},
		}
	}

	// preparing information for callback and recording to db
	cl := dm.CallbackDataDTO{
		TelegramID:  telegramID,
		ChatID:      chatID,
		ServiceName: dm.InlEntranceStudent,
		RoleID:      dm.StudentRole,
	}

	if err := s.storage.CreateCallbackBuffer(ctx, cl); err != nil {
		return nil, fmt.Errorf("svc: asStudentCallback: %w", err)
	}

	return []dm.SendMsg[dm.InKb]{
		{
			Text:     fmt.Sprintf(dm.MultipleInvitations, len(regs)),
			Keyboard: kb,
		},
	}, nil
}
