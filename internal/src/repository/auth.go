package repository

import "database/sql"

type Auth struct {
	DB *sql.DB
}

func NewAuth(db *sql.DB) *Auth {
	return &Auth{
		DB: db,
	}
}
