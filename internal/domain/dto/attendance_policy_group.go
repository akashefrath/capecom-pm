package dto

type AttendancePolicyGroupResponse struct {
	BaseModelTop
	Name                 string `json:"name" db:"name"`
	AttendancePolicyName string `json:"attendance_policy_name" db:"attendance_policy_name"`
	AttendancePolicyUUID string `json:"attendance_policy_uuid" db:"attendance_policy_uuid"`

	BaseModelBottom
}

type AttendancePolicyGroupSingleRow struct {
	BaseModelTop
	Name string `json:"name" db:"name"`
	AttendancePolicyResponseGroup
	BaseModelBottom
}
type AttendancePolicyResponseGroup struct {
	Name string `json:"name" db:"policy_name"`
	Uuid string `json:"id" db:"policy_uuid"`

	MinWorkHoursMinutes   int `json:"min_work_hours_minutes,omitempty" db:"policy_min_work_hours_minutes"`
	HalfDayMinutes        int `json:"half_day_minutes,omitempty" db:"policy_half_day_minutes"`
	LateGraceMinutes      int `json:"late_grace_minutes,omitempty" db:"policy_late_grace_minutes"`
	EarlyExitGraceMinutes int `json:"early_exit_grace_minutes,omitempty" db:"policy_early_exit_grace_minutes"`
	MaxBreakMinutes       int `json:"max_break_minutes,omitempty" db:"policy_max_break_minutes"`
	AutoCheckoutTime      int `json:"auto_checkout_time,omitempty" db:"policy_auto_checkout_time"`
}

type AttendancePolicyGroupSingleResponse struct {
	BaseModelTop
	Name                          string `json:"name" db:"name"`
	AttendancePolicyResponseGroup `json:"attendance_policy_response"`
	BaseModelBottom
}

func (a AttendancePolicyGroupSingleRow) AttendancePolicyGroupSingleResponse() AttendancePolicyGroupSingleResponse {
	return AttendancePolicyGroupSingleResponse{
		BaseModelTop: a.BaseModelTop,
		Name:         a.Name,
		AttendancePolicyResponseGroup: AttendancePolicyResponseGroup{
			Name:                  a.AttendancePolicyResponseGroup.Name,
			Uuid:                  a.AttendancePolicyResponseGroup.Uuid,
			MinWorkHoursMinutes:   a.MinWorkHoursMinutes,
			HalfDayMinutes:        a.HalfDayMinutes,
			LateGraceMinutes:      a.LateGraceMinutes,
			EarlyExitGraceMinutes: a.EarlyExitGraceMinutes,
			MaxBreakMinutes:       a.MaxBreakMinutes,
			AutoCheckoutTime:      a.AutoCheckoutTime,
		},
		BaseModelBottom: a.BaseModelBottom,
	}

}

type CreateAttendancePolicyGroupRequest struct {
	Name                 string `json:"name" form:"name" binding:"required,min=3,max=50"`
	AttendancePolicyUUID string `json:"attendance_policy_uuid" form:"attendance_policy_uuid" binding:"required,uuid"`
}

type UpdateAttendancePolicyGroupRequest struct {
	Name                 string `json:"name" form:"name" binding:"required,min=3,max=50"`
	AttendancePolicyUUID string `json:"attendance_policy_uuid" form:"attendance_policy_uuid" binding:"required,uuid"`
}
