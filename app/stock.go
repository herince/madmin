package madmin

import (
	"time"
	"github.com/shopspring/decimal"
)

type Stock interface {
	ID() string

	IsExpirable() bool
	ExpirationDate() time.Time
	ExpirationDate(time.Time)

	MinQuantity() decimal.Decimal
	MinQuantity(uint16)

	Distributor() string
}

type Medicine struct {
	id string
	expirable bool
}

func (m *Medicine) ID() string
func (m *Medicine) IsExpirable() bool {}
func (m *Medicine) GetExpirationDate() (time.Time, error) {}
func (m *Medicine) SetExpirationDate(time.Time) {}
func (m *Medicine) GetMinQuantity() (decimal.Decimal, error) {}
func (m *Medicine) SetMinQuantity(uint16) {}
func (m *Medicine) Distributor() string

type Feed struct {
	id string
}

func (f *Feed) ID() string
func (f *Feed) IsExpirable() bool {
	return false; // for now cow feed is unexpirable!
}
func (f *Feed) GetExpirationDate() (time.Time, error) {}
func (f *Feed) SetExpirationDate(time.Time) {}
func (f *Feed) GetMinQuantity() (decimal.Decimal, error) {}
func (f *Feed) SetMinQuantity(uint16) {}
func (f *Feed) Distributor() string


type Accessory struct {
	id string
}

func (a *Accessory) ID() string
func (a *Accessory) IsExpirable() bool {
	return false;
}
func (a *Accessory) GetExpirationDate() (time.Time, error) {}
func (a *Accessory) SetExpirationDate(time.Time) {}
func (a *Accessory) GetMinQuantity() (decimal.Decimal, error) {}
func (a *Accessory) SetMinQuantity(uint16) {}
func (a *Accessory) Distributor() string