package store_test

import (
	"Social-app/internal/model"
	"Social-app/internal/store"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	testCase := []*model.User{
		{Username: "test", Email: "test@gmail.com", Password: "pusinu48"},
		{Username: "test1", Email: "test1@gmail.com", Password: "pusinu48"},
	}
	for _, u := range testCase {
		assert.NoError(t, s.User().Create(u))
	}
}

func TestUserRepository_GetSeveralByUsername(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	testCase := []struct {
		part   string
		offset int
		limit  int
	}{
		{"t", 0, 5}, {"te", 0, 10},
	}
	for _, u := range testCase {
		uSlice, err := s.User().GetSeveralByUsername(u.part, u.limit, u.offset)
		assert.NotNil(t, uSlice)
		assert.NoError(t, err)
	}
}

func TestUserRepository_GetPasswordByUsername(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	u := &model.User{Username: "test"}
	assert.NoError(t, s.User().GetPasswordByUsername(u))
}

func TestUserRepository_GetPasswordByEmail(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	u := &model.User{Email: "test@gmail.com"}
	assert.NoError(t, s.User().GetPasswordByEmail(u))

}

func TestUserRepository_GetByID(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	u := &model.User{Username: "test"}
	_ = s.User().GetPasswordByUsername(u)
	assert.NoError(t, s.User().GetByID(u))
	fmt.Println(u)
}

func TestUserRepository_ExistByID(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	u := &model.User{Username: "test"}
	_ = s.User().GetPasswordByUsername(u)
	boolean, err := s.User().ExistByID(u.ID)
	assert.Equal(t, true, boolean)
	assert.NoError(t, err)
}

func TestUserRepository_ExistByUsername(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	boolean, err := s.User().ExistByUsername("test")
	assert.Equal(t, true, boolean)
	assert.NoError(t, err)
}

func TestUserRepository_ExistByEmail(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	boolean, err := s.User().ExistByEmail("test@gmail.com")
	assert.Equal(t, true, boolean)
	assert.NoError(t, err)
}

func TestUserRepository_SubscribeUser(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	u1, u2 := &model.User{Username: "test"}, &model.User{Username: "test1"}
	_, _ = s.User().GetPasswordByUsername(u1), s.User().GetPasswordByUsername(u2)
	assert.NoError(t, s.User().SubscribeUser(u1.ID, u2.ID))
}

func TestUserRepository_UnSubscribeUser(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	u1, u2 := &model.User{Username: "test"}, &model.User{Username: "test1"}
	_, _ = s.User().GetPasswordByUsername(u1), s.User().GetPasswordByUsername(u2)
	assert.NoError(t, s.User().UnSubscribeUser(u1.ID, u2.ID))
}

func TestUserRepository_SubscribedOrSubscriber(t *testing.T) {
	db, _ := TestDB(t, databaseURL)
	s := store.New(db)
	u1, u2 := &model.User{Username: "test"}, &model.User{Username: "test1"}
	_, _ = s.User().GetPasswordByUsername(u1), s.User().GetPasswordByUsername(u2)
	assert.NoError(t, s.User().SubscribedOrSubscriber(u1, u2.ID))
}

func TestUserRepository_Delete(t *testing.T) {
	db, teardown := TestDB(t, databaseURL)
	defer teardown("users")
	s := store.New(db)
	u := &model.User{Username: "test"}
	_ = s.User().GetPasswordByUsername(u)
	assert.NoError(t, s.User().Delete(u.ID))
}
