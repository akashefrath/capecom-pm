package models

type User struct {
	BaseModel

	Name  string
	Email string

	Phone       *string
	CountryCode *int

	PasswordHash string

	EmployeeID *string

	GroupID       int64
	DesignationID int64
	DepartmentID  int64
	IsAdmin       bool `gorm:"-"`
}
