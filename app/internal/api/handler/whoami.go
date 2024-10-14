package handler

import (
	"ac_bot/internal/helper"
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) Whoami(ctx context.Context, b *bot.Bot, update *models.Update) {
	if helper.UpdateIsNil(update) {
		log.Println("hdl: Whoami: is nil")
		return
	}

	msg := h.svc.Whoami(ctx, helper.TgUser(update.Message.From))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg.Text,
		ParseMode: models.ParseModeMarkdown,
	}); err != nil {
		log.Printf("[chatID: %d] hdl: Whoami: try to send message: %s", update.Message.Chat.ID, err.Error())
	}
}
