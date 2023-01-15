package models

import (
	"time"
)

type File struct {
	Name       string
	CreateadAt time.Time
	UpdatedAt  time.Time
}
