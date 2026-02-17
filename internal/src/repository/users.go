package repository

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
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

func (r *User) GetUserMinimalByEmail(email string) (*models.MinimalUser, error) {
	var user models.MinimalUser

	err := r.DB.Get(&user, `SELECT id,uuid,status,password_hash,is_admin FROM users WHERE email = ? AND deleted_at IS NULL AND status = ?`, email, models.SessionStatusActive)
	if err != nil {
		println(email)
		println(err.Error())
		return nil, domainerrors.ErrUserNotFound
	}

	return &user, err

}

func (r *User) FinUserByID(id int64) (*dto.User, error) {
	var res dto.User
	q := `SELECT id,uuid,name,email,phone,country_code,employee_id,status,created_at,updated_at,deleted_at,is_admin FROM users WHERE id = ?`

	err := r.DB.Get(&res, q, id)
	return &res, err
}

func (r *User) GetActiveUserUuidByID(id int64) *string {
	var userUuid string = ""
	q := `SELECT uuid FROM users WHERE id = ? AND status`
	_ = r.DB.Get(&userUuid, q, id, models.StatusActive)
	return &userUuid

}

func (r *User) GetActiveUserIDByUuid(uuid string) *int64 {
	var id int64
	q := `SELECT id FROM users WHERE uuid = ? AND status`
	_ = r.DB.Get(&id, q, id, models.StatusActive)
	return &id
}

func (r *User) FindUserStatus(id *int64) (*string, error) {
	var status string
	q := `SELECT status FROM users WHERE id = ?`
	_ = r.DB.Get(&status, q, id)
	return &status, nil

}
