package app

import (
	"testing"
)

func expectTypeAndExpirationDate(t *testing.T, aStockType stockType) {
	stockItem, err := defaultExpirableStockItem(aStockType)

	checkNewItemCreating(t, stockItem, err)

	if stockItem.Type() != aStockType {
		t.Fatal(`
			Incorrect stock type. Got %s, expected %s.`,
			stockItem.Type(),
			aStockType,
		)
	}
}

func expectTypeAndNoExpirationDate(t *testing.T, aStockType stockType) {
	stockItem, err := defaultUnexpirableStockItem(aStockType)

	checkNewItemCreating(t, stockItem, err)

	if stockItem.Type() != aStockType {
		t.Fatal(`
			Incorrect stock type. Got %s, expected %s.`,
			stockItem.Type(),
			aStockType,
		)
	}
}

func TestNewStock_WithValidData(t *testing.T) {
	expectTypeAndExpirationDate(t, MEDICINE)
	expectTypeAndExpirationDate(t, FEED)
	expectTypeAndNoExpirationDate(t, ACCESSORY)
}

func expectExpirationDate(t *testing.T, aStockType stockType) {
	dto := NewStockDTO{
		"name",
		aStockType,
		"1",
		"",
		"",
		"",
	}
	stockItem, err := NewStock(&dto)

	checkNewItemCreating(t, stockItem, err)

	if err == nil {
		t.Fatal(`
			NewStock does not return error on expirable item without expiration date
		`)
	}
}

func TestNewMedStock_WithNoExpirationDate(t *testing.T) {
	expectExpirationDate(t, MEDICINE)
}

func TestNewFeedStock_WithNoExpirationDate(t *testing.T) {
	expectExpirationDate(t, FEED)
}

func expectNoExpirationDate(t *testing.T, aStockType stockType) {
	dto := NewStockDTO{
		"name",
		aStockType,
		"1",
		"2030-01-01T00:00:00.000Z",
		"",
		"",
	}
	stockItem, err := NewStock(&dto)

	checkNewItemCreating(t, stockItem, err)

	if err == nil {
		t.Fatal(`
			NewStock does not return error on unexpirable item with expiration date set
		`)
	}
}

func TestNewAccessory_WithExpirationDate(t *testing.T) {
	expectNoExpirationDate(t, ACCESSORY)
}

func TestNewStock_WithInvalidType(t *testing.T) {
	dto := NewStockDTO{
		"name",
		stockType(5),
		"1",
		"2030-01-01T00:00:00.000Z",
		"",
		"",
	}
	stockItem, err := NewStock(&dto)

	checkNewItemCreating(t, stockItem, err)

	if err == nil {
		t.Fatal(`
			NewStock does not return error for unsuported stock type in DTO
		`)
	}
}
