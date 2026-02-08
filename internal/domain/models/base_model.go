package models

import "time"

type BaseModel struct {
	ID uint64

	UUID string

	Status string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	CreatedBy *uint64
}
