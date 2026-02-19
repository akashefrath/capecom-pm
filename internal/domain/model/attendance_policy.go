package models

type AttendancePolicy struct {
	BaseModel
	MinWorkHoursMinutes   int    `db:"min_work_hours_minutes" json:"min_work_hours_minutes"`
	HalfDayMinutes        int    `db:"half_day_minutes" json:"half_day_minutes"`
	LateGraceMinutes      int    `db:"late_grace_minutes" json:"late_grace_minutes"`
	EarlyExitGraceMinutes int    `db:"early_exit_grace_minutes" json:"early_exit_grace_minutes"`
	MaxBreakMinutes       int    `db:"max_break_minutes" json:"max_break_minutes"`
	AutoCheckoutTime      int    `db:"auto_checkout_time" json:"auto_checkout_time"`
	CreatedBy             *int64 `db:"created_by" json:"created_by"`
}
