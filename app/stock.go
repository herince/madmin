package app

import (
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"reflect"
	"time"
)

type stockType int

const (
	MEDICINE stockType = iota
	FEED
	ACCESSORY
)

// StockTyper is an interface that wraps the StockType method.
// StockType returns the type of a Stock object
type StockTyper interface {
	StockType() stockType
}

func (st stockType) StockType() stockType {
	return st
}

type Stock interface {
	Id() string
	Type() stockType

	Name() string
	SetName(string)

	IsExpirable() bool
	ExpirationDate() time.Time
	SetExpirationDate(time.Time)

	MinQuantity() decimal.Decimal
	SetMinQuantity(decimal.Decimal)

	Distributor() Distributor
	SetDistributor(distributor Distributor)
}

// NewStock creates a new valid Stock object.
func NewStock(dto *NewStockDTO) (Stock, error) {

	switch dto.Type {
	case MEDICINE:
		m, err := NewMedicine(dto)
		return m, err
	case FEED:
		f, err := NewFeed(dto)
		return f, err
	case ACCESSORY:
		a, err := NewAccessory(dto)
		return a, err
	default:
		return nil, errors.New("Invalid stock type.")
	}
}

type defaultStock struct {
	id             string
	name           string
	minQuantity    decimal.Decimal
	expirationDate time.Time
	distributor    Distributor
}

func (ds defaultStock) Id() string {
	return f.id
}
func (ds defaultStock) Name() string {
	return f.name
}
func (ds defaultStock) SetName(newName string) {
	f.name = newName
}
func (ds defaultStock) IsExpirable() bool {
	return true
}
func (ds defaultStock) ExpirationDate() time.Time {
	return f.expirationDate
}
func (ds defaultStock) SetExpirationDate(expirationDate time.Time) {
	f.expirationDate = expirationDate
}
func (ds defaultStock) MinQuantity() decimal.Decimal {
	return f.minQuantity
}
func (ds defaultStock) SetMinQuantity(quantity decimal.Decimal) {
	f.minQuantity = quantity
}
func (ds defaultStock) Distributor() Distributor {
	return f.distributor
}
func (ds defaultStock) SetDistributor(d Distributor) {
	f.distributor = d
}

type medicine struct {
	defaultStock
}

func NewMedicine(dto *NewStockDTO) (*medicine, error) {
	id, err := newUUID()
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(*dto)

	name := v.FieldByName("Name").String()

	quantityString := v.FieldByName("MinQuantity").String()
	var quantity decimal.Decimal
	if quantityString != "" {
		quantity, err = decimal.NewFromString(quantityString)
		if err != nil {
			return nil, err
		}
	} else {
		quantity = decimal.New(0, 0)
	}

	dateString := v.FieldByName("ExpirationDate").String()
	if dateString == "" {
		return nil, errors.New("No expiration date set for medicine")
	}
	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, dateString)
	if err != nil {
		return nil, err
	}

	distributorName := v.FieldByName("Distributor").String()
	distributor := Distributor(distributorName)

	return &medicine{defaultStock{id: id, name: name, minQuantity: quantity, expirationDate: date, distributor: distributor}}, err
}

func (m medicine) Type() stockType {
	return MEDICINE
}

type feed struct {
	defaultStock
}

func NewFeed(dto *NewStockDTO) (*feed, error) {
	id, err := newUUID()
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(*dto)

	name := v.FieldByName("Name").String()

	quantityString := v.FieldByName("MinQuantity").String()
	var quantity decimal.Decimal
	if quantityString != "" {
		quantity, err = decimal.NewFromString(quantityString)
		if err != nil {
			return nil, err
		}
	} else {
		quantity = decimal.New(0, 0)
	}

	dateString := v.FieldByName("ExpirationDate").String()
	if dateString == "" {
		return nil, errors.New("No expiration date set for feed")
	}
	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, dateString)
	if err != nil {
		return nil, err
	}

	distributorName := v.FieldByName("Distributor").String()
	distributor := Distributor(distributorName)

	return &feed{defaultStock{id: id, name: name, minQuantity: quantity, expirationDate: date, distributor: distributor}}, err
}

func (f feed) Type() stockType {
	return FEED
}

type accessory struct {
	defaultStock
}

func NewAccessory(dto *NewStockDTO) (*accessory, error) {
	id, err := newUUID()
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(*dto)

	name := v.FieldByName("Name").String()

	quantityString := v.FieldByName("MinQuantity").String()
	var quantity decimal.Decimal
	if quantityString != "" {
		quantity, err = decimal.NewFromString(quantityString)
		if err != nil {
			return nil, err
		}
	} else {
		quantity = decimal.New(0, 0)
	}

	dateString := v.FieldByName("ExpirationDate").String()
	if dateString != "" {
		if dateString != "" {
			err = errors.New("Error in creating stock item: Expiration date set for an accessory.")
		}
	}

	distributorName := v.FieldByName("Distributor").String()
	distributor := Distributor(distributorName)

	return &accessory{defaultStock{id: id, name: name, minQuantity: quantity, distributor: distributor}}, err
}

func (a accessory) Type() stockType {
	return ACCESSORY
}
func (a accessory) IsExpirable() bool {
	return false
}
func (a accessory) ExpirationDate() time.Time {
	panic("Error - trying to read accessory's expiration date. Accessories do not expire.")
	return time.Time{}
}
func (a accessory) SetExpirationDate(time.Time) {
	panic("Error - trying to set accessory's expiration date. Accessories do not expire.")
}