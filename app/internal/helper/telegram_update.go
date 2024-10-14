package helper

import (
	"github.com/go-telegram/bot/models"
)

func UpdateIsNil(update *models.Update) bool {
	if update == nil || update.Message == nil || update.Message.From == nil {
		return true
	}
	return false
}
