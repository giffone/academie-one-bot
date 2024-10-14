package domain

type Admin struct {
	TelegramID     string `json:"telegram_id"`
	Name           string `json:"user_name"`
	RoleID         string `json:"role_id"`
	OrganizationID string `json:"org_id"`
}

type AdminDTO struct {
	TelegramID     int64
	Name           string
	RoleID         int32
	OrganizationID int32
}

type AdminCardDTO struct {
	TelegramID        int64
	Name              string
	RoleID            int32
	RoleName          string
	OrganizationID    int32
	OrganizationTitle string
}
