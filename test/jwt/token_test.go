package jwt_test

import (
	"Social-app/internal/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVerify(t *testing.T) {
	testCase := []struct {
		token string
	}{
		{
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOiIyMDIzLTAzLTAyIDE4OjEzOjUyIn0.GpgRG2Zsdjtz_-mJsKtnAyHAeZtqeQAXsh0IE_6MujE",
		},
	}
	emptyPayload := &jwt.Payload{}
	for _, val := range testCase {
		payload, err := jwt.Verify(val.token)
		if payload != emptyPayload && payload.Exp < time.Now().Format(time.DateTime) {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
