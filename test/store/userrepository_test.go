package store_test

import (
	"Social-app/internal/model"
	"Social-app/internal/store"
	"fmt"
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
	err := s.User().Create(u)
	assert.NoError(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, clearTables := store.TestDB(t, databaseURL)
	defer clearTables("users")
	s := store.New(db)
	u := &model.User{
		Username: "user",
		Email:    "",
		Password: "pusinu48",
	}
	_ = s.User().Create(u)
	u.Email = "dsfadsfasdfadsf"
	err := s.User().GetPasswordByEmail(u)
	fmt.Println(err)
	assert.NotNil(t, u)
	assert.NoError(t, err)
}
