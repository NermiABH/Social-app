package store

import (
	"Social-app/internal/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) CreateUser(u *model.User) error {
	err := u.BeforeCreate()
	if err != nil {
		return err
	}
	return r.store.db.QueryRow("INSERT INTO users (email, encrypted_password) VALUES($1, $2) RETURNING id",
		u.Email, u.EncryptedPassword).Scan(&u.ID)
}
