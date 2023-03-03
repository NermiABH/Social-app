package store_test

import (
	"Social-app/internal/model"
	"Social-app/internal/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db, clearTables := store.TestDB(t, databaseURL)
	defer clearTables("users")
	s := store.New(db)
	u := &model.User{
		Username: "user",
		Email:    "test@gmail.com",
		Password: "pusinu48",
	}
	err := s.User().CreateUser(u)
	assert.NoError(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, clearTables := store.TestDB(t, databaseURL)
	defer clearTables("users")
	s := store.New(db)
	u := &model.User{
		Username: "user",
		Email:    "test@gmail.com",
		Password: "pusinu48",
	}
	_ = s.User().CreateUser(u)
	u, err := s.User().FindByEmail("test@gmail.com")
	assert.NotNil(t, u)
	assert.NoError(t, err)
}
