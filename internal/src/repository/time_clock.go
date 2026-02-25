package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/akashefrath/capecom-pm/internal/database"
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/akashefrath/capecom-pm/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TimeClock struct {
	DB                *sqlx.DB
	DBTx              *database.Database
	AttendanceSummary *AttendanceSummary
}

func NewTimeClock(db *sqlx.DB, dbTX *database.Database, summary *AttendanceSummary) *TimeClock {
	return &TimeClock{DB: db, DBTx: dbTX, AttendanceSummary: summary}
}

func (r *TimeClock) TimePunch(userID int64, req dto.TimeClockRequest, punchType string) (*int64, error) {
	var lastID int64
	err := r.DBTx.WithTx(context.Background(), func(tx *sqlx.Tx) error {
		var summaryID *int64
		pendingSummary, err := r.AttendanceSummary.GetCurrentPendingSummaryWithID(userID)

		if pendingSummary != nil && punchType == models.LogIn {
			return domainerrors.CantPerformThis
		}

		currentSummary, err := r.AttendanceSummary.GetCurrentSummaryForUserIf(tx, userID)
		if errors.Is(err, sql.ErrNoRows) && punchType != models.LogIn {
			return domainerrors.CantPerformThis
		}

		if punchType == models.LogIn {
			summaryID, err = r.AttendanceSummary.Create(tx, userID, time.Now())
			if err != nil {
				return err
			}
		} else {
			summaryID = &currentSummary.ID
			if summaryID == nil || *summaryID == 0 {
				return domainerrors.ErrInternal
			}
			lastLog, err := r.GetLastTimeDataWithSummaryIDTx(summaryID, tx)
			if err != nil {
				return err
			}
			punchTime := time.Now() // or request.PunchTime (better)

			duration := punchTime.Sub(lastLog.LogTime).Seconds()
			if duration < 0 {
				duration = 0
			}
			delta := int64(duration)

			// copy previous totals
			workSec := currentSummary.TotalWorkInSec
			breakSec := currentSummary.TotalBrakeInSec

			switch lastLog.LogType {

			case models.LogIn:
				if punchType == models.BrakeIn || punchType == models.LogOut || punchType == models.TimeOut {
					workSec += delta
				}

			case models.BrakeOut:
				if punchType == models.BrakeIn || punchType == models.LogOut || punchType == models.TimeOut {
					workSec += delta
				}

			case models.BrakeIn:
				if punchType == models.BrakeOut {
					breakSec += delta
				}
			}

			logStatus := "PENDING"
			if punchType == models.LogOut || punchType == models.TimeOut {
				logStatus = "COMPLETED"
			}

			err = r.AttendanceSummary.Update(tx, *summaryID, logStatus, workSec, breakSec, punchType)
			if err != nil {
				return err
			}

		}

		q := `INSERT INTO attendance_logs (uuid,user_id, log_type, source, latitude, longitude, device_id, remarks,attendance_summary_id) 
	      VALUES (?,?, ?, ?, ?, ?, ?, ?,?)`
		result, err := tx.Exec(q, uuid.NewString(), userID, punchType, req.Source, req.Latitude, req.Longitude, req.DeviceID, req.Remarks, summaryID)
		if err != nil {
			return err
		}

		lastID, err = result.LastInsertId()

		if err != nil {
			return err
		}

		return err

	})

	return &lastID, err
}
func (r *TimeClock) TimeOutPunch(userID int64, req dto.TimeClockRequest, punchType string, timeOutTime time.Time) (*int64, error) {
	var lastID int64
	err := r.DBTx.WithTx(context.Background(), func(tx *sqlx.Tx) error {
		var summaryID *int64
		pendingSummary, err := r.AttendanceSummary.GetCurrentPendingSummaryWithID(userID)

		if pendingSummary != nil && punchType == models.LogIn {
			return domainerrors.CantPerformThis
		}

		currentSummary, err := r.AttendanceSummary.GetCurrentSummaryWithID(userID)
		if errors.Is(err, sql.ErrNoRows) && punchType != models.LogIn {
			return domainerrors.CantPerformThis
		}

		if punchType == models.LogIn {
			summaryID, err = r.AttendanceSummary.Create(tx, userID, time.Now())
			if err != nil {
				return err
			}
		} else {
			summaryID = &currentSummary.ID
			if summaryID == nil || *summaryID == 0 {
				return domainerrors.ErrInternal
			}
			lastLog, err := r.GetLastTimeDataWithSummaryIDTx(summaryID, tx)
			if err != nil {
				return err
			}
			punchTime := time.Now() // or request.PunchTime (better)

			duration := punchTime.Sub(lastLog.LogTime).Seconds()
			if duration < 0 {
				duration = 0
			}
			delta := int64(duration)

			// copy previous totals
			workSec := currentSummary.TotalWorkInSec
			breakSec := currentSummary.TotalBrakeInSec

			switch lastLog.LogType {

			case models.LogIn:

				workSec += delta

			case models.BrakeOut:

				workSec += delta

			case models.BrakeIn:

				breakSec += delta

			}

			logStatus := "COMPLETED"

			err = r.AttendanceSummary.Update(tx, *summaryID, logStatus, workSec, breakSec, punchType)
			if err != nil {
				return err
			}

		}

		q := `INSERT INTO attendance_logs (uuid,user_id, log_type, source, latitude, longitude, device_id, remarks,attendance_summary_id,log_time) 
	      VALUES (?,?, ?, ?, ?, ?, ?, ?,?,?)`
		result, err := tx.Exec(q, uuid.NewString(), userID, punchType, req.Source, req.Latitude, req.Longitude, req.DeviceID, req.Remarks, summaryID, timeOutTime)
		if err != nil {
			return err
		}

		lastID, err = result.LastInsertId()

		if err != nil {
			return err
		}

		return err

	})

	return &lastID, err
}

