package postgres

import (
	"context"
	"fmt"
	"time"
)

type Notification interface {
	CreateEventListener(ctx context.Context, telegramID int64) error
	GetEventListeners(ctx context.Context) ([]int64, error)
	DeleteEventListener(ctx context.Context, telegramID int64) error
}

func (s *storage) CreateEventListener(ctx context.Context, telegramID int64) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx2,
		`INSERT INTO telegram.event_listeners (administrator_id)
VALUES ($1)
ON CONFLICT (administrator_id) DO NOTHING;`,
		telegramID,
	)
	if err != nil {
		return fmt.Errorf("db: CreateEventListener: exec: %w", err)
	}

	return nil
}

func (s *storage) GetEventListeners(ctx context.Context) ([]int64, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx2,
		`SELECT administrator_id
FROM telegram.event_listeners;`,
	)
	if err != nil {
		return nil, fmt.Errorf("db: GetEventListeners: query: %w", err)
	}

	defer rows.Close()

	lis := make([]int64, 0, 10)

	for rows.Next() {
		l := int64(0)

		err = rows.Scan(
			&l,
		)
		if err != nil {
			return nil, fmt.Errorf("db: GetEventListeners: scan: %w", err)
		}

		lis = append(lis, l)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db: GetOrgList: rows err: %w", err)
	}

	return lis, nil
}

func (s *storage) DeleteEventListener(ctx context.Context, telegramID int64) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx2,
		`DELETE FROM telegram.event_listeners
WHERE administrator_id = $1;`,
		telegramID,
	)
	if err != nil {
		return fmt.Errorf("db: DeleteEventListener: exec: %w", err)
	}

	return nil
}
