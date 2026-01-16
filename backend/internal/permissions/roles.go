package permissions

import "strings"

type RoleSet map[string]bool

func ParseRoles(role string) RoleSet {
	roles := RoleSet{}
	for _, part := range strings.FieldsFunc(strings.ToLower(role), func(r rune) bool {
		switch r {
		case ',', '/', ';', '|', '+':
			return true
		default:
			return r == ' '
		}
	}) {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		roles[trimmed] = true
	}
	return roles
}

func IsAdmin(role string) bool {
	roles := ParseRoles(role)
	return roles["админ"] || roles["admin"] || roles["владелец"] || roles["owner"]
}

func IsModerator(role string) bool {
	roles := ParseRoles(role)
	return roles["модератор"] || roles["moderator"]
}

func IsLeader(role string) bool {
	roles := ParseRoles(role)
	return roles["руководитель"]
}

func CanManageMembers(role string) bool {
	return IsAdmin(role)
}

func CanManageTasks(role string) bool {
	return IsAdmin(role) || IsModerator(role) || IsLeader(role)
}

func CanReviewCompletions(role string) bool {
	return IsAdmin(role) || IsModerator(role) || IsLeader(role)
}

func MustIncludeModerator(role string) bool {
	return IsLeader(role)
}

func IsDeveloper(role string) bool {
	roles := ParseRoles(role)
	return roles["разработчик"]
}
