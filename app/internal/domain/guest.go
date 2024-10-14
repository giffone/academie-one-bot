package domain

import "time"

type Guest struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	InviteCode     string `json:"invite_code"`
	OrganizationID string `json:"org_id"`
}

type GuestDTO struct {
	TelegramID int64
	FirstName  string
	LastName   string
	InviteID   int32
}

type GuestInvite struct {
	Code  string `json:"invite_code"`
	Title string `json:"invite_title"`
}

type GuestInviteDTO struct {
	ID     int32
	Code   string
	Title  string
	Expire time.Time
	Admin  AdminCardDTO
}

type GuestRegDTO struct {
	ID         int32
	TelegramID int64
	FirstName  string
	LastName   string
	Invite     struct {
		ID           int32
		Code         string
		Expire       time.Time
		Organization struct {
			ID    int32
			Title string
		}
		Admin struct {
			ID   int32
			Name string
		}
	}
}

type GuestEntranceCardDTO struct {
	RoleName          string
	FirstName         string
	LastName          string
	OrganizationTitle string
	AdminName         string
	InviteExpireDate  time.Time
}
