package service

import (
	dm "ac_bot/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Student interface {
	CreateInviteForStudents(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error)
	StudentRegForm(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error)
}

func (s *service) CreateInviteForStudents(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error) {
	var invite dm.StudentInvite

	err := json.Unmarshal(data.UserData, &invite)
	if err != nil {
		return nil, fmt.Errorf("svc: CreateInviteForStudents: unmarshal: %s", err.Error())
	}

	// check for legal data
	if invite.Code == "" {
		return nil, &dm.CliMsg{
			Message: "Invite code is empty.",
			Err:     fmt.Errorf("svc: CreateInviteForStudents: invite code is empty from user: %d", data.TelegramUser.ID),
		}
	}

	if _, err := time.Parse("2006-01-02", invite.ExpireDate); err != nil {
		return nil, &dm.CliMsg{
			Message: "Can not parse Date format.",
			Err:     fmt.Errorf("svc: CreateInviteForStudents: can not parse Date format from user: %d", data.TelegramUser.ID),
		}
	}

	// verify that the request is from an administrator and get his card
	card, err := s.isAdmin(ctx, data.TelegramUser.ID)
	if err != nil {
		return nil, err
	}

	// use org_id (student_invite belongs to an organization)
	if err = s.storage.CreateInviteForStudents(ctx, invite, card.OrganizationID); err != nil {
		return nil, err
	}

	return []dm.SendMsg[dm.InKb]{
		{
			Text: fmt.Sprintf(dm.InviteForStudentAdmin, invite.Code, invite.ExpireDate),
		},
	}, nil
}

func (s *service) StudentRegForm(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error) {
	var userData dm.Student

	err := json.Unmarshal(data.UserData, &userData)
	if err != nil {
		return nil, fmt.Errorf("svc: StudentRegForm: unmarshal: %s", err.Error())
	}

	// org_id {string -> int}
	oId, err := strconv.ParseInt(userData.OrganizationID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("svc: StudentRegForm: parseInt: organization id: %s", err.Error())
	}

	// check invite
	iId, err := s.storage.GetStudentInviteID(ctx, userData.InviteCode, int32(oId))
	if err != nil {
		// not found
		if err == pgx.ErrNoRows {
			return []dm.SendMsg[dm.InKb]{
				{
					Text: dm.InviteIDWrong,
				},
			}, nil
		}

		return nil, err
	}

	dto := dm.StudentDTO{
		TelegramID: data.TelegramUser.ID,
		Login:      strings.ToLower(userData.Login),
		FirstName:  userData.FirstName,
		LastName:   userData.LastName,
		InviteID:   iId,
	}

	// write to DB
	if err = s.storage.RegisterStudent(ctx, dto); err != nil {
		return nil, err
	}

	return []dm.SendMsg[dm.InKb]{
		{
			Text: dm.Created,
		},
	}, nil
}
