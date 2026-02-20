package repository

import (
	"time"

	"github.com/akashefrath/capecom-pm/internal/database"
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

func (t *AttendanceSummary) Create(tx *sqlx.Tx, userID int64, logDate time.Time) (*int64, error) {
	q := `INSERT INTO attendance_summary (uuid, user_id, log_date, log_status,status) 
		   VALUES  (?,?,?,?,?)
			`
	result, err := tx.Exec(q, uuid.NewString(), userID, logDate, "pending", models.StatusActive)
	if err != nil {
		return nil, err
	}
	lastId, err := result.LastInsertId()

	return &lastId, err
}

func (t *AttendanceSummary) Update(tx *sqlx.Tx, summaryID int64, logStatus string, totalWorkInSec int64, totalBrakeInSec int64) error {

	q := `UPDATE attendance_summary
           SET log_status = ?,total_work_in_sec=?,total_brake_in_sec=?
            WHERE id = ?
			`
	_, err := tx.Exec(q, logStatus, totalWorkInSec, totalBrakeInSec, summaryID)
	if err != nil {
		return err
	}

	return err
}

func (t *AttendanceSummary) GetCurrentSummaryForUserIf(tx *sqlx.Tx, id int64) (*dto.AttendanceSummaryResponse, error) {
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
