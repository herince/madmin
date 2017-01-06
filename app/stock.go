package madmin

import (
	"github.com/shopspring/decimal"
	"time"
)

type StockType int

const (
	MedicineStock StockType = iota
	FeedStock
	AccessoryStock
)

type Stock interface {
	ID() string
	Type() StockType

	Name() string
	Name(string)

	IsExpirable() bool
	ExpirationDate() time.Time
	ExpirationDate(time.Time)

	MinQuantity() decimal.Decimal
	MinQuantity(uint16)

	Distributor() string
}

// NewStock creates an unexpirable Stock object with the most basic needed information
func NewStock(s StockType,
	id string,
	name string) *Stock {

	return nil
}

// NewStock creates an expirable Stock object with the most basic needed information
func NewStock(s StockType,
	id string,
	name string,
	expirationDate time.Time) *Stock {

	return nil
}

// full way to add a new stock
func NewStock(s StockType,
	id string, name string,
	expirable bool,
	expirationDate time.Time,
	minQuantity decimal.Decimal,
	distributor string) *Stock {

	return nil
}

type Medicine struct {
	id             string
	name           string
	minQuantity    decimal.Decimal
	expirationDate time.Time
}

func (m *Medicine) ID() string {
	return m.id
}
func (m *Medicine) Type() StockType {
	return MedicineStock
}
func (m *Medicine) Name() string {
	return m.name
}
func (m *Medicine) Name(newName string) {
	m.name = newName
}
func (m *Medicine) IsExpirable() bool {
	return true
}
func (m *Medicine) ExpirationDate() time.Time {
	return m.expirationDate
}
func (m *Medicine) ExpirationDate(expirationDate time.Time) {
	m.expirationDate = expirationDate
}
func (m *Medicine) MinQuantity() decimal.Decimal {
	return m.minQuantity
}
func (m *Medicine) MinQuantity(minQuantity uint16) {
	m.minQuantity = minQuantity
}
func (m *Medicine) Distributor() string {
	return ""
}

type Feed struct {
	id             string
	name           string
	expirationDate time.Time
	minQuantity    decimal.Decimal
}

func (f *Feed) ID() string {
	return f.id
}
func (f *Feed) Type() StockType {
	return FeedStock
}
func (f *Feed) Name() string {
	return f.name
}
func (f *Feed) Name(newName string) {
	f.name = newName
}
func (f *Feed) IsExpirable() bool {
	return true
}
func (f *Feed) ExpirationDate() (time.Time, error) {
	return f.expirationDate, nil
}
func (f *Feed) ExpirationDate(time.Time) {
	return f.expirationDate
}
func (f *Feed) MinQuantity() decimal.Decimal {
	return f.minQuantity
}
func (f *Feed) MinQuantity(quantity uint16) {
	f.minQuantity = minQuantity
}
func (f *Feed) Distributor() string {
	return ""
}

type Accessory struct {
	id          string
	name        string
	minQuantity uint64
}

func (a *Accessory) ID() string {
	return a.id
}
func (a *Accessory) Type() StockType {
	return AccessoryStock
}
func (a *Accessory) Name() string {
	return a.name
}
func (a *Accessory) Name(newName string) {
	a.name = newName
}
func (a *Accessory) IsExpirable() bool {
	return false
}
func (a *Accessory) ExpirationDate() time.Time {
	return time.Time{}
}
func (a *Accessory) ExpirationDate(time.Time) {}
func (a *Accessory) MinQuantity() decimal.Decimal {
	return decimal.New(minQuantity, 0)
}
func (a *Accessory) MinQuantity(quantity decimal.Decimal) {
	a.minQuantity = quantity.IntPart()
}
func (a *Accessory) Distributor() string {
	return ""
}
