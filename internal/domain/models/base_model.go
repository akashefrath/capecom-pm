package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusBlocked  = "blocked"
	StatusArchived = "archived"
)

type BaseModel struct {
	ID uint64

	UUID string

	Status string `default:"active"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	CreatedBy *uint64
}

func NewBase(createdBy *uint64) BaseModel {
	now := time.Now()

	return BaseModel{
		UUID:      uuid.NewString(),
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: createdBy,
	}
}
