package jwt

import (
	"fmt"
	"time"
)

const (
	AccessTime  = time.Hour
	RefreshTime = time.Hour * 24 * 180
)

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func New() *Tokens {
	return &Tokens{}
}

func (t *Tokens) CreateTokens(id int) {
	t.Access, t.Refresh = Create(id, AccessTime), Create(id, RefreshTime)
}

func (t *Tokens) RecreateTokens() error {
	fmt.Print()
	payload, err := Verify(t.Refresh)
	if err != nil {
		return err
	}
	t.CreateTokens(payload.UserID)
	return nil
}
