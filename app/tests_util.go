package app

import "testing"

func checkNewItemCreating(t *testing.T, item interface{}, err error) {
	if err == nil && item == nil {
		t.Fatalf(`
			NewUser returns nil User and no error
		`)
	}
	if err != nil && item != nil {
		t.Fatalf(`
			NewUser returns nil User and no error
		`)
	}
}
