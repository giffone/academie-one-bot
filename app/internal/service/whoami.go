package service

import (
	"ac_bot/internal/domain"
	"bytes"
	"context"
	"strconv"

	"github.com/go-telegram/bot"
)

func (s *service) Whoami(ctx context.Context, user domain.TelegramUser) *domain.SendMsg[any] {
	buf := bytes.Buffer{}

	buf.WriteString("ID: ")
	buf.WriteString(strconv.FormatInt(user.ID, 10))
	buf.WriteRune('\n')
	buf.WriteString("First name: ")
	buf.WriteString(bot.EscapeMarkdown(user.FirstName))
	buf.WriteRune('\n')
	buf.WriteString("Last name: ")
	buf.WriteString(bot.EscapeMarkdown(user.LastName))
	buf.WriteRune('\n')
	buf.WriteString("User name: ")
	buf.WriteString(bot.EscapeMarkdown(user.Login))

	return &domain.SendMsg[any]{
		Text: buf.String(),
	}
}
