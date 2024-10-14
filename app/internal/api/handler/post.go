package handler

import (
	"ac_bot/internal/helper"
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) PostStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	if helper.UpdateIsNil(update) {
		log.Println("hdl: PostStart: is nil")
		return
	}

	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID

	msg, err := h.svc.PostStart(ctx,telegramID)
	if err != nil {
		sendErrMessageClient(ctx, b, chatID, err)
		return
	}

	for _, m := range msg {
		if _, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      m.Text,
			ParseMode: models.ParseModeMarkdown,
		}); err != nil {
			log.Printf("[chatID: %d] hdl: PostStart: try to send message: %s", chatID, err.Error())
		}
	}
}

func (h *Handler) PostStop(ctx context.Context, b *bot.Bot, update *models.Update) {
	if helper.UpdateIsNil(update) {
		log.Println("hdl: PostStop: is nil")
		return
	}

	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID

	msg, err := h.svc.PostStop(ctx,telegramID)
	if err != nil {
		sendErrMessageClient(ctx, b, chatID, err)
		return
	}

	for _, m := range msg {
		if _, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      m.Text,
			ParseMode: models.ParseModeMarkdown,
		}); err != nil {
			log.Printf("[chatID: %d] hdl: PostStop: try to send message: %s", chatID, err.Error())
		}
	}
}
