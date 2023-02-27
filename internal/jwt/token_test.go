package jwt_test

import (
	"Social-app/internal/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerify(t *testing.T) {
	testCase := []struct {
		isValid bool
		token   string
	}{
		{
			isValid: true,
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOiIyMDIzLTA4LTI2IDExOjQxOjA2In0.n7P4-rqQuD0KBUT0EJx5cjVlRvJd_3Fqnl55TuWBdPU",
		},
	}
	for _, val := range testCase {
		if val.isValid {
			assert.NoError(t, jwt.Verify(val.token))
		} else {
			assert.Error(t, jwt.Verify(val.token))
		}
	}
}
