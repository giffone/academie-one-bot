package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) Callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.CallbackQuery == nil {
		log.Println("hdl: Callback: update or callback is nil!!!")
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	telegramID := update.CallbackQuery.From.ID // user_id

	// is the id of the bot from which the message came(not the user_id). id must be 7470860586 The Academie One  academie_one_bot
	if update.CallbackQuery.Message.Message.From.ID != h.botID {
		log.Printf("hdl: Callback: bot_id [%d] and incoming bot_id [%d] not equal", h.botID, update.CallbackQuery.Message.Message.From.ID)
		return
	}

	// answering callback query first to let Telegram know that we received the callback query,
	// and we're handling it.
	if _, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	}); err != nil {
		log.Printf("[chatID: %d] hdl: AnswerCallbackQuery: try to send message: %s", chatID, err.Error())
		return
	}

	// remove word "button=5" and get only "5"
	s := strings.Split(update.CallbackQuery.Data, "=")
	buttonIdStr := ""
	if len(s) != 2 {
		sendErrMessageClient(ctx, b, chatID, fmt.Errorf("hdl: Callback: button data is incorrect: %s", update.CallbackQuery.Data))
		return
	} else {
		buttonIdStr = s[1]
	}

	buttonId, err := strconv.ParseInt(buttonIdStr, 10, 64)
	if err != nil {
		sendErrMessageClient(ctx, b, chatID, fmt.Errorf("hdl: Callback: parseInt: id: %s", err.Error()))
		return
	}

	msg, err := h.svc.RegistrationQRCallBack(ctx, telegramID, chatID, int32(buttonId))
	if err != nil {
		sendErrMessageClient(ctx, b, chatID, fmt.Errorf("hdl: Callback: %w", err))
		return
	}

	for _, m := range msg {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   m.Text,
			// ReplyMarkup: models.InlineKeyboardMarkup{
			// 	InlineKeyboard: helper.ParseInlineKb(m.Keyboard),
			// },
		}); err != nil {
			log.Printf("[chatID: %d] hdl: Callback: try to send message: %s", chatID, err.Error())
		}
	}
}
