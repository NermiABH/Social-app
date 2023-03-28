package store_test

import (
	"Social-app/internal/model"
	"Social-app/internal/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostRepository_Create(t *testing.T) {
	db, teardown := TestDB(t, databaseURL)
	defer teardown("post")
	s, u := store.New(db), &model.User{Username: "postuser", Email: "postuser@gmail.com", Password: "pusinu48"}
	_ = s.User().Create(u)
	text := "something text"
	p := &model.Post{Text: &text, AuthorID: u.ID}
	assert.NoError(t, s.Post().Create(p))
}
