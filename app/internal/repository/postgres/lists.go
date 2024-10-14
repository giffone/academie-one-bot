package postgres

import (
	dm "ac_bot/internal/domain"
	"context"
	"fmt"
	"time"
)

type Lists interface {
	GetOrgList(ctx context.Context) ([]dm.OrganizationDTO, error)

	GetRoleList(ctx context.Context, adminOnly bool) ([]dm.Role, error)
	// GetRole(ctx context.Context, telegramID int64) (*domain.Role, error)
}

func (s *storage) GetOrgList(ctx context.Context) ([]dm.OrganizationDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx2,
		`SELECT id, name, title
FROM telegram.organization;`,
	)
	if err != nil {
		return nil, fmt.Errorf("db: GetOrgList: query: %w", err)
	}

	defer rows.Close()

	orgs := make([]dm.OrganizationDTO, 0, 10)

	for rows.Next() {
		org := dm.OrganizationDTO{}

		err = rows.Scan(
			&org.ID,
			&org.Name,
			&org.Title,
		)
		if err != nil {
			return nil, fmt.Errorf("db: GetOrgList: scan: %w", err)
		}

		orgs = append(orgs, org)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db: GetOrgList: rows err: %w", err)
	}

	return orgs, nil
}

func (s *storage) GetRoleList(ctx context.Context, adminOnly bool) ([]dm.Role, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, name
FROM telegram.role;`

	if adminOnly {
		query = `SELECT id, name
FROM telegram.role
WHERE admin = true;`
	}

	rows, err := s.pool.Query(ctx2, query)
	if err != nil {
		return nil, fmt.Errorf("db: GetRoleList: query: %w", err)
	}

	defer rows.Close()

	roles := make([]dm.Role, 0, 10)

	for rows.Next() {
		role := dm.Role{}

		err = rows.Scan(
			&role.ID,
			&role.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("db: GetRoleList: scan: %w", err)
		}

		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db: GetRoleList: rows err: %w", err)
	}

	return roles, nil
}

// func (s *storage) GetRole(ctx context.Context, telegramID int64) (*domain.Role, error) {
// 	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	var role domain.Role

// 	row := s.pool.QueryRow(ctx2,
// 		`SELECT r.id, r.name
// FROM telegram.role r
// JOIN telegram.administrator a ON r.id = a.role_id
// WHERE a.telegram_id = $1;`,
// 		telegramID,
// 	)

// 	if err := row.Scan(&role.ID, &role.Name); err != nil {
// 		if err == pgx.ErrNoRows {
// 			return nil, err
// 		}

// 		return nil, fmt.Errorf("db: GetRole: scan: %w", err)
// 	}

// 	return &role, nil
// }
