package ticketrepo

import (
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type TicketRepo struct {
	DB *gorm.DB
}

func NewTicketRepo(db *gorm.DB) *TicketRepo {
	return &TicketRepo{DB: db}
}

func (r *TicketRepo) Create(ticket *models.Ticket) (*dto.TicketResponse, error) {
	var mysqlErr *mysql.MySQLError
	if err := r.DB.Create(ticket).Error; err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, domainerrors.NewWithCode(http.StatusConflict, domainerrors.ErrDuplicateTicket.Error(), "repo", "CreateTicket")
		}
		return nil, err
	}
	return r.findByID(ticket.ID)
}

func (r *TicketRepo) FindByUUID(uuid string) (*dto.TicketResponse, error) {
	var result dto.TicketResponse
	err := r.DB.Raw(r.selectQuery("t.uuid"), uuid).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrTicketNotFound
	}
	return &result, nil
}

func (r *TicketRepo) findByID(id uint64) (*dto.TicketResponse, error) {
	var result dto.TicketResponse
	err := r.DB.Raw(r.selectQuery("t.id"), id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrTicketNotFound
	}
	return &result, nil
}

func (r *TicketRepo) GetTicketTypeIDByUUID(uuid string) (int64, error) {
	var id int64
	err := r.DB.Raw(`SELECT id FROM ticket_types WHERE uuid = ? AND status = 'active' AND deleted_at IS NULL`, uuid).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, domainerrors.ErrTicketTypeNotFound
	}
	return id, nil
}

func (r *TicketRepo) GetTicketInternalIDByUUID(uuid string) (int64, error) {
	var id int64
	err := r.DB.Raw(`SELECT id FROM tickets WHERE uuid = ? AND deleted_at IS NULL`, uuid).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, domainerrors.ErrTicketNotFound
	}
	return id, nil
}

func (r *TicketRepo) IsUserProjectMember(projectID int64, userID int64) (bool, error) {
	var count int64
	err := r.DB.Raw(
		`SELECT COUNT(*) FROM project_members WHERE project_id = ? AND user_id = ? AND deleted_at IS NULL`,
		projectID, userID,
	).Scan(&count).Error
	return count > 0, err
}

func (r *TicketRepo) GenerateCode(projectID int64) (string, error) {
	var projectCode string
	err := r.DB.Raw(`SELECT project_code FROM projects WHERE id = ? AND deleted_at IS NULL`, projectID).Scan(&projectCode).Error
	if err != nil || projectCode == "" {
		return "", err
	}

	var count int64
	err = r.DB.Raw(`SELECT COUNT(*) FROM tickets WHERE project_id = ?`, projectID).Scan(&count).Error
	if err != nil {
		return "", err
	}

	return projectCode + "-" + fmt.Sprintf("%d", count+1), nil
}

func (r *TicketRepo) selectQuery(whereCol string, extra ...string) string {
	q := `SELECT t.uuid AS id, p.id AS internal_project_id, p.uuid AS project_id, t.code, t.title, t.description, t.branch,
		tt.uuid AS ticket_type_id, tt.name AS ticket_type_name,
		ua.uuid AS assigned_to, ua.name AS assigned_to_name,
		ub.uuid AS assigned_by, ub.name AS assigned_by_name,
		t.start_date, t.internal_start_date, t.end_date, t.internal_end_date,
		t.estimated_hours, t.internal_estimated_hours,
		t.lifecycle_status, t.priority,
		pt.uuid AS parent_id,
		t.due_date, t.status, t.created_at, t.updated_at
		FROM tickets t
		INNER JOIN projects p ON p.id = t.project_id AND p.deleted_at IS NULL
		INNER JOIN ticket_types tt ON tt.id = t.ticket_type_id
		LEFT JOIN users ua ON ua.id = t.assigned_to AND ua.deleted_at IS NULL
		LEFT JOIN users ub ON ub.id = t.assigned_by AND ub.deleted_at IS NULL
		LEFT JOIN tickets pt ON pt.id = t.parent_id AND pt.deleted_at IS NULL
		WHERE t.deleted_at IS NULL`
	if whereCol != "" {
		q += ` AND ` + whereCol + ` = ?`
	}
	for _, e := range extra {
		q += ` ` + e
	}
	return q
}

func (r *TicketRepo) GetAllByProjectID(projectID int64, pagination string) (*[]dto.TicketResponse, error) {
	var results []dto.TicketResponse
	err := r.DB.Raw(r.selectQuery("t.project_id", "ORDER BY t.created_at DESC", pagination), projectID).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (r *TicketRepo) CountByProjectID(projectID int64) (int64, error) {
	var count int64
	err := r.DB.Raw(`SELECT COUNT(*) FROM tickets WHERE project_id = ? AND deleted_at IS NULL`, projectID).Scan(&count).Error
	return count, err
}

func (r *TicketRepo) Update(uuid string, updates map[string]any) (*dto.TicketResponse, error) {
	res := r.DB.Table("tickets").Where("uuid = ? AND deleted_at IS NULL", uuid).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrTicketNotFound
	}
	return r.FindByUUID(uuid)
}

func (r *TicketRepo) Delete(uuid string) error {
	res := r.DB.Exec(`UPDATE tickets SET deleted_at = NOW() WHERE uuid = ? AND deleted_at IS NULL`, uuid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrTicketNotFound
	}
	return nil
}

func (r *TicketRepo) IsUserAssignedToTicket(uuid string, id string) (int64, error) {
	var ticketID int64 = 0
	err := r.DB.Raw(`SELECT t.id FROM tickets t
               LEFT JOIN users us ON users.uuid = tickets.assigned_to 
              WHERE t.uuid = ? AND t.assigned_to= us.id `,
		uuid, id).Scan(&ticketID).Error

	return ticketID, err

}
