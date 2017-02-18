package app

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
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
	// IsExpirable() should always be chacked before trying to call ExpirationDate()
	// or SetExpitrationDate(). Trying to get or set expiration date of an unexpirable
	// stock causes panic.
	ExpirationDate() time.Time
	SetExpirationDate(time.Time)

	MinQuantity() decimal.Decimal
	SetMinQuantity(decimal.Decimal)

	Quantity() decimal.Decimal
	SetQuantity(decimal.Decimal)

	DistributorID() string
	SetDistributorID(string)

	Update(StockDTO) error
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
	quantity       decimal.Decimal
	expirationDate time.Time
	distributorID  string
}

func (ds *defaultStock) ID() string {
	return ds.id
}
func (ds *defaultStock) Name() string {
	return ds.name
}
func (ds *defaultStock) SetName(name string) {
	ds.name = name
}
func (ds *defaultStock) IsExpirable() bool {
	return true
}
func (ds *defaultStock) ExpirationDate() time.Time {
	return ds.expirationDate
}
func (ds *defaultStock) SetExpirationDate(expirationDate time.Time) {
	ds.expirationDate = expirationDate
}
func (ds *defaultStock) MinQuantity() decimal.Decimal {
	return ds.minQuantity
}
func (ds *defaultStock) SetMinQuantity(quantity decimal.Decimal) {
	ds.minQuantity = quantity
}
func (ds *defaultStock) Quantity() decimal.Decimal {
	return ds.quantity
}
func (ds *defaultStock) SetQuantity(quantity decimal.Decimal) {
	ds.quantity = quantity
}
func (ds *defaultStock) DistributorID() string {
	return ds.distributorID
}
func (ds *defaultStock) SetDistributorID(id string) {
	ds.distributorID = id
}

func (ds *defaultStock) Update(dto StockDTO) error {
	if ds.ID() != dto.ID {
		return errors.New("trying to update stock with different id")
	}

	ds.SetName(dto.Name)

	if ds.IsExpirable() {
		date, err := validDateFromString(dto.ExpirationDate)
		if err != nil {
			return err
		}
		ds.SetExpirationDate(date)
	}

	quantity, err := validQuantityFromString(dto.Quantity)
	if err != nil {
		return err
	}
	ds.SetQuantity(quantity)

	if dto.MinQuantity == "" {
		ds.SetMinQuantity(decimal.Zero)
	} else {
		minQuantity, err := validQuantityFromString(dto.MinQuantity)
		if err != nil {
			return err
		}
		ds.SetMinQuantity(minQuantity)
	}

	ds.SetDistributorID(dto.DistributorID)

	return nil
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

	ds := &defaultStock{id: id}

	ds.SetName(dto.Name)

	date, err := validDateFromString(dto.ExpirationDate)
	if err != nil {
		return nil, err
	}
	ds.SetExpirationDate(date)

	quantity, err := validQuantityFromString(dto.Quantity)
	if err != nil {
		return nil, err
	}
	ds.SetQuantity(quantity)

	if dto.MinQuantity == "" {
		ds.SetMinQuantity(decimal.Zero)
	} else {
		minQuantity, err := validQuantityFromString(dto.MinQuantity)
		if err != nil {
			return nil, err
		}
		ds.SetMinQuantity(minQuantity)
	}

	ds.SetDistributorID(dto.DistributorID)

	return &medicine{defaultStock: *ds}, err
}

func (m *medicine) Type() stockType {
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

	name := dto.Name

	quantityString := dto.Quantity
	if quantityString == "" {
		return nil, errors.New("no quantity set for stock")
	}

	quantity, err := decimal.NewFromString(quantityString)
	if err != nil {
		return nil, err
	}

	var (
		minQuantityString = dto.MinQuantity
		minQuantity       decimal.Decimal
	)
	if minQuantityString != "" {
		minQuantity, err = decimal.NewFromString(minQuantityString)
		if err != nil {
			return nil, err
		}
	} else {
		minQuantity = decimal.New(0, 0)
	}

	dateString := dto.ExpirationDate
	if dateString == "" {
		return nil, errors.New("no expiration date set for feed")
	}
	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, dateString)
	if err != nil {
		return nil, err
	}

	distributorID := dto.DistributorID

	return &feed{defaultStock{id: id, name: name, quantity: quantity, minQuantity: minQuantity, expirationDate: date, distributorID: distributorID}}, err
}

func (f *feed) Type() stockType {
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

	name := dto.Name

	quantityString := dto.Quantity
	if quantityString == "" {
		return nil, errors.New("no quantity set for stock")
	}

	quantity, err := decimal.NewFromString(quantityString)
	if err != nil {
		return nil, err
	}

	var (
		minQuantityString = dto.MinQuantity
		minQuantity       decimal.Decimal
	)
	if minQuantityString != "" {
		minQuantity, err = decimal.NewFromString(minQuantityString)
		if err != nil {
			return nil, err
		}
	} else {
		minQuantity = decimal.New(0, 0)
	}

	dateString := dto.ExpirationDate
	if dateString != "" {
		if dateString != "" {
			return nil, errors.New("error in creating stock item: expiration date set for an accessory")
		}
	}

	distributorID := dto.DistributorID

	return &accessory{defaultStock{id: id, name: name, quantity: quantity, minQuantity: minQuantity, distributorID: distributorID}}, err
}

func (a *accessory) Type() stockType {
	return ACCESSORY
}
func (a *accessory) IsExpirable() bool {
	return false
}
func (a *accessory) ExpirationDate() time.Time {
	panic("Error - trying to read accessory's expiration date. Accessories do not expire.")
}
func (a *accessory) SetExpirationDate(time.Time) {
	panic("Error - trying to set accessory's expiration date. Accessories do not expire.")
}

func compareStock(first, second Stock) bool {
	return first.ID() == second.ID() &&
		first.Type() == second.Type() &&
		first.Name() == second.Name() &&
		first.ExpirationDate() == second.ExpirationDate() &&
		first.Quantity().Cmp(second.Quantity()) == 0 &&
		first.MinQuantity().Cmp(second.MinQuantity()) == 0 &&
		first.DistributorID() == second.DistributorID()
}
