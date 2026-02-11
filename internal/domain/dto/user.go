package dto

import "capecom-pm/internal/domain/models"

type CreateUserRequest struct {
	Name          string   `json:"name" form:"name" binding:"required,min=2,max=120"`
	Email         string   `json:"email" form:"email" binding:"required,email,max=255"`
	Phone         *string  `json:"phone" form:"phone" binding:"omitempty,max=15"`
	CountryCode   *int     `json:"country_code" form:"country_code" binding:"omitempty"`
	Password      string   `json:"password" form:"password" binding:"required,min=6,max=255"`
	EmployeeID    *string  `json:"employee_id" form:"employee_id" binding:"omitempty,max=50"`
	GroupID       string   `json:"group_id" form:"group_id" binding:"required"`
	DesignationID string   `json:"designation_id" form:"designation_id" binding:"required"`
	DepartmentID  string   `json:"department_id" form:"department_id" binding:"required"`
	RoleIDs       []string `json:"role_ids" form:"role_ids" binding:"required,min=1"`
}

type UserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	//Phone       *string `json:"phone"`
	//CountryCode *int    `json:"country_code"`
	EmployeeID  *string      `json:"employee_id"`
	Group       UtilNameId   `json:"user_group" gorm:"embedded;embeddedPrefix:user_group_"`
	Designation UtilNameId   `json:"designation" gorm:"embedded;embeddedPrefix:designation_"`
	Department  UtilNameId   `json:"department" gorm:"embedded;embeddedPrefix:department_"`
	Roles       []UtilNameId `json:"roles,omitempty" gorm:"-"`

	models.BaseResponse
}

type UtilNameId struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
