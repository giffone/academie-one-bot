package handler

import (
	"ac_bot/internal/domain"

	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-telegram/bot"
)

func sendErrMessageClient(ctx context.Context, b *bot.Bot, chatID int64, err error) {
	// message for client
	msg := ""

	unwrappedErr := err

	for unwrappedErr != nil {
		// if have message for client
		if m, ok := unwrappedErr.(*domain.CliMsg); ok {
			msg = m.Message
			break
		}
		// next error
		unwrappedErr = errors.Unwrap(unwrappedErr)
	}

	if msg == "" {
		msg = fmt.Sprintf(domain.ErrorMsgClient, chatID, time.Now().Format(time.DateTime))
	}

	if err != nil {
		log.Printf("[chatID: %d] error is: %s", chatID, err.Error())
	}

	// message for client
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})
	if err != nil {
		log.Printf("[chatID: %d] try send err message: %s", chatID, err.Error())
	}
}
