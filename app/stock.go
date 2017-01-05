package madmin

import (
	"time"
)

type Stock interface {
	GetExpirationDate() (time.Time, error)
	SetExpirationDate(time.Time)

	GetMinQuantity() error
	SetMinQuantity(uint16)
}

