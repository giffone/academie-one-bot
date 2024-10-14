package handler

import (
	"ac_bot/internal/helper"
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) About(ctx context.Context, b *bot.Bot, update *models.Update) {
	if helper.UpdateIsNil(update) {
		log.Println("hdl: About: is nil")
		return
	}

	chatID := update.Message.Chat.ID

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Here is more about of us:",
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					models.InlineKeyboardButton{
						Text:         "Info",
						URL:          "https://www.academie.one/",
					},
					models.InlineKeyboardButton{
						Text: "Intra",
						URL:  "https://zero.academie.one/",
					},
				},
				{
					models.InlineKeyboardButton{
						Text: "Instagram",
						URL:  "https://www.instagram.com/the.academie/",
					},
					models.InlineKeyboardButton{
						Text: "Facebook",
						URL:  "https://www.facebook.com/academie.one",
					},
				},
				{
					models.InlineKeyboardButton{
						Text: "Telegram",
						URL:  "https://t.me/theacademie",
					},
					models.InlineKeyboardButton{
						Text: "Youtube",
						URL:  "https://www.youtube.com/@academieone",
					},
				},
			},
		},
	}); err != nil {
		log.Printf("[chatID: %d] hdl: About: try to send message: %s", chatID, err.Error())
	}
}
