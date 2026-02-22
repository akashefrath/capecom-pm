package repository

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShiftSystemGroup struct {
	DB *sqlx.DB
}

func NewShiftSystemGroup(db *sqlx.DB) *ShiftSystemGroup {
	return &ShiftSystemGroup{DB: db}
}

func (r *ShiftSystemGroup) Create(req dto.CreateShiftSystemGroupRequest, shiftSystemID int64) (*int64, error) {
	q := `INSERT INTO shift_system_groups (uuid, name, shift_system_id) VALUES (?, ?, ?)`
	result, err := r.DB.Exec(q, uuid.New().String(), req.Name, shiftSystemID)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	return &id, err
}

func (r *ShiftSystemGroup) Update(uuid string, req dto.UpdateShiftSystemGroupRequest, shiftSystemID int64) error {
	q := `UPDATE shift_system_groups SET name = ?, shift_system_id = ? WHERE uuid = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(q, req.Name, shiftSystemID, uuid)
	return err
}

func (r *ShiftSystemGroup) Delete(uuid string) error {
	q := `UPDATE shift_system_groups SET deleted_at = NOW() WHERE uuid = ?`
	_, err := r.DB.Exec(q, uuid)
	return err
}

func (r *ShiftSystemGroup) GetAll() ([]dto.ShiftSystemGroupResponse, error) {
	var groups = make([]dto.ShiftSystemGroupResponse, 0)
	q := `SELECT id, uuid, name, shift_system_id, status FROM shift_system_groups WHERE deleted_at IS NULL`
	err := r.DB.Select(&groups, q)
	return groups, err
}

func (r *ShiftSystemGroup) GetByUUID(uuid string) (*dto.ShiftSystemGroupResponse, error) {
	var group dto.ShiftSystemGroupResponse
	q := `SELECT id, uuid, name, shift_system_id, status FROM shift_system_groups WHERE uuid = ? AND deleted_at IS NULL`
	err := r.DB.Get(&group, q, uuid)
	return &group, err
}

func (r *ShiftSystemGroup) GetByID(id int64) (*dto.ShiftSystemGroupResponse, error) {
	var group dto.ShiftSystemGroupResponse
	q := `SELECT id, uuid, name, shift_system_id, status FROM shift_system_groups WHERE id = ? AND deleted_at IS NULL`
	err := r.DB.Get(&group, q, id)
	return &group, err
}

func (r *ShiftSystemGroup) AssignUsers(groupID int64, userUUIDs []string) error {
	if len(userUUIDs) == 0 {
		return nil
	}

	q := `UPDATE users SET shift_system_group_id = ? WHERE uuid IN (?) AND deleted_at IS NULL`
	query, args, err := sqlx.In(q, groupID, userUUIDs)
	if err != nil {
		return err
	}
	query = r.DB.Rebind(query)
	_, err = r.DB.Exec(query, args...)
	return err
}

func (r *ShiftSystemGroup) RemoveUsers(userUUIDs []string) error {
	if len(userUUIDs) == 0 {
		return nil
	}

	q := `UPDATE users SET shift_system_group_id = NULL WHERE uuid IN (?) AND deleted_at IS NULL`
	query, args, err := sqlx.In(q, userUUIDs)
	if err != nil {
		return err
	}
	query = r.DB.Rebind(query)
	_, err = r.DB.Exec(query, args...)
	return err
}

func (r *ShiftSystemGroup) GetIDByUUID(uuid string) (*int64, error) {
	var id int64
	q := `SELECT id FROM shift_system_groups WHERE uuid = ? AND deleted_at IS NULL`
	err := r.DB.Get(&id, q, uuid)
	return &id, err
}

func (r *ShiftSystemGroup) GetUsersInGroup(groupUUID string) ([]dto.UserMinimalResponse, error) {
	var users = make([]dto.UserMinimalResponse, 0)
	q := `SELECT u.uuid, u.name, u.employee_id
	      FROM users u
	      INNER JOIN shift_system_groups ssg ON u.shift_system_group_id = ssg.id
	      WHERE ssg.uuid = ? AND u.deleted_at IS NULL AND ssg.deleted_at IS NULL`
	err := r.DB.Select(&users, q, groupUUID)
	return users, err
}
