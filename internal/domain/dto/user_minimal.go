package dto

type UserMinimalResponse struct {
	UUID       string  `json:"id" db:"uuid"`
	Name       string  `json:"name" db:"name"`
	EmployeeID *string `json:"employee_id" db:"employee_id"`
}
