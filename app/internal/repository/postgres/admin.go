package postgres

import (
	dm "ac_bot/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Administrator interface {
	RegisterAdmin(ctx context.Context, dto dm.AdminDTO) error
	GetAdministratorCard(ctx context.Context, adminID int64) (*dm.AdminCardDTO, error)

	IsAdmin(ctx context.Context, telegramID int64) (*dm.AdminDTO, error)
}

func (s *storage) RegisterAdmin(ctx context.Context, dto dm.AdminDTO) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx2,
		`INSERT INTO telegram.administrator (telegram_id, name, role_id, organization_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (telegram_id)
DO UPDATE SET
name = EXCLUDED.name,
role_id = EXCLUDED.role_id;`,
		dto.TelegramID,
		dto.Name,
		dto.RoleID,
		dto.OrganizationID,
	)
	if err != nil {
		return fmt.Errorf("db: RegisterAdmin: exec: %w", err)
	}

	return nil
}

func (s *storage) GetAdministratorCard(ctx context.Context, adminID int64) (*dm.AdminCardDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var card dm.AdminCardDTO

	row := s.pool.QueryRow(ctx2,
		`SELECT a.telegram_id, a.name, a.role_id, r.name, a.organization_id, o.title
FROM telegram.administrator a
JOIN telegram.role r ON a.role_id = r.id
join telegram.organization o ON a.organization_id = o.id
WHERE a.telegram_id = $1;`,
		adminID,
	)

	err := row.Scan(
		&card.TelegramID,
		&card.Name,
		&card.RoleID,
		&card.RoleName,
		&card.OrganizationID,
		&card.OrganizationTitle,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("db: GetAdministratorCard: scan: %w", err)
	}

	return &card, nil
}

func (s *storage) IsAdmin(ctx context.Context, telegramID int64) (*dm.AdminDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var dto dm.AdminDTO

	row := s.pool.QueryRow(ctx2,
		`SELECT telegram_id, name, role_id, organization_id
FROM telegram.administrator
WHERE telegram_id = $1;`,
		telegramID,
	)

	if err := row.Scan(
		&dto.TelegramID,
		&dto.Name,
		&dto.RoleID,
		&dto.OrganizationID,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("db: IsAdmin: scan: %w", err)
	}

	return &dto, nil
}
