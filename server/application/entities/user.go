package entities

type User struct {
	Username     string
	Password     string
	Role         uint8
	RecoveryCode string
}

func (u *User) IsAdmin() bool {
	return u.Role == 0
}
