package dto

import "time"

type BaseModelTop struct {
	ID   uint64 `db:"id" json:"-"`
	UUID string `db:"uuid" json:"id"`
}

type BaseModelBottom struct {
	Status string `db:"status" default:"active" json:"status,omitempty"`

	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`

	CreatedBy *uint64 `db:"created_by" json:"-"`
}
