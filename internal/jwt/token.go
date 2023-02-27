package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	SecretKey            = "pusinu48"
	ErrorStructToken     = errors.New("struct of token is not valid")
	ErrorTokenIsNotValid = errors.New("token is not valid")
	ErrorTokenIsOutdated = errors.New("token is outdated")
	ErrorTokenIsEmpy     = errors.New("token is empty")
	varHeader            = &Header{Alg: "HS256", Typ: "JWT"}
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}
type Payload struct {
	UserID int    `json:"user_id"`
	Exp    string `json:"exp"`
}

// Create ...
func Create(id int, lifeTime time.Duration) string {
	header, _ := json.Marshal(varHeader)
	payload, _ := json.Marshal(
		Payload{
			UserID: id,
			Exp:    time.Now().Add(lifeTime).Format(time.DateTime),
		})
	unsignedToken := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	return unsignedToken + "." + generateHmac256(unsignedToken, SecretKey)
}

// Verify checks the validity of the token.
func Verify(token string) error {
	if token == "" {
		return ErrorTokenIsEmpy
	}
	jwtParts := strings.Split(token, ".")
	if len(jwtParts) != 3 {
		return ErrorStructToken
	}
	payload := &Payload{}
	payloadJson, err := base64.RawURLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		return ErrorTokenIsNotValid
	}
	if err = json.Unmarshal(payloadJson, payload); err != nil {
		return ErrorTokenIsNotValid
	}
	if payload.Exp < time.Now().Format(time.DateTime) {
		return ErrorTokenIsOutdated
	}
	if generateHmac256(jwtParts[0]+"."+jwtParts[1], SecretKey) != jwtParts[2] {
		return ErrorTokenIsNotValid
	}
	return nil
}

func generateHmac256(message, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
