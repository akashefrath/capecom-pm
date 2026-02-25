package dto

type CreateShiftSystemGroupRequest struct {
	Name          string `json:"name" form:"name" binding:"required"`
	ShiftSystemID string `json:"shift_system_id" form:"shift_system_id" binding:"required"`
}

type UpdateShiftSystemGroupRequest struct {
	Name          string `json:"name" form:"name" binding:"required"`
	ShiftSystemID string `json:"shift_system_id" form:"shift_system_id" binding:"required"`
}

type ShiftSystemGroupResponse struct {
	BaseModelTop
	Name          string `json:"name" db:"name"`
	ShiftSystemID string `json:"shift_system_id,omitempty" db:"shift_system_id"`
	BaseModelBottom
}

type AssignUsersToShiftGroupRequest struct {
	UserUUIDs []string `json:"user_ids" form:"user_ids" binding:"required,min=1"`
}