func (r *TimeClock) GetByID(id int64) (*dto.TimeClockResponse, error) {
	var result dto.TimeClockResponse
	q := `SELECT al.id,al.uuid, al.user_id, al.log_time, al.log_type, al.source, al.latitude, al.longitude, al.device_id, al.remarks,al.attendance_summary_id , users.uuid as user_uuid
	      FROM attendance_logs  al
	      LEFT JOIN users ON al.user_id = users.id
		  WHERE al.id = ?
	      `
	err := r.DB.Get(&result, q, id)
	return &result, err
}

func (r *TimeClock) GetUsersLastLog(id int64) (*string, error) {
	s, e := utils.GetTodayRange()
	var logType *string
	q := `
	SELECT log_type
	FROM attendance_logs
	WHERE user_id = ?
	  AND created_at >= ?
	  AND created_at < ?
	ORDER BY created_at DESC
	LIMIT 1
	`
	err := r.DB.Get(&logType, q, id, s, e)
	return logType, err

}

func (r *TimeClock) GetTodayDetails(id *int64) (*dto.AttendanceDetailsListWithSummary, error) {
	var attendances = make([]dto.AttendanceDetails, 0)
	var summary *dto.AttendanceSummaryResponse
	s, e := utils.GetTodayRange()
	q := `SELECT uuid,log_time,log_type,source,remarks,attendance_summary_id
		FROM attendance_logs 
		WHERE user_id = ? 	 
	AND created_at >= ?
	  AND created_at < ?
	ORDER BY created_at ASC`
	err := r.DB.Select(&attendances, q, id, s, e)
	if len(attendances) != 0 {
		summary, err = r.AttendanceSummary.GetCurrentSummaryWithID(attendances[0].AttendanceSummaryID)
	}

	pendingSummary, err := r.AttendanceSummary.GetCurrentPendingSummaryWithID(*id)

	return &dto.AttendanceDetailsListWithSummary{
		TimeLogs:          attendances,
		AttendanceSummary: summary,
		PendingSummary:    pendingSummary,
	}, err

}

func (r *TimeClock) GetLastTimeDataWithSummaryID(summaryID *int64) (*dto.TimeClockResponse, error) {

	var result dto.TimeClockResponse
	q := `SELECT al.id,al.uuid, al.user_id, al.log_time, al.log_type, al.source, al.latitude, al.longitude, al.device_id, al.remarks,al.attendance_summary_id , users.uuid as user_uuid
	      FROM attendance_logs  al
		  WHERE al.attendance_summary_id = ?
		  ORDER BY created_at DESC 
		  LIMIT 1
	      `
	err := r.DB.Get(&result, q, summaryID)
	return &result, err
}
func (r *TimeClock) GetLastTimeDataWithSummaryIDTx(summaryID *int64, tx *sqlx.Tx) (*dto.TimeClockResponse, error) {

	var result dto.TimeClockResponse
	q := `SELECT al.id,al.uuid, al.user_id, al.log_time, al.log_type, al.source, al.latitude, al.longitude, al.device_id, al.remarks,al.attendance_summary_id 
	      FROM attendance_logs  al
		  WHERE al.attendance_summary_id = ?
		  ORDER BY created_at DESC 
		  LIMIT 1
	      `
	err := tx.Get(&result, q, summaryID)
	return &result, err
}
