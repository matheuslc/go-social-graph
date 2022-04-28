package entity

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	validUsername := "validuser"
	u, err := NewUser(validUsername)

	if err != nil && u.Username == validUsername {
		t.Errorf("Could not create a new user with a valid username")
	}

}

func TestInvalidNewUser(t *testing.T) {
	invalidUsername := "invalidinvalidduser"
	_, err := NewUser(invalidUsername)

	if err == nil {
		t.Errorf("Expect an error when trying to create a new user with userame greater than %d", UsernameMaxLength)
	}
}
