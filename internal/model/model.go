package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID                int    `json:"id,omitempty"`
	Username          string `json:"username,omitempty"`
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

func (u *User) BeforeCreate() error {
	s, err := encryptString(u.Password)
	if err != nil {
		return err
	}
	u.EncryptedPassword = s
	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
	u.EncryptedPassword = ""
}
func encryptString(s string) (string, error) {
	encrypt, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encrypt), nil
}
