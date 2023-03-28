package jwt

import (
	"time"
)

const (
	AccessTime  = time.Hour
	RefreshTime = time.Hour * 24 * 180
)

type Tokens struct {
	UserID  int
	Access  string
	Refresh string
}

func (t *Tokens) ConvertMap() []map[string]any {
	return []map[string]any{
		{"type": "access",
			"token":    t.Access,
			"user:id":  t.UserID,
			"timelife": AccessTime},
		{"type": "refresh",
			"token":    t.Refresh,
			"user:id":  t.UserID,
			"timelife": RefreshTime},
	}
}

func New() *Tokens {
	return &Tokens{}
}

func (t *Tokens) CreateTokens(id int) {
	t.UserID, t.Access, t.Refresh = id, Create(id, AccessTime), Create(id, RefreshTime)
}

func (t *Tokens) RecreateTokens() error {
	payload, err := Verify(t.Refresh)
	if err != nil {
		return err
	}
	t.CreateTokens(payload.UserID)
	return nil
}
