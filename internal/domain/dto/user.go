package dto

type User struct {
	BaseModelTop

	Name        string   `db:"name" json:"name"`
	Email       string   `db:"email" json:"email"`
	Phone       *string  `db:"phone" json:"phone"`
	CountryCode *int     `db:"country_code" json:"country_code"`
	EmployeeID  *string  `db:"employee_id" json:"employee_id"`
	IsAdmin     bool     `db:"is_admin" json:"is_admin"`
	Roles       []IDName `db:"roles" json:"roles"`
	BaseModelBottom
}

type IDName struct {
	BaseModelTop
	Name string `db:"name" json:"name"`
}

type CreateUserRequest struct {
	Name       string   ` json:"name" form:"name" binding:"required"`
	Email      string   ` json:"email" form:"email" binding:"required,email"`
	Phone      *string  ` json:"phone" form:"phone" `
	EmployeeID *string  ` json:"employee_id" form:"employee_id"`
	Roles      []string `json:"role_ids" form:"role_ids"`
	Password   string   `json:"password" form:"password"`
	//CountryCode *int    `db:"country_code" json:"country_code"`
}
