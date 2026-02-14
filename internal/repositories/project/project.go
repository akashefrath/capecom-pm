package projectrepo

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	DB *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{DB: db}
}

func (r *ProjectRepo) Create(project *models.Project) (*dto.ProjectResponse, error) {
	var mysqlErr *mysql.MySQLError
	if err := r.DB.Create(project).Error; err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, domainerrors.NewWithCode(http.StatusConflict, domainerrors.ErrDuplicateProject.Error(), "repo", "CreateProject")
		}
		return nil, err
	}
	return r.findByID(project.ID)
}

func (r *ProjectRepo) FindByUUID(uuid string) (*dto.ProjectResponse, error) {
	var result dto.ProjectResponse

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Raw(r.selectQuery("p.uuid", "", ",COUNT(t.id) AS ticket_count", " LEFT JOIN tickets t on t.project_id = p.id"), uuid).Scan(&result).Error
		println(result.InternalID)
		err = tx.Raw(`SELECT p.uuid AS id, p.title , p.asset_type ,p.description ,p.content,
       f.uuid as file_id, f.storage_key as file_storage_key, f.mime_type as file_mime_type, f.size_bytes as files_size_bytes, f.storage as storage 
       FROM project_assets as p 
       LEFT JOIN files f ON f.id = p.file_id 
       WHERE p.project_id = ? AND p.deleted_at IS NULL`,
			result.InternalID).Scan(&result.Assets).Error

		return err
	})

	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrProjectNotFound
	}
	return &result, nil
}

func (r *ProjectRepo) findByID(id uint64) (*dto.ProjectResponse, error) {
	var result dto.ProjectResponse
	err := r.DB.Raw(r.selectQuery("p.id", "", "", ""), id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.Id == "" {
		return nil, domainerrors.ErrProjectNotFound
	}
	return &result, nil
}

func (r *ProjectRepo) GetAll(pg *common.Pagination) (*[]dto.ProjectResponse, error) {
	var results []dto.ProjectResponse
	err := r.DB.Raw(r.selectQuery("", pg.BuildPaginationQuery(), "", "")).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (r *ProjectRepo) GetInternalIDByUUID(uuid string) (int64, error) {
	var id int64
	err := r.DB.Raw(`SELECT id FROM projects WHERE uuid = ? AND deleted_at IS NULL`, uuid).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, domainerrors.ErrProjectNotFound
	}
	return id, nil
}

func (r *ProjectRepo) Update(uuid string, updates map[string]any) (*dto.ProjectResponse, error) {
	res := r.DB.Table("projects").Where("uuid = ? AND deleted_at IS NULL", uuid).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrProjectNotFound
	}
	return r.FindByUUID(uuid)
}

func (r *ProjectRepo) Delete(uuid string) error {
	res := r.DB.Exec(`UPDATE projects SET deleted_at = NOW() WHERE uuid = ? AND deleted_at IS NULL`, uuid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrProjectNotFound
	}
	return nil
}

func (r *ProjectRepo) selectQuery(whereCol string, pagination string, extraSelect string, extraJoin string) string {
	q := `SELECT p.id AS internal_id, p.uuid AS id, p.project_name, p.project_code,
		c.uuid AS client_id, p.client_name_snapshot AS client_name,
		p.lifecycle_status, p.start_date, p.internal_start_date,
		p.end_date, p.internal_end_date, p.estimated_hours,
		p.internal_estimated_hours, p.primary_repo_url,
		p.status, p.created_at, p.updated_at
		%s
		FROM projects p
		LEFT JOIN clients c ON c.id = p.client_id AND c.deleted_at IS NULL
		%s
		WHERE p.deleted_at IS NULL`
	if whereCol != "" {
		q += ` AND ` + whereCol + ` = ?`
	}
	if pagination != "" {
		q += ` ` + pagination
	}
	finalQ := q
	if extraSelect != "" && extraJoin != "" {
		finalQ = fmt.Sprintf(q, extraSelect, extraJoin)

	} else {
		finalQ = fmt.Sprintf(q, "", "")
	}

	return finalQ
}
