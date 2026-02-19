package repository

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShiftSystem struct {
	DB *sqlx.DB
}

func NewShiftSystem(db *sqlx.DB) *ShiftSystem {
	return &ShiftSystem{DB: db}
}

func (r *ShiftSystem) Create(req dto.CreateShiftSystemRequest) (*int64, error) {
	if req.IsDefault {
		unsetQ := `UPDATE shift_system SET is_default = FALSE WHERE is_default = TRUE AND deleted_at IS NULL`
		_, err := r.DB.Exec(unsetQ)
		if err != nil {
			return nil, err
		}
	}

	q := `INSERT INTO shift_system (uuid, name, start_time, end_time, checkin_early, checkin_late, checkout_early, checkout_late, is_overnight, is_default) 
	      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(q, uuid.New().String(), req.Name, req.StartTime, req.EndTime, req.CheckinEarly, req.CheckinLate, req.CheckoutEarly, req.CheckoutLate, req.IsOvernight, req.IsDefault)

	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	return &id, err
}

func (r *ShiftSystem) Update(uuid string, req dto.UpdateShiftSystemRequest) error {
	if req.IsDefault {
		unsetQ := `UPDATE shift_system SET is_default = FALSE WHERE is_default = TRUE AND deleted_at IS NULL AND uuid != ?`
		_, err := r.DB.Exec(unsetQ, uuid)
		if err != nil {
			return err
		}
	}

	q := `UPDATE shift_system SET name = ?, start_time = ?, end_time = ?, checkin_early = ?, checkin_late = ?, checkout_early = ?, checkout_late = ?, is_overnight = ?, is_default = ? 
	      WHERE uuid = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(q, req.Name, req.StartTime, req.EndTime, req.CheckinEarly, req.CheckinLate, req.CheckoutEarly, req.CheckoutLate, req.IsOvernight, req.IsDefault, uuid)
	return err
}

func (r *ShiftSystem) Delete(uuid string) error {
	q := `UPDATE shift_system SET deleted_at = NOW() WHERE uuid = ?`
	_, err := r.DB.Exec(q, uuid)
	return err
}

func (r *ShiftSystem) GetAll() ([]dto.ShiftSystemResponse, error) {
	var shifts = make([]dto.ShiftSystemResponse, 0)
	q := `SELECT id, uuid, name, start_time, end_time, checkin_early, checkin_late, checkout_early, checkout_late, is_overnight, is_default, status
	      FROM shift_system WHERE deleted_at IS NULL`
	err := r.DB.Select(&shifts, q)
	return shifts, err
}

func (r *ShiftSystem) GetAllUtils() ([]dto.ShiftSystemResponse, error) {
	var shifts = make([]dto.ShiftSystemResponse, 0)
	q := `SELECT uuid, name
	      FROM shift_system WHERE deleted_at IS NULL AND status = ?`
	err := r.DB.Select(&shifts, q, models.StatusActive)
	return shifts, err
}

func (r *ShiftSystem) GetByUUID(uuid string) (*dto.ShiftSystemResponse, error) {
	var shift dto.ShiftSystemResponse
	q := `SELECT id, uuid, name, start_time, end_time, checkin_early, checkin_late, checkout_early, checkout_late, is_overnight, is_default, status
	      FROM shift_system WHERE uuid = ? AND deleted_at IS NULL`
	err := r.DB.Get(&shift, q, uuid)
	return &shift, err
}

func (r *ShiftSystem) GetByID(id int64) (*dto.ShiftSystemResponse, error) {
	var shift dto.ShiftSystemResponse
	q := `SELECT id, uuid, name, start_time, end_time, checkin_early, checkin_late, checkout_early, checkout_late, is_overnight, is_default, status
	      FROM shift_system WHERE id = ? AND deleted_at IS NULL`
	err := r.DB.Get(&shift, q, id)
	return &shift, err
}

func (r *ShiftSystem) SetDefault(uuid string) error {
	unsetQ := `UPDATE shift_system SET is_default = FALSE WHERE is_default = TRUE AND deleted_at IS NULL`
	_, err := r.DB.Exec(unsetQ)
	if err != nil {
		return err
	}

	setQ := `UPDATE shift_system SET is_default = TRUE WHERE uuid = ? AND deleted_at IS NULL`
	_, err = r.DB.Exec(setQ, uuid)
	return err
}

func (r *ShiftSystem) GetIDByUUID(uuid string) (*int64, error) {
	var id int64
	q := `SELECT id FROM shift_system WHERE uuid = ? AND deleted_at IS NULL`
	err := r.DB.Get(&id, q, uuid)
	return &id, err
}
