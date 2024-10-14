package helper

import (
	dm "ac_bot/internal/domain"
	"fmt"
	"strings"
)

const (
	formatGuest = `статус: *%s*
	имя: *%s* *%s*
	организация: *%s*
	приглашен: *%s*
	окончание: *%s*`

	formatStudent = `допуск: *%s*
	статус: *%s*
	имя: *%s* *%s*
	организация: *%s*
	поток: *%s*
	окончание: *%s*`
)

func GuestEntrCardMkdw(dto *dm.GuestEntranceCardDTO) string {

	return fmt.Sprintf(
		formatGuest,
		escapeMarkdown(dto.RoleName),
		escapeMarkdown(dto.LastName),
		escapeMarkdown(dto.FirstName),
		escapeMarkdown(dto.OrganizationTitle),
		escapeMarkdown(dto.AdminName),
		escapeMarkdown(dto.InviteExpireDate.Format("2006-01-02 15:04:05")),
	)
}

func StudentEntrCardMkdw(dto *dm.StudentEntranceCardDTO) string {
	approve := "ДА"
	if !dto.RegistrationApproved {
		approve = "НЕТ"
	}

	return fmt.Sprintf(
		formatStudent,
		approve,
		escapeMarkdown(dto.RoleName),
		escapeMarkdown(dto.LastName),
		escapeMarkdown(dto.FirstName),
		escapeMarkdown(dto.OrganizationTitle),
		escapeMarkdown(dto.InviteCode),
		escapeMarkdown(dto.InviteExpireDate.Format("2006-01-02")),
	)
}

func UnknownEntrCardMkdw() string {
	return `допуск: НЕТ
статус: \-
имя: \-
организация: \-`
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		`_`, `\_`,
		`*`, `\*`,
		`[`, `\[`,
		`]`, `\]`,
		`(`, `\(`,
		`)`, `\)`,
		`~`, `\~`,
		`>`, `\>`,
		`#`, `\#`,
		`+`, `\+`,
		`-`, `\-`,
		`=`, `\=`,
		`|`, `\|`,
		`{`, `\{`,
		`}`, `\}`,
		`.`, `\.`,
		`!`, `\!`,
	)
	return replacer.Replace(text)
}
