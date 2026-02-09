package mastersreo

import (
	domainerrors "capecom-pm/internal/domain/error"

	"gorm.io/gorm"
)

type MasterDataRepo struct {
	DB *gorm.DB
}

func NewMasterDataRepo(db *gorm.DB) *MasterDataRepo {
	return &MasterDataRepo{
		DB: db,
	}
}

type MasterDataIDs struct {
	GroupID       int64
	DesignationID int64
	DepartmentID  int64
	RoleIDs       []int64
}

// GetUserRelatedIDs fetches all related IDs in a single optimized query
func (r *MasterDataRepo) GetUserRelatedIDs(
	groupUUID string,
	designationUUID string,
	departmentUUID string,
	roleUUIDs []string,
) (*MasterDataIDs, error) {
	result := &MasterDataIDs{}

	// Single query using UNION ALL to fetch all IDs at once
	query := `
		SELECT 'group' as type, id, uuid FROM user_groups WHERE uuid = ? AND status = 'active'
		UNION ALL
		SELECT 'designation' as type, id, uuid FROM designations WHERE uuid = ? AND status = 'active'
		UNION ALL
		SELECT 'department' as type, id, uuid FROM departments WHERE uuid = ? AND status = 'active'
		UNION ALL
		SELECT 'role' as type, id, uuid FROM roles WHERE uuid IN ? AND status = 'active'
	`

	type row struct {
		Type string
		ID   int64
		UUID string
	}

	var rows []row
	err := r.DB.Raw(query, groupUUID, designationUUID, departmentUUID, roleUUIDs).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// Map results
	foundGroup := false
	foundDesignation := false
	foundDepartment := false

	for _, r := range rows {
		switch r.Type {
		case "group":
			result.GroupID = r.ID
			foundGroup = true
		case "designation":
			result.DesignationID = r.ID
			foundDesignation = true
		case "department":
			result.DepartmentID = r.ID
			foundDepartment = true
		case "role":
			result.RoleIDs = append(result.RoleIDs, r.ID)
		}
	}

	// Validate all required data was found
	if !foundGroup {
		return nil, domainerrors.ErrGroupNotFound
	}
	if !foundDesignation {
		return nil, domainerrors.ErrDesignationNotFound
	}
	if !foundDepartment {
		return nil, domainerrors.ErrDepartmentNotFound
	}
	if len(result.RoleIDs) != len(roleUUIDs) {
		return nil, domainerrors.ErrInvalidRoles
	}

	return result, nil
}
