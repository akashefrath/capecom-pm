package models

type User struct {
	BaseModel

	Name  string
	Email string

	Phone       *string
	CountryCode *int

	PasswordHash string

	EmployeeID *string

	GroupID       uint64
	DesignationID uint64
	DepartmentID  uint64
}
