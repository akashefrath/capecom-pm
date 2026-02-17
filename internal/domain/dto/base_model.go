package dto

import "time"

type BaseModelTop struct {
	ID   uint64 `db:"id" json:"-"`
	UUID string `db:"uuid" json:"id"`
}

type BaseModelBottom struct {
	Status string `db:"status" default:"active" json:"status"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`

	CreatedBy *uint64 `db:"created_by" json:"-"`
}
