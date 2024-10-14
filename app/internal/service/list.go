package service

import (
	dm "ac_bot/internal/domain"
	"bytes"
	"context"
	"fmt"
	"log"
)

type Lists interface {
	GetList(ctx context.Context, telegramID int64, path string) (string, error)
}

func (s *service) GetList(ctx context.Context, telegramID int64, path string) (string, error) {
	// verify that the request is from an administrator
	_, err := s.isAdmin(ctx, telegramID)
	if err != nil {
		return "", fmt.Errorf("svc: GetList: %w", err) // wrap!
	}

	switch path {
	case dm.InPathOrgs:
		return s.getOrgListF(ctx, "id: %d\ntitle: %s\n\n", 0)
	case dm.InPathRoles:
		return s.getRoleListF(ctx, "id: %d\nname: %s\n\n", 0, false)
	}

	return "", nil
}

func (s *service) getOrgListF(ctx context.Context, format string, sep rune) (string, error) {
	list, err := s.storage.GetOrgList(ctx)
	if err != nil {
		return "", err
	}

	if len(list) == 0 {
		log.Println("empty list of organizations!!!")
		return "", nil
	}

	buf := bytes.Buffer{}
	line := ""
	lastIndex := len(list) - 1

	for i, member := range list {
		line = fmt.Sprintf(format, member.ID, member.Title)
		buf.WriteString(line)
		if sep != 0 && i < lastIndex {
			buf.WriteRune(sep)
		}
	}

	return buf.String(), nil
}

func (s *service) getRoleListF(ctx context.Context, format string, sep rune, admin bool) (string, error) {
	list, err := s.storage.GetRoleList(ctx, admin)
	if err != nil {
		return "", err
	}

	if len(list) == 0 {
		log.Println("empty list of roles!!!")
		return "", nil
	}

	buf := bytes.Buffer{}
	line := ""
	lastIndex := len(list) - 1

	for i, role := range list {
		line = fmt.Sprintf(format, role.ID, role.Name)
		buf.WriteString(line)
		if sep != 0 && i < lastIndex {
			buf.WriteRune(sep)
		}
	}

	return buf.String(), nil
}
