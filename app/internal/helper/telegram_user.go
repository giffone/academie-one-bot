package helper

import (
	dm "ac_bot/internal/domain"

	"github.com/go-telegram/bot/models"
)

func TgUser(user *models.User) dm.TelegramUser {
	return dm.TelegramUser{
		ID:        user.ID,
		Login:     user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsBot:     user.IsBot,
	}
}
