package domain

type NewNotification struct {
	RoleID     int32 // student or guest
	CauseID    int32 // regID for guest or inviteID for student
	TelegramID int64
}
