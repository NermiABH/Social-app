package store_test

import (
	"Social-app/internal/model"
	"Social-app/internal/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostRepository_CreatePost(t *testing.T) {
	db, clearTables := store.TestDB(t, databaseURL)
	defer clearTables("users", "post")
	s := store.New(db)
	u := &model.User{
		Username: "user",
		Email:    "test@gmail.com",
		Password: "pusinu48",
	}
	_ = s.User().CreateUser(u)
	p := &model.Post{
		AuthorID: u.ID,
		Object:   "Something url",
		Text:     "Something text",
	}
	err := s.Post().CreatePost(p)
	assert.NoError(t, err)

}
