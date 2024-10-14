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

type Student interface {
	RegisterStudent(ctx context.Context, dto dm.StudentDTO) error
	GetStudentRegistrations(ctx context.Context, telegramID int64) ([]dm.StudentRegDTO, error)

	CreateInviteForStudents(ctx context.Context, invite dm.StudentInvite, memberID int32) error
	GetStudentInviteID(ctx context.Context, code string, orgID int32) (int32, error)

	CreateEntranceStudent(ctx context.Context, telegramID int64, inviteID int32) (time.Time, error)
	GetStudentEntranceCard(ctx context.Context, telegramID int64, inviteID int32) (*dm.StudentEntranceCardDTO, error)

	IsStudent(ctx context.Context, telegramID int64) (*dm.StudentDTO, error)
}

func (s *storage) RegisterStudent(ctx context.Context, dto dm.StudentDTO) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log.Println("StudentDTO is", dto)

	_, err := s.pool.Exec(ctx2,
		`INSERT INTO telegram.regform_student (telegram_id, login, first_name, last_name, invite_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (telegram_id, invite_id)
DO UPDATE SET
login = EXCLUDED.login,
first_name = EXCLUDED.first_name,
last_name = EXCLUDED.last_name;`,
		dto.TelegramID,
		dto.Login,
		dto.FirstName,
		dto.LastName,
		dto.InviteID,
	)
	if err != nil {
		return fmt.Errorf("db: RegisterStudent: exec: %w", err)
	}

	return nil
}

func (s *storage) GetStudentRegistrations(ctx context.Context, telegramID int64) ([]dm.StudentRegDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx2,
		`SELECT r.telegram_id, r.login, r.first_name, r.last_name, r.approved, i.id, i.code, i.expire, o.id, o.title
FROM telegram.regform_student r
JOIN telegram.invite_student i ON r.invite_id = i.id
JOIN telegram.organization o ON i.organization_id = o.id
WHERE r.telegram_id = $1
AND CURRENT_TIMESTAMP > i.created
AND CURRENT_DATE < i.expire;`,
		telegramID,
	)
	if err != nil {
		return nil, fmt.Errorf("db: GetStudentRegistrations: query: %w", err)
	}

	defer rows.Close()

	regs := make([]dm.StudentRegDTO, 0, 5)

	for rows.Next() {
		reg := dm.StudentRegDTO{}

		err = rows.Scan(
			&reg.TelegramID,
			&reg.Login,
			&reg.FirstName,
			&reg.LastName,
			&reg.Approved,
			&reg.Invite.ID,
			&reg.Invite.Code,
			&reg.Invite.Expire,
			&reg.Invite.Organization.ID,
			&reg.Invite.Organization.Title,
		)
		if err != nil {
			return nil, fmt.Errorf("db: GetStudentRegistrations: scan: %w", err)
		}

		regs = append(regs, reg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db: GetStudentRegistrations: rows err: %w", err)
	}

	return regs, nil
}

func (s *storage) CreateInviteForStudents(ctx context.Context, invite dm.StudentInvite, memberID int32) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx2,
		`INSERT INTO telegram.invite_student (code, organization_id, expire)
VALUES ($1, $2, $3);`,
		strings.ToLower(invite.Code),
		memberID,
		invite.ExpireDate,
	)
	if err != nil {
		return fmt.Errorf("db: CreateInviteStudent: exec: %w", err)
	}

	return nil
}

func (s *storage) GetStudentInviteID(ctx context.Context, code string, orgID int32) (int32, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var invite int32

	row := s.pool.QueryRow(ctx2,
		`SELECT id
FROM telegram.invite_student
WHERE LOWER(code) = $1
AND organization_id = $2;`,
		strings.ToLower(code),
		orgID,
	)

	if err := row.Scan(&invite); err != nil {
		if err == pgx.ErrNoRows {
			return 0, err
		}

		return 0, fmt.Errorf("db: GetStudentInviteID: scan: %w", err)
	}

	return invite, nil
}

func (s *storage) CreateEntranceStudent(ctx context.Context, telegramID int64, inviteID int32) (time.Time, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var expire time.Time

	row := s.pool.QueryRow(ctx2,
		`INSERT INTO telegram.entrance_student (telegram_id, invite_id, created, expire)
VALUES ($1, $2, DEFAULT, DEFAULT)
RETURNING expire;`,
		telegramID,
		inviteID,
	)

	if err := row.Scan(&expire); err != nil {
		if err == pgx.ErrNoRows {
			return time.Unix(0, 0), err
		}

		return time.Unix(0, 0), fmt.Errorf("db: CreateEntranceStudent: scan: %w", err)
	}

	return expire, nil
}

func (s *storage) GetStudentEntranceCard(ctx context.Context, telegramID int64, inviteID int32) (*dm.StudentEntranceCardDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var entr dm.StudentEntranceCardDTO

	row := s.pool.QueryRow(ctx2,
		`SELECT rl.name, i.code, r.first_name, r.last_name, o.title, r.approved, i.expire
FROM telegram.regform_student r
JOIN telegram.invite_student i ON r.invite_id = i.id
JOIN telegram.organization o ON i.organization_id = o.id
join telegram.role rl ON r.role_id = rl.id 
WHERE r.invite_id = $1
AND r.telegram_id = $2;`,
		inviteID,
		telegramID,
	)

	if err := row.Scan(
		&entr.RoleName,
		&entr.InviteCode,
		&entr.FirstName,
		&entr.LastName,
		&entr.OrganizationTitle,
		&entr.RegistrationApproved,
		&entr.InviteExpireDate,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("db: GetStudentEntranceCard: scan: %w", err)
	}

	log.Println("StudentEntranceCard", entr)

	return &entr, nil
}

func (s *storage) IsStudent(ctx context.Context, telegramID int64) (*dm.StudentDTO, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var dto dm.StudentDTO

	row := s.pool.QueryRow(ctx2,
		`SELECT telegram_id, login, first_name, last_name, invite_id
FROM telegram.regform_student
WHERE approved = true
AND telegram_id = $1;`,
		telegramID,
	)

	if err := row.Scan(
		&dto.TelegramID,
		&dto.Login,
		&dto.FirstName,
		&dto.LastName,
		&dto.InviteID,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("db: IsStudent: scan: %w", err)
	}

	return &dto, nil
}
