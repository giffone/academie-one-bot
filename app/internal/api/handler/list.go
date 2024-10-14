package handler

import (
	"ac_bot/internal/helper"
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) GetList(ctx context.Context, b *bot.Bot, update *models.Update) {
	if helper.UpdateIsNil(update) {
		log.Println("hdl: GetList: is nil")
		return
	}

	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID
	path := update.Message.Text

	list, err := h.svc.GetList(ctx, telegramID, path)
	if err != nil {
		sendErrMessageClient(ctx, b, chatID, fmt.Errorf("hdl: GetList: %w", err))
		return
	}

	if _, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      list,
		ParseMode: models.ParseModeMarkdown,
	}); err != nil {
		log.Printf("[chatID: %d] hdl: GetList: try to send message: %s", chatID, err.Error())
	}
}
