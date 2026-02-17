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
	ID uint64 `db:"id" json:"-"`

	UUID string `db:"uuid" json:"id"`

	Status string `db:"status" default:"active" json:"status"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`

	CreatedBy *uint64 `db:"created_by" json:"-"`
}
type BaseModelNoCB struct {
	ID uint64 `db:"id"`

	UUID string `db:"uuid"`

	Status string `db:"status"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
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
