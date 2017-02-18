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

func checkNewItemCreatingWithOKStatus(t *testing.T, item interface{}, ok bool) {
	if ok && item == nil {
		t.Fatalf(`
			NewUser returns nil User and no error
		`)
	}
	if !ok && item != nil {
		t.Fatalf(`
			NewUser returns nil User and no error
		`)
	}
}

func defaultExpirableStockItem(aStockType stockType) (Stock, error) {
	dto := NewStockDTO{
		"name",
		aStockType,
		"1",
		"2030-01-01T00:00:00.000Z",
		"",
		"",
	}
	return NewStock(&dto)
}

func defaultUnexpirableStockItem(aStockType stockType) (Stock, error) {
	dto := NewStockDTO{
		"name",
		aStockType,
		"1",
		"",
		"",
		"",
	}
	return NewStock(&dto)
}