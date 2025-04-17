package entities

import "fmt"

type User struct {
	ID           string
	Username     string
	Password     string
	RoleID       string
	RecoveryCode string
}

func (u *User) IsAdmin() bool {
	return u.RoleID == "0"
}

func (u *User) SetPassword(hash string) error {
	if len(hash) < 8 {
		return fmt.Errorf("password too short")
	}
	u.Password = hash
	return nil
}

func NewUser(id, username, password string, roleID, recoveryCode string) (User, error) {
	if len(id) < 1 || len(username) < 3 || len(password) < 8 {
		return User{}, fmt.Errorf("bad user format")
	}
	return User{ID: id, Username: username, Password: password, RoleID: roleID, RecoveryCode: recoveryCode}, nil
}
