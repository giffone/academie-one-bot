package postgres

import (
	dm "ac_bot/internal/domain"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Guest interface {
	RegisterGuest(ctx context.Context, dto dm.GuestDTO) error
	GetGuestRegistrations(ctx context.Context, telegramID int64) ([]dm.GuestRegDTO, error)

	CreateInviteForGuest(ctx context.Context, invite dm.GuestInvite, adminID int64) (time.Time, error)
	GetGuestInviteCard(ctx context.Context, code string, organizationID int32) (*dm.GuestInviteDTO, error)

	CreateEntranceGuest(ctx context.Context, regID int32) error
	GetGuestEntranceCard(ctx context.Context, telegramID int64, regID int32) (*dm.GuestEntranceCardDTO, error)

	// GetEntranceGuestExpireDate(ctx context.Context, regID int32) (time.Time, error)
}

func (s *storage) RegisterGuest(ctx context.Context, dto dm.GuestDTO) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx2,
		`INSERT INTO telegram.regform_guest (telegram_id, first_name, last_name, invite_id)
VALUES ($1, $2, $3, $4);`,
		dto.TelegramID,
		dto.FirstName,
		dto.LastName,
		dto.InviteID,
	)
	if err != nil {
		return fmt.Errorf("db: RegisterGuest: exec: %w", err)
	}

	return nil
}

func (s *storage) GetGuestRegistrations(ctx context.Context, telegramID int64) ([]dm.GuestRegDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx2,
		`SELECT r.id, r.telegram_id, r.first_name, r.last_name, i.id, i.code, i.expire, o.id, o.title, a.telegram_id, a.name
FROM telegram.regform_guest r
JOIN telegram.invite_guest i ON r.invite_id = i.id
JOIN telegram.administrator a ON i.administrator_id = a.telegram_id
JOIN telegram.organization o ON a.organization_id = o.id
WHERE r.telegram_id = $1
AND CURRENT_TIMESTAMP > i.created
AND CURRENT_DATE < i.expire;`,
		telegramID,
	)
	if err != nil {
		return nil, fmt.Errorf("db: GetGuestRegistrations: query: %w", err)
	}

	defer rows.Close()

	regs := make([]dm.GuestRegDTO, 0, 5)

	for rows.Next() {
		reg := dm.GuestRegDTO{}
		err = rows.Scan(
			&reg.ID,
			&reg.TelegramID,
			&reg.FirstName,
			&reg.LastName,
			&reg.Invite.ID,
			&reg.Invite.Code,
			&reg.Invite.Expire,
			&reg.Invite.Organization.ID,
			&reg.Invite.Organization.Title,
			&reg.Invite.Admin.ID,
			&reg.Invite.Admin.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("db: GetGuestRegistrations: scan: %w", err)
		}

		regs = append(regs, reg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db: GetGuestRegistrations: rows err: %w", err)
	}

	return regs, nil
}

func (s *storage) CreateInviteForGuest(ctx context.Context, invite dm.GuestInvite, adminID int64) (time.Time, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var expire time.Time

	row := s.pool.QueryRow(ctx2,
		`INSERT INTO telegram.invite_guest (code, title, administrator_id, created, expire)
VALUES ($1, $2, $3, DEFAULT, DEFAULT)
RETURNING expire;`,
		strings.ToLower(invite.Code),
		invite.Title,
		adminID,
	)

	if err := row.Scan(&expire); err != nil {
		if err == pgx.ErrNoRows {
			return time.Unix(0, 0), err
		}

		return time.Unix(0, 0), fmt.Errorf("db: CreateInviteForGuest: scan: %w", err)
	}

	return expire, nil
}

func (s *storage) GetGuestInviteCard(ctx context.Context, code string, organizationID int32) (*dm.GuestInviteDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var invite dm.GuestInviteDTO

	row := s.pool.QueryRow(ctx2,
		`SELECT i.id, i.code, i.title, i.expire, a.telegram_id, a.name, r.id ,r."name" ,o.id ,o.title 
FROM telegram.invite_guest i
JOIN telegram.administrator a ON i.administrator_id = a.telegram_id
JOIN telegram.organization o ON a.organization_id = o.id
JOIN telegram.role r ON a.role_id = r.id
WHERE i.code = $1
AND o.id = $2
AND i.expire > CURRENT_DATE;`,
		strings.ToLower(code),
		organizationID,
	)

	if err := row.Scan(
		&invite.ID,
		&invite.Code,
		&invite.Title,
		&invite.Expire,
		&invite.Admin.TelegramID,
		&invite.Admin.Name,
		&invite.Admin.RoleID,
		&invite.Admin.RoleName,
		&invite.Admin.OrganizationID,
		&invite.Admin.OrganizationTitle,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("db: GetGuestInviteID: scan: %w", err)
	}

	return &invite, nil
}

func (s *storage) CreateEntranceGuest(ctx context.Context, regID int32) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx2,
		`INSERT INTO telegram.entrance_guest (regform_id)
VALUES ($1);`,
		regID,
	)
	if err != nil {
		return fmt.Errorf("db: CreateEntranceGuest: scan: %w", err)
	}

	return nil
}

// func (s *storage) GetEntranceGuestExpireDate(ctx context.Context, regID int32) (time.Time, error) {
// 	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	var expire time.Time

// 	row := s.pool.QueryRow(ctx2,
// 		`SELECT i.expire
// FROM entrance_guest e
// JOIN regform_guest r on e.regform_id = r.id
// JOIN invite_guest i on r.invite_id = i.id
// WHERE e.regform_id  = $1;`,
// 		regID,
// 	)

// 	if err := row.Scan(&expire); err != nil {
// 		if err == pgx.ErrNoRows {
// 			return time.Unix(0, 0), err
// 		}

// 		return time.Unix(0, 0), fmt.Errorf("db: GetEntranceGuestExpireDate: scan: %w", err)
// 	}

// 	return expire, nil
// }

func (s *storage) GetGuestEntranceCard(ctx context.Context, telegramID int64, regID int32) (*dm.GuestEntranceCardDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var entr dm.GuestEntranceCardDTO

	row := s.pool.QueryRow(ctx2,
		`SELECT rl.name, r.first_name, r.last_name, o.title, a.name, i.expire
FROM telegram.regform_guest r
JOIN telegram.invite_guest i on r.invite_id = i.id
JOIN telegram.administrator a on i.administrator_id = a.telegram_id
JOIN telegram.organization o on a.organization_id = o.id
join telegram.role rl ON r.role_id = rl.id
WHERE r.id = $1
AND r.telegram_id = $2;`,
		regID,
		telegramID,
	)

	if err := row.Scan(
		&entr.RoleName,
		&entr.FirstName,
		&entr.LastName,
		&entr.OrganizationTitle,
		&entr.AdminName,
		&entr.InviteExpireDate,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("db: GetGuestEntranceCard: scan: %w", err)
	}

	log.Println("GuestEntranceCard", entr)

	return &entr, nil
}
