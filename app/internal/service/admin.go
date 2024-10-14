package service

import (
	dm "ac_bot/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Admin interface {
	RegisterAdmin(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error)
	PostStart(ctx context.Context, telegramID int64) ([]dm.SendMsg[dm.InKb], error)
	PostStop(ctx context.Context, telegramID int64) ([]dm.SendMsg[dm.InKb], error)
}

func (s *service) RegisterAdmin(ctx context.Context, data *dm.FormData) ([]dm.SendMsg[dm.InKb], error) {
	// verify that the request is from an administrator
	if _, err := s.isAdmin(ctx, data.TelegramUser.ID); err != nil {
		return nil, fmt.Errorf("svc: RegisterAdmin: %w", err) // wrap!
	}

	// prepare data for db
	var newA dm.Admin

	if err := json.Unmarshal(data.UserData, &newA); err != nil {
		return nil, fmt.Errorf("svc: RegisterAdmin: unmarshal: %s", err.Error())
	}

	tId, err := strconv.ParseInt(newA.TelegramID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("svc: RegisterAdmin: parseInt: telegram id: %s", err.Error())
	}

	rId, err := strconv.ParseInt(newA.RoleID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("svc: RegisterAdmin: parseInt: role id: %s", err.Error())
	}

	mId, err := strconv.ParseInt(newA.OrganizationID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("svc: RegisterAdmin: parseInt: org id: %s", err.Error())
	}

	dto := dm.AdminDTO{
		TelegramID:     tId,
		Name:           newA.Name,
		RoleID:         int32(rId),
		OrganizationID: int32(mId),
	}

	// check for legal data
	if dto.TelegramID == 0 || dto.RoleID == 0 {
		return nil, fmt.Errorf("svc: RegisterAdmin: telegram_id or role_id equal 0")
	}

	// db
	if err = s.storage.RegisterAdmin(ctx, dto); err != nil {
		return nil, err
	}

	return []dm.SendMsg[dm.InKb]{
		{
			Text: dm.Created,
		},
	}, nil
}

func (s *service) isAdmin(ctx context.Context, telegramID int64) (*dm.AdminCardDTO, error) {
	if telegramID <= 0 {
		return nil, errors.New("svc: isAdmin: ID is less or equal 0")
	}

	// db
	card, err := s.storage.GetAdministratorCard(ctx, telegramID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	// check id
	if card == nil || card.RoleID != dm.AdminRole {
		return nil, &dm.CliMsg{
			Message: dm.AdminRightNeed,
			Err:     fmt.Errorf("svc: isAdmin: user: %d", telegramID),
		}
	}

	return card, nil
}

func (s *service) PostStart(ctx context.Context, telegramID int64) ([]dm.SendMsg[dm.InKb], error) {
	_, err := s.storage.IsAdmin(ctx, telegramID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err = s.storage.CreateEventListener(ctx, telegramID); err != nil {
		return nil, err
	}

	return []dm.SendMsg[dm.InKb]{
		{
			Text: dm.Started,
		},
	}, nil
}

func (s *service) PostStop(ctx context.Context, telegramID int64) ([]dm.SendMsg[dm.InKb], error) {
	_, err := s.storage.IsAdmin(ctx, telegramID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err = s.storage.DeleteEventListener(ctx, telegramID); err != nil {
		return nil, err
	}

	return []dm.SendMsg[dm.InKb]{
		{
			Text: dm.Stopped,
		},
	}, nil
}
