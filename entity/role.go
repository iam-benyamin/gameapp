package entity

type Role uint8

const (
	UserRole Role = iota + 1
	AdminRole
)

const (
	userRoleStr  = "user"
	adminRoleStr = "admin"
)

func (r Role) String() string {
	switch r {
	case UserRole:
		return userRoleStr
	case AdminRole:
		return adminRoleStr
	}

	return ""
}

func MapToRoleEntity(roleStr string) Role {
	switch roleStr {
	case userRoleStr:
		return UserRole
	case adminRoleStr:
		return AdminRole
	}

	return Role(0)
}
