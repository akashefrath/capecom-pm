package repository

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AttendancePolicy struct {
	DB *sqlx.DB
}

func NewAttendancePolicy(db *sqlx.DB) *AttendancePolicy {
	return &AttendancePolicy{DB: db}
}

func (r *AttendancePolicy) Create(req dto.CreateAttendancePolicyRequest) (*int64, error) {
	if req.IsDefault {
		unsetQ := `UPDATE attendance_policies SET is_default = FALSE WHERE is_default = TRUE AND deleted_at IS NULL`
		_, err := r.DB.Exec(unsetQ)
		if err != nil {
			return nil, err
		}
	}

	q := `INSERT INTO attendance_policies (uuid,name, min_work_hours_minutes, half_day_minutes, late_grace_minutes, early_exit_grace_minutes, max_break_minutes, auto_checkout_time, is_default) 
	      VALUES (?, ?,?, ?, ?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(q, uuid.New().String(), req.Name, req.MinWorkHoursMinutes, req.HalfDayMinutes, req.LateGraceMinutes, req.EarlyExitGraceMinutes, req.MaxBreakMinutes, req.AutoCheckoutTime, req.IsDefault)

	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	return &id, err
}

func (r *AttendancePolicy) Update(uuid string, req dto.CreateAttendancePolicyRequest) error {
	if req.IsDefault {
		unsetQ := `UPDATE attendance_policies SET is_default = FALSE WHERE is_default = TRUE AND deleted_at IS NULL AND uuid != ?`
		_, err := r.DB.Exec(unsetQ, uuid)
		if err != nil {
			return err
		}
	}

	q := `UPDATE attendance_policies SET name = ?, min_work_hours_minutes = ?, half_day_minutes = ?, late_grace_minutes = ?, early_exit_grace_minutes = ?, max_break_minutes = ?, auto_checkout_time = ?, is_default = ? 
	      WHERE uuid = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(q, req.Name, req.MinWorkHoursMinutes, req.HalfDayMinutes, req.LateGraceMinutes, req.EarlyExitGraceMinutes, req.MaxBreakMinutes, req.AutoCheckoutTime, req.IsDefault, uuid)
	return err
}

func (r *AttendancePolicy) Delete(uuid string) error {
	q := `UPDATE attendance_policies SET deleted_at = NOW() WHERE uuid = ?`
	_, err := r.DB.Exec(q, uuid)
	return err
}

func (r *AttendancePolicy) GetAll() ([]dto.AttendancePolicyResponse, error) {
	var policies = make([]dto.AttendancePolicyResponse, 0)
	q := `SELECT id, uuid,name, min_work_hours_minutes, half_day_minutes, late_grace_minutes, early_exit_grace_minutes, max_break_minutes, auto_checkout_time, is_default, status
	      FROM attendance_policies WHERE deleted_at IS NULL`
	err := r.DB.Select(&policies, q)
	return policies, err
}
func (r *AttendancePolicy) GetAllUtils() ([]dto.AttendancePolicyResponse, error) {
	var policies = make([]dto.AttendancePolicyResponse, 0)
	q := `SELECT uuid,name
	      FROM attendance_policies WHERE deleted_at IS NULL AND status = ?`
	err := r.DB.Select(&policies, q, models.StatusActive)
	return policies, err
}
func (r *AttendancePolicy) GetByUUID(uuid string) (*dto.AttendancePolicyResponse, error) {
	var policy dto.AttendancePolicyResponse
	q := `SELECT id, uuid,name, min_work_hours_minutes, half_day_minutes, late_grace_minutes, early_exit_grace_minutes, max_break_minutes, auto_checkout_time, is_default, status
	      FROM attendance_policies WHERE uuid = ? AND deleted_at IS NULL`
	err := r.DB.Get(&policy, q, uuid)
	return &policy, err
}
func (r *AttendancePolicy) GetByID(id int64) (*dto.AttendancePolicyResponse, error) {
	var policy dto.AttendancePolicyResponse
	q := `SELECT id, uuid,name, min_work_hours_minutes, half_day_minutes, late_grace_minutes, early_exit_grace_minutes, max_break_minutes, auto_checkout_time, is_default, status
	      FROM attendance_policies WHERE id = ? AND deleted_at IS NULL`
	err := r.DB.Get(&policy, q, id)
	return &policy, err
}

func (r *AttendancePolicy) SetDefault(uuid string) error {
	unsetQ := `UPDATE attendance_policies SET is_default = FALSE WHERE is_default = TRUE AND deleted_at IS NULL`
	_, err := r.DB.Exec(unsetQ)
	if err != nil {
		return err
	}

	setQ := `UPDATE attendance_policies SET is_default = TRUE WHERE uuid = ? AND deleted_at IS NULL`
	_, err = r.DB.Exec(setQ, uuid)
	return err
}
