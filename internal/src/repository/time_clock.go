package repository

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TimeClock struct {
	DB *sqlx.DB
}

func NewTimeClock(db *sqlx.DB) *TimeClock {
	return &TimeClock{DB: db}
}

func (r *TimeClock) ClockIn(employeeID int64, req dto.TimeClockRequest) (*int64, error) {
	q := `INSERT INTO attendance_logs (uuid,employee_id, log_type, source, latitude, longitude, device_id, remarks) 
	      VALUES (?,?, 'IN', ?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(q, uuid.NewString(), employeeID, req.Source, req.Latitude, req.Longitude, req.DeviceID, req.Remarks)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	return &id, err
}

func (r *TimeClock) ClockOut(employeeID int64, req dto.TimeClockRequest) (*int64, error) {
	q := `INSERT INTO attendance_logs (uuid,employee_id, log_type, source, latitude, longitude, device_id, remarks) 
	      VALUES (?,?, 'OUT', ?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(q, uuid.NewString(), employeeID, req.Source, req.Latitude, req.Longitude, req.DeviceID, req.Remarks)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	return &id, err
}

func (r *TimeClock) BreakIn(employeeID int64, req dto.TimeClockRequest) (*int64, error) {
	q := `INSERT INTO attendance_logs (uuid,employee_id, log_type, source, latitude, longitude, device_id, remarks) 
	      VALUES (?,?, 'BREAK_IN', ?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(q, uuid.NewString(), employeeID, req.Source, req.Latitude, req.Longitude, req.DeviceID, req.Remarks)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	return &id, err
}

func (r *TimeClock) BreakOut(employeeID int64, req dto.TimeClockRequest) (*int64, error) {
	q := `INSERT INTO attendance_logs (uuid,employee_id, log_type, source, latitude, longitude, device_id, remarks) 
	      VALUES (?,?, 'BREAK_OUT', ?, ?, ?, ?, ?)`
	result, err := r.DB.Exec(q, uuid.NewString(), employeeID, req.Source, req.Latitude, req.Longitude, req.DeviceID, req.Remarks)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	return &id, err
}

func (r *TimeClock) GetByID(id int64) (*dto.TimeClockResponse, error) {
	var result dto.TimeClockResponse
	q := `SELECT al.id,al.uuid, al.employee_id, al.log_time, al.log_type, al.source, al.latitude, al.longitude, al.device_id, al.remarks , users.uuid as user_uuid
	      FROM attendance_logs  al
	      LEFT JOIN users ON al.employee_id = users.id
		  WHERE al.id = ?
	      `
	err := r.DB.Get(&result, q, id)
	return &result, err
}

func (r *TimeClock) GetUsersLastLog(id int64) (*string, error) {
	var logType *string
	q := `
	SELECT log_type
	FROM attendance_logs
	WHERE employee_id = ?
	  AND created_at >= CURDATE()
	  AND created_at < CURDATE() + INTERVAL 1 DAY
	ORDER BY created_at DESC
	LIMIT 1
	`
	err := r.DB.Get(&logType, q, id)
	return logType, err

}
