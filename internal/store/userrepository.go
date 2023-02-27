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
	return r.store.db.QueryRow("INSERT INTO users (username, email, encrypted_password) VALUES($1, $2, $3) RETURNING id",
		u.Username, u.Email, u.EncryptedPassword).Scan(&u.ID)
}

func (r *UserRepository) FindByUsername(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email=$1",
		email).Scan(&u.Email, &u.Password, &u.EncryptedPassword); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email=$1",
		email).Scan(&u.Email, &u.Password, &u.EncryptedPassword); err != nil {
		return nil, err
	}
	return u, nil
}
