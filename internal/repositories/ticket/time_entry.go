package ticketrepo

import (
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"time"

	"gorm.io/gorm"
)

type TimeEntryRepo struct {
	DB *gorm.DB
}

func NewTimeEntryRepo(db *gorm.DB) *TimeEntryRepo {
	return &TimeEntryRepo{DB: db}
}

func (r *TimeEntryRepo) Create(entry *models.TimeEntry) (*dto.TimeEntryResponse, error) {
	if err := r.DB.Create(entry).Error; err != nil {
		return nil, err
	}
	return r.findByID(entry.ID)
}

func (r *TimeEntryRepo) FindByUUID(uuid string) (*dto.TimeEntryResponse, error) {
	var result dto.TimeEntryResponse
	err := r.DB.Raw(r.selectQuery("te.uuid"), uuid).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrTimeEntryNotFound
	}
	return &result, nil
}

func (r *TimeEntryRepo) findByID(id uint64) (*dto.TimeEntryResponse, error) {
	var result dto.TimeEntryResponse
	err := r.DB.Raw(r.selectQuery("te.id"), id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrTimeEntryNotFound
	}
	return &result, nil
}

func (r *TimeEntryRepo) GetAllByTicketID(ticketID int64, pagination string) (*[]dto.TimeEntryResponse, error) {
	var results []dto.TimeEntryResponse
	err := r.DB.Raw(r.selectQuery("te.ticket_id", "ORDER BY te.work_date DESC, te.created_at DESC", pagination), ticketID).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (r *TimeEntryRepo) CountByTicketID(ticketID int64) (int64, error) {
	var count int64
	err := r.DB.Raw(`SELECT COUNT(*) FROM time_entries WHERE ticket_id = ? AND deleted_at IS NULL`, ticketID).Scan(&count).Error
	return count, err
}

func (r *TimeEntryRepo) Update(uuid string, updates map[string]any) (*dto.TimeEntryResponse, error) {
	res := r.DB.Table("time_entries").Where("uuid = ? AND deleted_at IS NULL", uuid).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrTimeEntryNotFound
	}
	return r.FindByUUID(uuid)
}

func (r *TimeEntryRepo) Delete(uuid string) error {
	res := r.DB.Exec(`UPDATE time_entries SET deleted_at = NOW() WHERE uuid = ? AND deleted_at IS NULL`, uuid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrTimeEntryNotFound
	}
	return nil
}

func (r *TimeEntryRepo) GetTicketAssignedUserID(ticketID int64) (*int64, error) {
	var userID *int64
	err := r.DB.Raw(`SELECT assigned_to FROM tickets WHERE id = ? AND deleted_at IS NULL`, ticketID).Scan(&userID).Error
	if err != nil {
		return nil, err
	}
	return userID, nil
}

func (r *TimeEntryRepo) GetTimeEntryOwnerID(uuid string) (*int64, error) {
	var userID *int64
	err := r.DB.Raw(`SELECT user_id FROM time_entries WHERE uuid = ? AND deleted_at IS NULL`, uuid).Scan(&userID).Error
	if err != nil {
		return nil, err
	}
	if userID == nil {
		return nil, domainerrors.ErrTimeEntryNotFound
	}
	return userID, nil
}

func (r *TimeEntryRepo) GetTotalHoursByTicketID(ticketID int64) (float64, error) {
	var hours float64
	err := r.DB.Raw(`SELECT SUM(hours) FROM time_entries WHERE ticket_id = ? AND deleted_at IS NULL`, ticketID).Scan(&hours).Error
	return hours, err
}

func (r *TimeEntryRepo) GetTotalHoursByProjectID(projectID int64) (float64, error) {
	var hours float64
	err := r.DB.Raw(`SELECT SUM(hours) FROM time_entries WHERE project_id = ? AND deleted_at IS NULL`, projectID).Scan(&hours).Error
	return hours, err
}

func (r *TimeEntryRepo) GetTotalHoursByUserID(userID int64) (float64, error) {
	var hours float64
	err := r.DB.Raw(`SELECT SUM(hours) FROM time_entries WHERE user_id = ? AND deleted_at IS NULL`, userID).Scan(&hours).Error
	return hours, err
}
func (r *TimeEntryRepo) GetTotalHoursByUserIDByDate(userID int64, date time.Time) (*float64, error) {
	var hours *float64

	err := r.DB.Raw(`SELECT SUM(hours) FROM time_entries WHERE user_id = ? AND DATE(work_date) = DATE(?) AND deleted_at IS NULL`, userID, date).Scan(&hours).Error

	return hours, err
}
func (r *TimeEntryRepo) selectQuery(whereCol string, extra ...string) string {
	q := `SELECT te.uuid AS id, t.uuid AS ticket_id, t.code AS ticket_code,
		p.uuid AS project_id, u.uuid AS user_id, u.name AS user_name,
		te.work_date, te.hours, te.description, te.status,
		te.created_at, te.updated_at
		FROM time_entries te
		INNER JOIN tickets t ON t.id = te.ticket_id AND t.deleted_at IS NULL
		INNER JOIN projects p ON p.id = te.project_id AND p.deleted_at IS NULL
		INNER JOIN users u ON u.id = te.user_id AND u.deleted_at IS NULL
		WHERE te.deleted_at IS NULL`
	if whereCol != "" {
		q += ` AND ` + whereCol + ` = ?`
	}
	for _, e := range extra {
		q += ` ` + e
	}
	return q
}
