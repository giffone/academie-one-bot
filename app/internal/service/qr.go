package service

import (
	dm "ac_bot/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type QR interface {
	RegistrationQR(ctx context.Context, chatID int64, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error)
	RegistrationQRCallBack(ctx context.Context, telegramID, chatID int64, buttonID int32) ([]dm.SendMsg[dm.InKb], error)
}

func (s *service) RegistrationQR(ctx context.Context, chatID int64, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error) {
	var userData dm.QRData

	err := json.Unmarshal(data.UserData, &userData)
	if err != nil {
		return nil, fmt.Errorf("RegistrationQR: unmarshal: %s", err.Error())
	}

	if userData.QRText != s.webApp.QRText {
		return []dm.SendMsg[dm.InKb]{
			{
				Text: "Unknown QR code.",
				// Keyboard: nil,
			},
		}, nil
	}

	// check registration as student
	msg, err := s.asStudentQR(ctx, chatID, data)
	if err != nil {
		return nil, err
	}
	// found
	if msg != nil {
		return msg, nil
	}

	// check registration as guest
	msg, err = s.asGuestQR(ctx, chatID, data)
	if err != nil {
		return nil, err
	}
	// found
	if msg != nil {
		return msg, nil
	}

	// send info to security (post) chat
	go s.createNotification(0, 0, 0) // unknown

	// not found
	return []dm.SendMsg[dm.InKb]{
		{
			Text: dm.NoInvitation,
		},
	}, nil
}

func (s *service) RegistrationQRCallBack(ctx context.Context, telegramID, chatID int64, buttonID int32) ([]dm.SendMsg[dm.InKb], error) {
	exData, err := s.storage.DeleteCallbackBuffer(ctx, telegramID, chatID)
	// clicked the button again after already existing registration - nothing to delete
	if err != nil && err == pgx.ErrNoRows {
		log.Println("double click")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	switch exData.ServiceName {
	case dm.InlEntranceGuest:
		return s.asGuestEntr(ctx, telegramID, buttonID) // reg_id
	case dm.InlEntranceStudent:
		return s.asStudentEntr(ctx, telegramID, buttonID) // invite_id
	}

	return nil, fmt.Errorf("svc: RegistrationQRCallBack: unknown entrance type: %s", exData.ServiceName)
}
