package utilsrepository

import (
	"github.com/jmoiron/sqlx"
)

type Utils struct {
	DB *sqlx.DB
}

func NewUtils(db *sqlx.DB) *Utils {
	return &Utils{DB: db}
}
