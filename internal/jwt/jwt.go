package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	AccessTime  = 15 * time.Minute
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
	if err := Verify(t.Refresh); err != nil {
		return err
	}
	payload := &Payload{}
	data, _ := base64.RawURLEncoding.DecodeString(strings.TrimLeft(strings.TrimRight(t.Access, "."), "."))
	json.Unmarshal(data, payload)
	t.CreateTokens(payload.UserID)
	return nil
}
