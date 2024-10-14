package tgs

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramService interface {
	SendMessage(ctx context.Context, chatID int64, message string)
}

func New(b *bot.Bot) TelegramService {
	return &telegramService{
		b: b,
	}
}

type telegramService struct {
	b *bot.Bot
}

func (ts *telegramService) SendMessage(ctx context.Context, chatID int64, message string) {
	if _, err := ts.b.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
		ParseMode: models.ParseModeMarkdown,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					models.InlineKeyboardButton{
						Text:         "Ok",
						CallbackData: "Ok",
					},
				},
			},
		},
	}); err != nil {
		log.Printf("[chatID: %d] tsvc: try to send message: %s", chatID, err.Error())
	}
}
