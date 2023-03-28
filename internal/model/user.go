package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID                int
	Username          string
	Email             string
	Password          string
	EncryptedPassword string
	Userpic           string
	Name              *string
	Surname           *string
	Born              *string
	Created           string
	SubscriptionCount int
	SubscriberCount   int
	PostCount         int
	FavoritePostCount int
	Subscribed        bool
	Subscriber        bool
	Own               bool
}

func (u *User) ConvertMap() map[string]any {
	return map[string]any{
		"type": "user",
		"id":   u.ID,
		"attributes": map[string]any{
			"username":          u.Username,
			"email":             u.Email,
			"userpic":           u.Userpic,
			"name":              u.Name,
			"surname":           u.Surname,
			"born":              u.Born,
			"created":           u.Created,
			"subscriptionCount": u.SubscriptionCount,
			"subscriberCount":   u.SubscriberCount,
			"postCount":         u.PostCount,
			"favoritePostCount": u.FavoritePostCount,
		},
		"relationships": map[string]any{
			"subscribed": u.Subscribed,
			"subscriber": u.Subscriber,
			"own":        u.Own,
		},
	}
}

func (u *User) BeforeCreate() error {
	s, err := encryptString(u.Password)
	if err != nil {
		return err
	}
	u.EncryptedPassword = s
	return nil
}

func encryptString(s string) (string, error) {
	encrypt, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encrypt), nil
}

type Users []*User

func (u Users) ConvertMap() []map[string]any {
	users := make([]map[string]any, len(u))
	for i, user := range u {
		users[i] = map[string]any{
			"type": "user",
			"id":   user.ID,
			"attributes": map[string]any{
				"username": user.Username,
				"userpic":  user.Userpic,
			},
		}
	}
	return users
}
