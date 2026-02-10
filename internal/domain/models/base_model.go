package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusBlocked  = "blocked"
	StatusArchived = "archived"
)

type BaseModel struct {
	gorm.Model
	ID uint64

	UUID string

	Status string `default:"active"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CreatedBy *uint64
}
type BaseModelNoCB struct {
	gorm.Model
	ID uint64

	UUID string

	Status string `default:"active"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
