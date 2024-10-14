package domain

import "encoding/json"

type FormData struct {
	FormType     string          `json:"form_type"`
	UserData     json.RawMessage `json:"user_data"`
	TelegramUser TelegramUser    `json:"-"`
}
