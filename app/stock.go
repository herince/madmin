package app

import (
	"errors"
	"github.com/shopspring/decimal"
	"reflect"
	"strconv"
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

type medicine struct {
	id             string
	name           string
	minQuantity    decimal.Decimal
	expirationDate time.Time
	distributor    Distributor
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

	return &medicine{id: id, name: name, minQuantity: quantity, expirationDate: date, distributor: distributor}, err
}

func (m medicine) Id() string {
	return m.id
}
func (m medicine) Type() stockType {
	return MEDICINE
}
func (m medicine) Name() string {
	return m.name
}
func (m medicine) SetName(newName string) {
	m.name = newName
}
func (m medicine) IsExpirable() bool {
	return true
}
func (m medicine) ExpirationDate() time.Time {
	return m.expirationDate
}
func (m medicine) SetExpirationDate(expirationDate time.Time) {
	m.expirationDate = expirationDate
}
func (m medicine) MinQuantity() decimal.Decimal {
	return m.minQuantity
}
func (m medicine) SetMinQuantity(quantity decimal.Decimal) {
	m.minQuantity = quantity
}
func (m medicine) Distributor() Distributor {
	return m.distributor
}
func (m medicine) SetDistributor(d Distributor) {
	m.distributor = d
}

type feed struct {
	id             string
	name           string
	expirationDate time.Time
	minQuantity    decimal.Decimal
	distributor    Distributor
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

	return &feed{id: id, name: name, minQuantity: quantity, expirationDate: date, distributor: distributor}, err
}

func (f feed) Id() string {
	return f.id
}
func (f feed) Type() stockType {
	return FEED
}
func (f feed) Name() string {
	return f.name
}
func (f feed) SetName(newName string) {
	f.name = newName
}
func (f feed) IsExpirable() bool {
	return true
}
func (f feed) ExpirationDate() time.Time {
	return f.expirationDate
}
func (f feed) SetExpirationDate(expirationDate time.Time) {
	f.expirationDate = expirationDate
}
func (f feed) MinQuantity() decimal.Decimal {
	return f.minQuantity
}
func (f feed) SetMinQuantity(quantity decimal.Decimal) {
	f.minQuantity = quantity
}
func (f feed) Distributor() Distributor {
	return f.distributor
}
func (f feed) SetDistributor(d Distributor) {
	f.distributor = d
}

type accessory struct {
	id          string
	name        string
	minQuantity int64
	distributor Distributor
}

func NewAccessory(dto *NewStockDTO) (*accessory, error) {
	id, err := newUUID()
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(*dto)

	name := v.FieldByName("Name").String()

	quantityString := v.FieldByName("MinQuantity").String()
	var quantity int64
	if quantityString != "" {
		quantity, err = strconv.ParseInt(quantityString, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	dateString := v.FieldByName("ExpirationDate").String()
	if dateString != "" {
		if dateString != "" {
			err = errors.New("Error in creating stock item: Expiration date set for an accessory.")
		}
	}

	distributorName := v.FieldByName("Distributor").String()
	distributor := Distributor(distributorName)

	return &accessory{id: id, name: name, minQuantity: quantity, distributor: distributor}, err
}

func (a accessory) Id() string {
	return a.id
}
func (a accessory) Type() stockType {
	return ACCESSORY
}
func (a accessory) Name() string {
	return a.name
}
func (a accessory) SetName(newName string) {
	a.name = newName
}
func (a accessory) IsExpirable() bool {
	return false
}
func (a accessory) ExpirationDate() time.Time {
	return time.Time{}
}
func (a accessory) SetExpirationDate(time.Time) {}
func (a accessory) MinQuantity() decimal.Decimal {
	return decimal.New(a.minQuantity, 0)
}
func (a accessory) SetMinQuantity(quantity decimal.Decimal) {
	a.minQuantity = quantity.IntPart()
}
func (a accessory) Distributor() Distributor {
	return a.distributor
}
func (a accessory) SetDistributor(d Distributor) {
	a.distributor = d
}
