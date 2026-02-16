package repository

import "database/sql"

type User struct {
	DB *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{
		DB: db,
	}
}

func (u *User) GetUserByEmail(email string) error {
	var uuid string
	var id int64
	var status string
	var password string
	err := u.DB.QueryRow("SELECT id,uuid,status,password_hash FROM users WHERE email = ? AND deleted_at IS NULL AND status = ?", email, "active").Scan(&id, &uuid, &status, &password)
	println(id)
	return err

}

func GetUserIDFromUuidQuery() string {
	return `SELECT id FROM users WHERE uuid  = ?`
}
