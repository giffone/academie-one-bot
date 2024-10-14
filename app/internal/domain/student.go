package domain

import "time"

type Student struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Login          string `json:"login"`
	InviteCode     string `json:"invite_code"`
	OrganizationID string `json:"org_id"`
}

type StudentDTO struct {
	TelegramID int64
	Login      string
	FirstName  string
	LastName   string
	InviteID   int32
}

type StudentInvite struct {
	Code       string `json:"invite_code"`
	ExpireDate string `json:"expire_date"`
}

type StudentRegDTO struct {
	TelegramID int64
	Login      string
	FirstName  string
	LastName   string
	Approved   bool
	Invite     struct {
		ID           int32
		Code         string
		Expire       time.Time
		Organization struct {
			ID    int32
			Title string
		}
	}
}

type StudentEntranceCardDTO struct {
	RoleName             string // student
	InviteCode           string // piscine23sep
	FirstName            string
	LastName             string
	OrganizationTitle    string
	RegistrationApproved bool
	InviteExpireDate     time.Time
}
