package handler

import (
	"ac_bot/internal/helper"
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	if helper.UpdateIsNil(update) {
		log.Println("hdl: Default: is nil")
		return
	}

	chatID := update.Message.Chat.ID

	msg, err := h.svc.Start(ctx, update.Message.From.ID)
	if err != nil {
		sendErrMessageClient(ctx, b, chatID, fmt.Errorf("hdl: Default: Start: %w", err))
		return
	}

	if _, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg.Text,
		ReplyMarkup: models.ReplyKeyboardMarkup{
			ResizeKeyboard: true,
			IsPersistent:   true, // always show
			Keyboard:       helper.ParseBottomKb(msg.Keyboard),
		},
	}); err != nil {
		log.Printf("[chatID: %d] hdl: Default: try to send message: %s", chatID, err.Error())
	}
}
