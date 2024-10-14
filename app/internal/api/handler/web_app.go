package handler

import (
	dm "ac_bot/internal/domain"
	"ac_bot/internal/helper"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) WebApp(ctx context.Context, b *bot.Bot, update *models.Update) {
	if helper.UpdateIsNil(update) {
		log.Println("hdl: WebApp: is nil")
		return
	}

	chatID := update.Message.Chat.ID

	var data dm.FormData

	err := json.Unmarshal([]byte(update.Message.WebAppData.Data), &data)
	if err != nil {
		sendErrMessageClient(ctx, b, chatID, fmt.Errorf("hdl: WebApp: unmarshal: %s", err.Error()))
		return
	}

	data.TelegramUser = helper.TgUser(update.Message.From)

	if data.TelegramUser.IsBot {
		err := dm.CliMsg{
			Message: dm.IsBot,
			Err:     fmt.Errorf("hdl: WebApp: is bot: %t", update.Message.From.IsBot),
		}
		sendErrMessageClient(ctx, b, chatID, &err)
		return
	}

	if data.TelegramUser.Login == "" {
		data.TelegramUser.Login = update.Message.From.FirstName
	}
	if data.TelegramUser.Login == "" {
		data.TelegramUser.Login = update.Message.From.LastName
	}

	var msg []dm.SendMsg[dm.InKb]

	switch data.FormType {
	case dm.FormTypeGuestRegForm:
		msg, err = h.svc.GuestRegForm(ctx, &data)
	case dm.FormTypeStudentRegForm:
		msg, err = h.svc.StudentRegForm(ctx, &data)
	case dm.FormTypeQR:
		msg, err = h.svc.RegistrationQR(ctx, chatID, &data)
	case dm.FormTypeCreateAdmin:
		msg, err = h.svc.RegisterAdmin(ctx, &data)
	case dm.FormTypeCreateInviteGuest:
		msg, err = h.svc.CreateInviteForGuest(ctx, &data)
	case dm.FormTypeCreateInviteStudent:
		msg, err = h.svc.CreateInviteForStudents(ctx, &data)
	default:
		sendErrMessageClient(ctx, b, chatID, fmt.Errorf("hdl: WebApp: unknown form type: %s", data.FormType))
		return
	}

	if err != nil {
		sendErrMessageClient(ctx, b, chatID, err)
		return
	}

	for _, m := range msg {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   m.Text,
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: helper.ParseInlineKb(m.Keyboard),
			},
		}); err != nil {
			log.Printf("[chatID: %d] hdl: WebApp: try to send message: %s", chatID, err.Error())
		}
	}
}
