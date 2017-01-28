package app

import (
	"github.com/shopspring/decimal"
	"time"
	//"fmt"
	"errors"
)

type StockType int

const (
	MedicineStock StockType = iota
	FeedStock
	AccessoryStock
)

type Stock interface {
	Id() string
	Type() StockType

	Name() string
	SetName(string)

	IsExpirable() bool
	ExpirationDate() time.Time
	SetExpirationDate(time.Time)

	MinQuantity() decimal.Decimal
	SetMinQuantity(decimal.Decimal)

	Distributor() string
}

// NewStock creates an unexpirable Stock object with the most basic needed information
func NewStock(t StockType,
	name string) (Stock, error) {

	switch t {
	case MedicineStock:
		newItem, err := NewMedicine()
		newItem.SetName(name)
		return *newItem, err
	case FeedStock:
		newItem, err := NewFeed()
		newItem.SetName(name)
		return *newItem, err
	case AccessoryStock:
		newItem, err := NewAccessory()
		newItem.SetName(name)
		return *newItem, err
	default:
		return nil, errors.New("Invalid stock type")
	}
}

// NewStock creates an expirable Stock object with the most basic needed information
/*
func NewStock(s StockType,
	name string,
	expirationDate time.Time) *Stock {

	return nil
}*/

// full way to add a new stock
/*
func NewStock(s StockType,
	name string,
	expirable bool,
	expirationDate time.Time,
	minQuantity decimal.Decimal,
	distributor string) *Stock {

	s := &Stock{}

	id, err := newUUID()
	if err != nil {
		panic(err)
	}

	fmt.Println(id)

	return nil
}*/

type medicine struct {
	id             string
	name           string
	minQuantity    decimal.Decimal
	expirationDate time.Time
}

func NewMedicine() (*medicine, error) {
	var (
		id, err = newUUID()
		m       = &medicine{id: id}
	)
	return m, err
}

func (m medicine) Id() string {
	return m.id
}
func (m medicine) Type() StockType {
	return MedicineStock
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
func (m medicine) Distributor() string {
	return ""
}

type feed struct {
	id             string
	name           string
	expirationDate time.Time
	minQuantity    decimal.Decimal
}

func NewFeed() (*feed, error) {
	var (
		id, err = newUUID()
		f       = &feed{id: id}
	)
	return f, err
}

func (f feed) Id() string {
	return f.id
}
func (f feed) Type() StockType {
	return FeedStock
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
func (f feed) Distributor() string {
	return ""
}

type accessory struct {
	id          string
	name        string
	minQuantity int64
}

func NewAccessory() (*accessory, error) {
	var (
		id, err = newUUID()
		a       = &accessory{id: id}
	)
	return a, err
}

func (a accessory) Id() string {
	return a.id
}
func (a accessory) Type() StockType {
	return AccessoryStock
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
func (a accessory) Distributor() string {
	return ""
}
