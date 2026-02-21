package dto

import "github.com/akashefrath/capecom-pm/internal/domain/common"

type AttendanceSummaryResponse struct {
	BaseModelTop
	UserID          int64  `json:"-" db:"user_id"`
	LogDate         string `json:"log_date,omitempty" db:"log_date"`
	TotalWorkInSec  int64  `json:"total_work_in_sec,omitempty" db:"total_work_in_sec"`
	TotalBrakeInSec int64  `json:"total_brake_in_sec,omitempty" db:"total_brake_in_sec"`
	LogStatus       string `json:"log_status,omitempty" db:"log_status"`
	BaseModelBottom
}

type AttendanceSummaryListQuery struct {
	From string `json:"from" form:"from" query:"from"`
	To   string `json:"to" form:"to" query:"to"`
	common.OrderBy
}
