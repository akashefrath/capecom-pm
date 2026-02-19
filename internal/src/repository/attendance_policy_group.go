package repository

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AttendancePolicyGroup struct {
	DB *sqlx.DB
}

func NewAttendancePolicyGroup(db *sqlx.DB) *AttendancePolicyGroup {
	return &AttendancePolicyGroup{DB: db}
}

func (r *AttendancePolicyGroup) GetPolicyIDByUUID(policyUUID string) (*int64, error) {
	var id int64
	q := `SELECT id FROM attendance_policies WHERE uuid = ? AND deleted_at IS NULL`
	err := r.DB.Get(&id, q, policyUUID)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *AttendancePolicyGroup) Create(req dto.CreateAttendancePolicyGroupRequest, policyID int64) (*int64, error) {
	q := `INSERT INTO attendance_policy_groups (uuid, name, attendance_policy_id) VALUES (?, ?, ?)`
	result, err := r.DB.Exec(q, uuid.New().String(), req.Name, policyID)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *AttendancePolicyGroup) Update(uuid string, req dto.UpdateAttendancePolicyGroupRequest, policyID int64) error {
	q := `UPDATE attendance_policy_groups SET name = ?, attendance_policy_id = ? WHERE uuid = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(q, req.Name, policyID, uuid)
	return err
}

func (r *AttendancePolicyGroup) Delete(uuid string) error {
	q := `UPDATE attendance_policy_groups SET deleted_at = NOW() WHERE uuid = ?`
	_, err := r.DB.Exec(q, uuid)
	return err
}

func (r *AttendancePolicyGroup) GetAll() ([]dto.AttendancePolicyGroupSingleResponse, error) {
	var groups = make([]dto.AttendancePolicyGroupSingleRow, 0)
	q := `SELECT 
    	  apg.id, 
    	  apg.uuid,
    	  apg.name,
    	  
    	  ap.uuid as policy_uuid,
	      ap.name as policy_name, 
	      ap.min_work_hours_minutes as policy_min_work_hours_minutes,
	      ap.half_day_minutes as policy_half_day_minutes, 
	      ap.late_grace_minutes as policy_late_grace_minutes,
	      ap.early_exit_grace_minutes as policy_early_exit_grace_minutes, 
	      ap.max_break_minutes as policy_max_break_minutes, 
	      ap.auto_checkout_time as policy_auto_checkout_time, 
    	  apg.status 
	      FROM attendance_policy_groups apg
	      LEFT JOIN attendance_policies ap ON ap.id = apg.attendance_policy_id
	      WHERE apg.deleted_at IS NULL`
	err := r.DB.Select(&groups, q)
	if err != nil {
		return nil, err
	}
	var groupsSend = make([]dto.AttendancePolicyGroupSingleResponse, len(groups))

	for i, group := range groups {
		groupsSend[i] = group.AttendancePolicyGroupSingleResponse()
	}
	return groupsSend, err
}

func (r *AttendancePolicyGroup) GetByUUID(uuid string) (*dto.AttendancePolicyGroupSingleResponse, error) {
	var group dto.AttendancePolicyGroupSingleRow

	q := `SELECT 
    	    apg.id, 
    		apg.uuid, 
    		apg.name, 
    		ap.uuid as policy_uuid,
    		ap.name as policy_name, 
    		ap.min_work_hours_minutes as policy_min_work_hours_minutes,
    		ap.half_day_minutes as policy_half_day_minutes, 
    		ap.late_grace_minutes as policy_late_grace_minutes,
    		ap.early_exit_grace_minutes as policy_early_exit_grace_minutes, 
    		ap.max_break_minutes as policy_max_break_minutes, 
    		ap.auto_checkout_time as policy_auto_checkout_time, 
    		
    		apg.status 
	      FROM attendance_policy_groups apg
	      LEFT JOIN attendance_policies ap ON ap.id = apg.attendance_policy_id
	      WHERE apg.uuid = ? AND apg.deleted_at IS NULL`
	err := r.DB.Get(&group, q, uuid)
	finalStruct := group.AttendancePolicyGroupSingleResponse()
	return &finalStruct, err
}

func (r *AttendancePolicyGroup) GetByID(id int64) (*dto.AttendancePolicyGroupSingleResponse, error) {
	var group dto.AttendancePolicyGroupSingleRow

	q := `SELECT 
    	    apg.id, 
    		apg.uuid, 
    		apg.name, 
    		ap.uuid as policy_uuid,
    		ap.name as policy_name, 
    		
    		ap.min_work_hours_minutes as policy_min_work_hours_minutes,
    		ap.half_day_minutes as policy_half_day_minutes, 
    		ap.late_grace_minutes as policy_late_grace_minutes,
    		ap.early_exit_grace_minutes as policy_early_exit_grace_minutes, 
    		ap.max_break_minutes as policy_max_break_minutes, 
    		ap.auto_checkout_time as policy_auto_checkout_time, 
    		
    		apg.status 
	      FROM attendance_policy_groups apg
	      LEFT JOIN attendance_policies ap ON ap.id = apg.attendance_policy_id
	      WHERE apg.id = ? AND apg.deleted_at IS NULL`
	err := r.DB.Get(&group, q, id)
	finalStruct := group.AttendancePolicyGroupSingleResponse()
	return &finalStruct, err
}
