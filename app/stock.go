package app

import (
	"errors"
	"github.com/shopspring/decimal"
	"reflect"
	"time"
)

type stockType int

// MEDICINE, FEED and ACCESSORY are the default stock types in madmin
const (
	MEDICINE stockType = iota
	FEED
	ACCESSORY
)

// StockTyper is an interface that wraps the StockType method.
// StockType returns the type of a stock item
type StockTyper interface {
	StockType() stockType
}

func (st stockType) StockType() stockType {
	return st
}

// Stock is a stock item interface
type Stock interface {
	ID() string
	Type() stockType

	Name() string
	SetName(string)

	IsExpirable() bool
	ExpirationDate() time.Time
	SetExpirationDate(time.Time)

	MinQuantity() decimal.Decimal
	SetMinQuantity(decimal.Decimal)

	Distributor() Distributor
	SetDistributor(d Distributor)
}

// NewStock creates a new valid Stock object.
// For now it only works for stock items of MEDICINE, FEED or ACCESSORY type.
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
		return nil, errors.New("invalid stock type")
	}
}

type defaultStock struct {
	id             string
	name           string
	minQuantity    decimal.Decimal
	expirationDate time.Time
	distributor    Distributor
}

func (ds defaultStock) ID() string {
	return ds.id
}
func (ds defaultStock) Name() string {
	return ds.name
}
func (ds defaultStock) SetName(name string) {
	ds.name = name
}
func (ds defaultStock) IsExpirable() bool {
	return true
}
func (ds defaultStock) ExpirationDate() time.Time {
	return ds.expirationDate
}
func (ds defaultStock) SetExpirationDate(expirationDate time.Time) {
	ds.expirationDate = expirationDate
}
func (ds defaultStock) MinQuantity() decimal.Decimal {
	return ds.minQuantity
}
func (ds defaultStock) SetMinQuantity(quantity decimal.Decimal) {
	ds.minQuantity = quantity
}
func (ds defaultStock) Distributor() Distributor {
	return ds.distributor
}
func (ds defaultStock) SetDistributor(d Distributor) {
	ds.distributor = d
}

type medicine struct {
	defaultStock
}

// NewMedicine creates a new stock item of medicine type (expirable stock item)
func NewMedicine(dto *NewStockDTO) (Stock, error) {
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
	distributor, err := NewDistributor(distributorName)
	if err != nil {
		return nil, err
	}

	return &medicine{defaultStock{id: id, name: name, minQuantity: quantity, expirationDate: date, distributor: distributor}}, err
}

func (m medicine) Type() stockType {
	return MEDICINE
}

type feed struct {
	defaultStock
}

// NewFeed creates a new stock item of feed type (expirable stock item)
func NewFeed(dto *NewStockDTO) (Stock, error) {
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
	distributor, err := NewDistributor(distributorName)
	if err != nil {
		return nil, err
	}

	return &feed{defaultStock{id: id, name: name, minQuantity: quantity, expirationDate: date, distributor: distributor}}, err
}

func (f feed) Type() stockType {
	return FEED
}

type accessory struct {
	defaultStock
}

// NewAccessory creates a new stock item of accessory type (unexpirable stock item).
// TODO: minimum quantity should be an integer, represented as decimal
// (should not have non-zero values after floating point)
func NewAccessory(dto *NewStockDTO) (Stock, error) {
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
			err = errors.New("error in creating stock item: Expiration date set for an accessory")
		}
	}

	distributorName := v.FieldByName("Distributor").String()
	distributor, err := NewDistributor(distributorName)
	if err != nil {
		return nil, err
	}

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
}
func (a accessory) SetExpirationDate(time.Time) {
	panic("Error - trying to set accessory's expiration date. Accessories do not expire.")
}
