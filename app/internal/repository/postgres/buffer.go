package postgres

import (
	dm "ac_bot/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type CallbackBuffer interface {
	CreateCallbackBuffer(ctx context.Context, data dm.CallbackDataDTO) error
	DeleteCallbackBuffer(ctx context.Context, telegramID int64, chatID int64) (*dm.CallbackDataDTO, error)
}

func (s *storage) CreateCallbackBuffer(ctx context.Context, data dm.CallbackDataDTO) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx2,
		`INSERT INTO telegram.callback_buffer (telegram_id, chat_id, svc_name, role_id)
VALUES ($1, $2, $3, $4);`,
		data.TelegramID,
		data.ChatID,
		data.ServiceName,
		data.RoleID,
	)
	if err != nil {
		return fmt.Errorf("db: CreateCallbackBuffer: exec: %w", err)
	}

	return nil
}

func (s *storage) DeleteCallbackBuffer(ctx context.Context, telegramID int64, chatID int64) (*dm.CallbackDataDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var dto dm.CallbackDataDTO

	row := s.pool.QueryRow(ctx2,
		`DELETE FROM telegram.callback_buffer
WHERE telegram_id = $1
AND chat_id = $2
RETURNING telegram_id, chat_id, svc_name, role_id;`,
		telegramID,
		chatID,
	)
	if err := row.Scan(
		&dto.TelegramID,
		&dto.ChatID,
		&dto.ServiceName,
		&dto.RoleID,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, err // nothing to delete
		}

		return nil, fmt.Errorf("db: DeleteCallbackBuffer: scan: %w", err)
	}

	return &dto, nil
}
