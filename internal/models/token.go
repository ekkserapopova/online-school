package models

import (
	"time"
)

type TokenPayload struct {
	ID  int
	Exp time.Time
}
