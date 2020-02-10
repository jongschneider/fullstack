package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {

	tests := []struct {
		shouldError bool
		user        *User
	}{
		{
			shouldError: true,
			user:        &User{},
		},
		{
			shouldError: true,
			user: &User{
				Nickname: "Jonny",
				Email:    "hedge",
			},
		},
		{
			shouldError: false,
			user: &User{
				Nickname: "Jonny",
				Email:    "hedge@google.com",
			},
		},
	}

	for _, test := range tests {
		var didError bool
		if err := test.user.Validate(); err != nil {
			didError = true
		}
		assert.Equal(t, test.shouldError, didError)
	}
}
