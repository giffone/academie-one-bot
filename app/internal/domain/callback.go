package domain

type CallbackDataDTO struct {
	TelegramID  int64
	ChatID      int64
	ServiceName string // varchar(10)
	RoleID      int32
}

type CallbackButton struct {
	Title string
	ObjID int32
}
