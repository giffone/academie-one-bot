package service

import (
	dm "ac_bot/internal/domain"
	"ac_bot/internal/helper"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *service) eventListener(notification <-chan dm.NewNotification) {
	log.Println("svc: event listener started!")

	for {
		select {
		case n := <-notification:
			go s.readNotification(s.ctx, n)
		case <-s.ctx.Done():
			log.Println("svc: eventListener: stopping due to application context cancellation")
			return
		}
	}
}

func (s *service) createNotification(telegramID int64, roleID, causeID int32) {
	log.Println("in createNotification")

	ctx2, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	select {
	case s.event.ch <- dm.NewNotification{
		RoleID:     roleID,
		CauseID:    causeID,
		TelegramID: telegramID,
	}:
		log.Println("svc: createNotification: sended to chanel")
	case <-ctx2.Done():
		log.Println("svc: createNotification: context timed out before sending to chanel")
	}
}

func (s *service) readNotification(ctx context.Context, data dm.NewNotification) {
	log.Println("in readNotification")
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// create student info card
	if data.RoleID == dm.StudentRole {
		dto, err := s.storage.GetStudentEntranceCard(ctx2, data.TelegramID, data.CauseID)
		if err != nil {
			if err != pgx.ErrNoRows {
				log.Printf("svc: readNotification: %s", err)
			}
			s.sendNotification(ctx, "")
			return
		}

		s.sendNotification(ctx, helper.StudentEntrCardMkdw(dto))
		return
	}

	// create guest info card
	if data.RoleID == dm.GuestRole {
		dto, err := s.storage.GetGuestEntranceCard(ctx2, data.TelegramID, data.CauseID)
		if err != nil {
			if err != pgx.ErrNoRows {
				log.Printf("svc: readNotification: %s", err)
			}
			s.sendNotification(ctx, "")
			return
		}

		s.sendNotification(ctx, helper.GuestEntrCardMkdw(dto))
		return
	}

	// unknown
	s.sendNotification(ctx, "")
}

func (s *service) sendNotification(ctx context.Context, msg string) {
	log.Println("in sendNotification")

	if msg == "" {
		msg = helper.UnknownEntrCardMkdw()
	}

	// find listeners
	lis, err := s.storage.GetEventListeners(ctx)
	if err != nil {
		log.Printf("svc: sendNotification: %s", err)
		return
	}

	if len(lis) == 0 {
		log.Println("svc: sendNotification: GetEventListeners: no listeners")
		return
	}

	// send notification
	for _, chatID := range lis {
		log.Printf("sending lis [%d] msg: %s\n", chatID, msg)
		s.event.tSvc.SendMessage(ctx, chatID, msg)
	}

}
