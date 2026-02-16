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
	ID uint64 `db:"id"`

	UUID string `db:"uuid"`

	Status string `db:"status" default:"active"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	CreatedBy *uint64
}
type BaseModelNoCB struct {
	ID uint64

	UUID string

	Status string `default:"active"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type BaseResponse struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
