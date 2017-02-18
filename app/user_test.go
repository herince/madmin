package app

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name             string
		password         string
		shouldCauseError bool
	}{
		{"name", "password", false},
		{"", "password", true},
		{"name", "", true},
	}

	for _, test := range tests {
		user, err := NewUser(test.name, test.password)
		if !test.shouldCauseError && err != nil {
			t.Fatalf(`
				NewUser returns an error for valid data.
				`)
		}
		if test.shouldCauseError && err == nil {
			t.Fatalf(`
				NewUser does not cause error for invalid data.
				Name=\"%s\", Password=\"%s\"
				`,
				test.name,
				test.password)
		}

		checkNewItemCreating(t, user, err)
	}
}

func TestSetPassword(t *testing.T) {
	tests := []struct {
		password         string
		shouldCauseError bool
	}{
		{"password", false},
		{"password123: парола", false},
		{"", true},
	}

	for _, test := range tests {
		user := &defaultUser{}
		err := user.SetPassword(test.password)
		if !test.shouldCauseError && err != nil {
			t.Fatalf(`
					Unexpected error %s. Getting error when setting a valid password.
					`, err)
		}
		if test.shouldCauseError && err == nil {
			t.Fatalf(`Expected error when setting invalid password.`)
		}

		if !test.shouldCauseError && err == nil {
			if ok := user.CheckPassword(test.password); !ok {
				t.Fatalf(`
				Error when setting password.
				Hash of passed password to CheckPassword() should equal hash of set password.
				`, )
			}
		}
	}
}
