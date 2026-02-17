package dto

type User struct {
	BaseModelTop

	Name        string  `db:"name" json:"name"`
	Email       string  `db:"email" json:"email"`
	Phone       *string `db:"phone" json:"phone"`
	CountryCode *int    `db:"country_code" json:"country_code"`
	EmployeeID  *string `db:"employee_id" json:"employee_id"`
	IsAdmin     bool    `db:"is_admin" json:"is_admin"`

	BaseModelBottom
}
