package repository

import (
	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

type User struct {
	DB *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{
		DB: db,
	}
}

func (u *User) GetUserMinimalByEmail(email string) (*models.MinimalUser, error) {
	var user models.MinimalUser

	err := u.DB.Get(&user, "SELECT id,uuid,status,password_hash FROM users WHERE email = ? AND deleted_at IS NULL AND status = ?", email, "active")
	if err != nil {
		return nil, domainerrors.ErrUserNotFound
	}

	return &user, err

}

func GetUserIDFromUuidQuery() string {
	return `SELECT id FROM users WHERE uuid  = ?`
}
