package repository

import (
	"time"

	"github.com/akashefrath/capecom-pm/internal/database"
	"github.com/akashefrath/capecom-pm/internal/domain/common"
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/akashefrath/capecom-pm/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AttendanceSummary struct {
	DB   *sqlx.DB
	DBTx *database.Database
}

func NewAttendanceSummary(db *sqlx.DB, dbTX *database.Database) *AttendanceSummary {
	return &AttendanceSummary{DB: db, DBTx: dbTX}
}

func (r *AttendanceSummary) Create(tx *sqlx.Tx, userID int64, logDate time.Time) (*int64, error) {
	q := `INSERT INTO attendance_summary (uuid, user_id, log_date, log_status,last_log_type,last_log_at,status) 
		   VALUES  (?,?,?,?,?,?,?)
			`
	result, err := tx.Exec(q, uuid.NewString(), userID, logDate, "pending", models.LogIn, time.Now(), models.StatusActive)

	if err != nil {
		return nil, err
	}
	lastId, err := result.LastInsertId()

	return &lastId, err
}

func (r *AttendanceSummary) Update(tx *sqlx.Tx, summaryID int64, logStatus string, totalWorkInSec int64, totalBrakeInSec int64, lastLogType string) error {
	println(totalWorkInSec)
	println(totalBrakeInSec)
	q := `UPDATE attendance_summary
           SET log_status = ?,total_work_in_sec=?,total_brake_in_sec=? , last_log_type=? ,last_log_at=? 
            WHERE id = ?
			`
	_, err := tx.Exec(q, logStatus, totalWorkInSec, totalBrakeInSec, lastLogType, time.Now(), summaryID)

	if err != nil {
		return err
	}

	return err
}

func (r *AttendanceSummary) GetCurrentSummaryForUserIf(tx *sqlx.Tx, id int64) (*dto.AttendanceSummaryResponse, error) {
	var attendanceSummary dto.AttendanceSummaryResponse
	s, e := utils.GetTodayRange()
	q := `SELECT id,uuid,user_id,log_date,total_work_in_sec,total_brake_in_sec,log_status FROM attendance_summary WHERE user_id = ? 	 
	AND created_at >= ?
	  AND created_at < ?
	ORDER BY created_at ASC
	LIMIT 1
	`
	err := tx.Get(&attendanceSummary, q, id, s, e)

	return &attendanceSummary, err

}
func (r *AttendanceSummary) GetCurrentSummaryWithID(id int64) (*dto.AttendanceSummaryResponse, error) {
	var attendanceSummary dto.AttendanceSummaryResponse

	q := `SELECT id,uuid,user_id,log_date,total_work_in_sec,total_brake_in_sec,log_status FROM attendance_summary WHERE id = ?`
	err := r.DB.Get(&attendanceSummary, q, id)

	return &attendanceSummary, err

}
func (r *AttendanceSummary) GetCurrentPendingSummaryWithID(id int64) (*dto.AttendanceSummaryResponse, error) {
	var attendanceSummary dto.AttendanceSummaryResponse

	q := `SELECT id,uuid,user_id,log_date,total_work_in_sec,total_brake_in_sec,log_status,created_at FROM attendance_summary WHERE user_id = ? AND log_status =?`
	err := r.DB.Get(&attendanceSummary, q, id, "PENDING")

	return &attendanceSummary, err

}

func (r *AttendanceSummary) GetByUUID(uuid string) (*dto.AttendanceSummaryResponse, error) {
	var attendanceSummary dto.AttendanceSummaryResponse
	q := `SELECT id, uuid, user_id, log_date, total_work_in_sec, total_brake_in_sec, log_status, status, created_at, updated_at 
	      FROM attendance_summary 
	      WHERE uuid = ? AND deleted_at IS NULL`
	err := r.DB.Get(&attendanceSummary, q, uuid)
	return &attendanceSummary, err
}

func (r *AttendanceSummary) GetList(pagination common.Pagination, query dto.AttendanceSummaryListQuery) ([]dto.AttendanceSummaryResponse, error) {
	results := make([]dto.AttendanceSummaryResponse, 0)

	q := `SELECT id, uuid, user_id, log_date, total_work_in_sec, total_brake_in_sec, log_status, status, created_at, updated_at 
	      FROM attendance_summary 
	      WHERE deleted_at IS NULL`

	var args []interface{}

	if query.From != "" {
		q += ` AND log_date >= ?`
		args = append(args, query.From)
	}

	if query.To != "" {
		q += ` AND log_date <= ?`
		args = append(args, query.To)
	}

	allowedColumns := []string{"created_at", "log_date", "total_work_in_sec", "total_brake_in_sec"}
	orderQuery := query.OrderBy.GetQuery(allowedColumns)

	if orderQuery == "" {
		orderQuery = " ORDER BY created_at DESC"
	}

	q += orderQuery + ` LIMIT ? OFFSET ?`

	args = append(args, pagination.Limit, pagination.Offset())

	err := r.DB.Select(&results, q, args...)
	return results, err
}

func (r *AttendanceSummary) GetCount(query dto.AttendanceSummaryListQuery) (int64, error) {
	var count int64
	q := `SELECT COUNT(*) FROM attendance_summary WHERE deleted_at IS NULL`

	var args []interface{}

	if query.From != "" {
		q += ` AND log_date >= ?`
		args = append(args, query.From)
	}

	if query.To != "" {
		q += ` AND log_date <= ?`
		args = append(args, query.To)
	}

	err := r.DB.Get(&count, q, args...)
	return count, err
}
