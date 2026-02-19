package dto

type CreateAttendancePolicyRequest struct {
	Name                  string `json:"name" form:"name" binding:"required"`
	MinWorkHoursMinutes   int    `json:"min_work_hours_minutes" form:"min_work_hours_minutes" binding:"required,min=0"`
	HalfDayMinutes        int    `json:"half_day_minutes" form:"half_day_minutes" binding:"required,min=0"`
	LateGraceMinutes      int    `json:"late_grace_minutes" form:"late_grace_minutes" binding:"required,min=0"`
	EarlyExitGraceMinutes int    `json:"early_exit_grace_minutes" form:"early_exit_grace_minutes" binding:"required,min=0"`
	MaxBreakMinutes       int    `json:"max_break_minutes" form:"max_break_minutes" binding:"required,min=0"`
	AutoCheckoutTime      int    `json:"auto_checkout_time" form:"auto_checkout_time" binding:"required,min=0"`
	IsDefault             bool   `json:"is_default" form:"is_default"`
}

type AttendancePolicyResponse struct {
	BaseModelTop
	Name                  string `json:"name" db:"name"`
	MinWorkHoursMinutes   int    `json:"min_work_hours_minutes,omitempty" db:"min_work_hours_minutes"`
	HalfDayMinutes        int    `json:"half_day_minutes,omitempty" db:"half_day_minutes"`
	LateGraceMinutes      int    `json:"late_grace_minutes,omitempty" db:"late_grace_minutes"`
	EarlyExitGraceMinutes int    `json:"early_exit_grace_minutes,omitempty" db:"early_exit_grace_minutes"`
	MaxBreakMinutes       int    `json:"max_break_minutes,omitempty" db:"max_break_minutes"`
	AutoCheckoutTime      int    `json:"auto_checkout_time,omitempty" db:"auto_checkout_time"`
	IsDefault             bool   `json:"is_default"   db:"is_default"`
	BaseModelBottom
}
