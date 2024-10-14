package domain

type TelegramUser struct {
	ID        int64
	Login     string
	FirstName string
	LastName  string
	IsBot     bool
}
