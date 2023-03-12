package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID                 int    `json:"id,omitempty"`
	Username           string `json:"username,omitempty"`
	Email              string `json:"email,omitempty"`
	Userpic            string `json:"userpic,omitempty"`
	Password           string `json:"password,omitempty"`
	EncryptedPassword  string `json:"-"`
	Name               string `json:"name,omitempty"`
	Surname            string `json:"surname,omitempty"`
	DateBirth          string `json:"date_of_birth,omitempty"`
	DateCreation       string `json:"date_of_creation,omitempty"`
	SubscriptionsCount int    `json:"subscriptions_count,omitempty"`
	SubscribersCount   int    `json:"subscribers_count,omitempty"`
	PostsCount         int    `json:"posts_count,omitempty"`
	FavoritesPosts     int    `json:"favorites_posts,omitempty"`
	IsSubscription     bool   `json:"is_subscription,omitempty"`
	IsSubscriber       bool   `json:"is_subscriber,omitempty"`
	IsOwn              bool   `json:"is_own,omitempty"`
	IsLiked            bool   `json:"is_liked,omitempty"`
	IsDisliked         bool   `json:"is_disliked,omitempty"`
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
