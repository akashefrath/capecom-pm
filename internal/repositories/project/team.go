package projectrepo

import (
	"capecom-pm/internal/domain/dto"

	"gorm.io/gorm"
)

type TeamRepo struct {
	DB *gorm.DB
}

func NewTeamRepo(db *gorm.DB) *TeamRepo {
	return &TeamRepo{DB: db}
}

// --- shared helpers ---

type existingRow struct {
	UserID    int64
	DeletedAt *string
}

// getExistingRows returns all rows (active + soft-deleted) for given project+users
func (r *TeamRepo) getExistingRows(table string, projectID int64, userIDs []int64) (active map[int64]bool, softDeleted map[int64]bool, err error) {
	var rows []existingRow
	err = r.DB.Raw(
		`SELECT user_id, deleted_at FROM `+table+` WHERE project_id = ? AND user_id IN ?`,
		projectID, userIDs,
	).Scan(&rows).Error
	if err != nil {
		return nil, nil, err
	}

	active = make(map[int64]bool)
	softDeleted = make(map[int64]bool)
	for _, row := range rows {
		if row.DeletedAt == nil {
			active[row.UserID] = true
		} else {
			softDeleted[row.UserID] = true
		}
	}
	return active, softDeleted, nil
}

// reactivateRows sets deleted_at=NULL and status=active for soft-deleted rows
func (r *TeamRepo) reactivateRows(table string, projectID int64, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	return r.DB.Exec(
		`UPDATE `+table+` SET deleted_at = NULL, status = 'active', updated_at = NOW() WHERE project_id = ? AND user_id IN ? AND deleted_at IS NOT NULL`,
		projectID, userIDs,
	).Error
}

// --- Managers ---

func (r *TeamRepo) BulkInsertManagers(rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}
	return r.DB.Table("project_managers").Create(rows).Error
}

func (r *TeamRepo) ReactivateManagers(projectID int64, userIDs []int64) error {
	return r.reactivateRows("project_managers", projectID, userIDs)
}

func (r *TeamRepo) GetExistingManagerRows(projectID int64, userIDs []int64) (active map[int64]bool, softDeleted map[int64]bool, err error) {
	return r.getExistingRows("project_managers", projectID, userIDs)
}

func (r *TeamRepo) BulkSoftDeleteManagers(projectID int64, userIDs []int64) (int64, error) {
	res := r.DB.Exec(
		`UPDATE project_managers SET deleted_at = NOW() WHERE project_id = ? AND user_id IN ? AND deleted_at IS NULL`,
		projectID, userIDs,
	)
	return res.RowsAffected, res.Error
}

func (r *TeamRepo) GetManagersByProjectID(projectID int64) ([]dto.ProjectManagerResponse, error) {
	var result []dto.ProjectManagerResponse
	err := r.DB.Raw(`
		SELECT pm.uuid AS id, u.uuid AS user_id, u.name AS user_name,
			pm.status, pm.created_at, pm.updated_at
		FROM project_managers pm
		INNER JOIN users u ON u.id = pm.user_id AND u.deleted_at IS NULL
		WHERE pm.project_id = ? AND pm.deleted_at IS NULL
		ORDER BY u.name ASC
	`, projectID).Scan(&result).Error
	if result == nil {
		result = make([]dto.ProjectManagerResponse, 0)
	}
	return result, err
}

// --- Members ---

func (r *TeamRepo) BulkInsertMembers(rows []map[string]any) error {
	if len(rows) == 0 {
		return nil
	}
	return r.DB.Table("project_members").Create(rows).Error
}

func (r *TeamRepo) ReactivateMembers(projectID int64, userIDs []int64, allocMap map[int64]float64) error {
	if len(userIDs) == 0 {
		return nil
	}
	// reactivate + update allocated_hours per user
	for _, uid := range userIDs {
		hours := allocMap[uid]
		if err := r.DB.Exec(
			`UPDATE project_members SET deleted_at = NULL, status = 'active', allocated_hours = ?, updated_at = NOW() WHERE project_id = ? AND user_id = ? AND deleted_at IS NOT NULL`,
			hours, projectID, uid,
		).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *TeamRepo) GetExistingMemberRows(projectID int64, userIDs []int64) (active map[int64]bool, softDeleted map[int64]bool, err error) {
	return r.getExistingRows("project_members", projectID, userIDs)
}

func (r *TeamRepo) BulkSoftDeleteMembers(projectID int64, userIDs []int64) (int64, error) {
	res := r.DB.Exec(
		`UPDATE project_members SET deleted_at = NOW() WHERE project_id = ? AND user_id IN ? AND deleted_at IS NULL`,
		projectID, userIDs,
	)
	return res.RowsAffected, res.Error
}

func (r *TeamRepo) GetMembersByProjectID(projectID int64) ([]dto.ProjectMemberResponse, error) {
	var result []dto.ProjectMemberResponse
	err := r.DB.Raw(`
		SELECT pm.uuid AS id, u.uuid AS user_id, u.name AS user_name,
			pm.allocated_hours, pm.status, pm.created_at, pm.updated_at
		FROM project_members pm
		INNER JOIN users u ON u.id = pm.user_id AND u.deleted_at IS NULL
		WHERE pm.project_id = ? AND pm.deleted_at IS NULL
		ORDER BY u.name ASC
	`, projectID).Scan(&result).Error
	if result == nil {
		result = make([]dto.ProjectMemberResponse, 0)
	}
	return result, err
}
