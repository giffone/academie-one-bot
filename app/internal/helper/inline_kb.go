package helper

import (
	"ac_bot/internal/domain"

	"github.com/go-telegram/bot/models"
)

func ParseInlineKb(keyboard domain.InKb) [][]models.InlineKeyboardButton {
	tKb := make([][]models.InlineKeyboardButton, len(keyboard))

	for i, kb := range keyboard {
		tKb[i] = make([]models.InlineKeyboardButton, len(kb))
		for j, k := range kb {
			tKb[i][j] = models.InlineKeyboardButton{
				Text:         k.Text,
				URL:          k.Url,
				CallbackData: k.CallbackData,
			}
		}
	}

	return tKb
}
