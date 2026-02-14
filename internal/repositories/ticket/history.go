package ticketrepo

import (
	"capecom-pm/internal/domain/dto"
	"capecom-pm/internal/domain/models"

	"gorm.io/gorm"
)

type HistoryRepo struct {
	DB *gorm.DB
}

func NewHistoryRepo(db *gorm.DB) *HistoryRepo {
	return &HistoryRepo{DB: db}
}

func (r *HistoryRepo) Create(history *models.TicketHistory) error {
	return r.DB.Table("ticket_history").Create(history).Error
}

func (r *HistoryRepo) GetAllByTicketID(ticketID int64, pagination string) (*[]dto.TicketHistoryResponse, error) {
	var results []dto.TicketHistoryResponse
	q := `SELECT th.uuid AS id, t.uuid AS ticket_id,
		u.uuid AS changed_by, u.name AS changed_by_name,
		th.field_name, th.old_value, th.new_value, th.note, th.created_at
		FROM ticket_history th
		INNER JOIN tickets t ON t.id = th.ticket_id
		INNER JOIN users u ON u.id = th.changed_by AND u.deleted_at IS NULL
		WHERE th.ticket_id = ?
		ORDER BY th.created_at DESC`
	if pagination != "" {
		q += ` ` + pagination
	}
	err := r.DB.Raw(q, ticketID).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (r *HistoryRepo) CountByTicketID(ticketID int64) (int64, error) {
	var count int64
	err := r.DB.Raw(`SELECT COUNT(*) FROM ticket_history WHERE ticket_id = ?`, ticketID).Scan(&count).Error
	return count, err
}
