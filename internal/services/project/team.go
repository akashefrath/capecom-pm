package projectsvc

import (
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	projectrepo "capecom-pm/internal/repositories/project"
	"net/http"

	"github.com/google/uuid"
)

type TeamService struct {
	teamRepo    *projectrepo.TeamRepo
	projectRepo *projectrepo.ProjectRepo
	userRepo    *repositories.UserRepo
	redisRepo   *cacherepo.RedisRepo
}

func NewTeamService(
	teamRepo *projectrepo.TeamRepo,
	projectRepo *projectrepo.ProjectRepo,
	userRepo *repositories.UserRepo,
	redisRepo *cacherepo.RedisRepo,
) *TeamService {
	return &TeamService{
		teamRepo:    teamRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
		redisRepo:   redisRepo,
	}
}

// --- helpers ---

func (s *TeamService) resolveProjectID(projectUUID string) (int64, error) {
	id, err := s.projectRepo.GetInternalIDByUUID(projectUUID)
	if err != nil {
		return 0, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrProjectNotFound.Error(), "team_service", "resolve_project")
	}
	return id, nil
}

func (s *TeamService) resolveUserIDs(uuids []string) (map[string]int64, error) {
	idMap, err := s.userRepo.GetActiveUserIDsByUuids(uuids)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusInternalServerError, domainerrors.ErrInternal.Error(), "team_service", "resolve_users")
	}
	if len(idMap) != len(uuids) {
		return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrUserNotFound.Error(), "team_service", "resolve_users")
	}
	return idMap, nil
}

func (s *TeamService) resolveCreatedBy(userUUID string) (*uint64, error) {
	if userUUID == "" {
		return nil, nil
	}
	userID, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
	if err != nil || userID == nil {
		return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "team_service", "get_user_id")
	}
	uid := uint64(*userID)
	return &uid, nil
}

func collectIDs(idMap map[string]int64) []int64 {
	ids := make([]int64, 0, len(idMap))
	for _, id := range idMap {
		ids = append(ids, id)
	}
	return ids
}

// --- Managers ---

func (s *TeamService) AssignManagers(projectUUID string, req dto.AssignManagersRequest, userUUID string) ([]dto.ProjectManagerResponse, error) {
	projectID, err := s.resolveProjectID(projectUUID)
	if err != nil {
		return nil, err
	}

	idMap, err := s.resolveUserIDs(req.UserUUIDs)
	if err != nil {
		return nil, err
	}

	createdBy, err := s.resolveCreatedBy(userUUID)
	if err != nil {
		return nil, err
	}

	allIDs := collectIDs(idMap)
	active, softDeleted, err := s.teamRepo.GetExistingManagerRows(projectID, allIDs)
	if err != nil {
		return nil, err
	}

	// reactivate soft-deleted ones
	var toReactivate []int64
	for _, uid := range allIDs {
		if softDeleted[uid] {
			toReactivate = append(toReactivate, uid)
		}
	}
	if err := s.teamRepo.ReactivateManagers(projectID, toReactivate); err != nil {
		return nil, err
	}

	// insert only truly new ones
	var rows []map[string]any
	for _, uid := range idMap {
		if active[uid] || softDeleted[uid] {
			continue
		}
		row := map[string]any{
			"uuid":       uuid.NewString(),
			"user_id":    uid,
			"project_id": projectID,
			"status":     "active",
		}
		if createdBy != nil {
			row["created_by"] = *createdBy
		}
		rows = append(rows, row)
	}

	if err := s.teamRepo.BulkInsertManagers(rows); err != nil {
		return nil, err
	}

	return s.teamRepo.GetManagersByProjectID(projectID)
}

func (s *TeamService) RemoveManagers(projectUUID string, req dto.AssignManagersRequest) ([]dto.ProjectManagerResponse, error) {
	projectID, err := s.resolveProjectID(projectUUID)
	if err != nil {
		return nil, err
	}

	idMap, err := s.resolveUserIDs(req.UserUUIDs)
	if err != nil {
		return nil, err
	}

	if _, err := s.teamRepo.BulkSoftDeleteManagers(projectID, collectIDs(idMap)); err != nil {
		return nil, err
	}

	return s.teamRepo.GetManagersByProjectID(projectID)
}

func (s *TeamService) GetManagers(projectUUID string) ([]dto.ProjectManagerResponse, error) {
	projectID, err := s.resolveProjectID(projectUUID)
	if err != nil {
		return nil, err
	}
	return s.teamRepo.GetManagersByProjectID(projectID)
}

// --- Members ---

func (s *TeamService) AssignMembers(projectUUID string, req dto.AssignMembersRequest, userUUID string) ([]dto.ProjectMemberResponse, error) {
	projectID, err := s.resolveProjectID(projectUUID)
	if err != nil {
		return nil, err
	}

	uuids := make([]string, len(req.Members))
	for i, m := range req.Members {
		uuids[i] = m.UserUUID
	}

	idMap, err := s.resolveUserIDs(uuids)
	if err != nil {
		return nil, err
	}

	createdBy, err := s.resolveCreatedBy(userUUID)
	if err != nil {
		return nil, err
	}

	allIDs := collectIDs(idMap)
	active, softDeleted, err := s.teamRepo.GetExistingMemberRows(projectID, allIDs)
	if err != nil {
		return nil, err
	}

	// build allocated_hours map for reactivation
	allocMap := make(map[int64]float64)
	for _, m := range req.Members {
		allocMap[idMap[m.UserUUID]] = m.AllocatedHours
	}

	// reactivate soft-deleted ones with updated hours
	var toReactivate []int64
	for _, uid := range allIDs {
		if softDeleted[uid] {
			toReactivate = append(toReactivate, uid)
		}
	}
	if err := s.teamRepo.ReactivateMembers(projectID, toReactivate, allocMap); err != nil {
		return nil, err
	}

	// insert only truly new ones
	var rows []map[string]any
	for _, m := range req.Members {
		uid := idMap[m.UserUUID]
		if active[uid] || softDeleted[uid] {
			continue
		}
		row := map[string]any{
			"uuid":            uuid.NewString(),
			"user_id":         uid,
			"project_id":      projectID,
			"allocated_hours": m.AllocatedHours,
			"status":          "active",
		}
		if createdBy != nil {
			row["created_by"] = *createdBy
		}
		rows = append(rows, row)
	}

	if err := s.teamRepo.BulkInsertMembers(rows); err != nil {
		return nil, err
	}

	return s.teamRepo.GetMembersByProjectID(projectID)
}

func (s *TeamService) RemoveMembers(projectUUID string, req dto.RemoveMembersRequest) ([]dto.ProjectMemberResponse, error) {
	projectID, err := s.resolveProjectID(projectUUID)
	if err != nil {
		return nil, err
	}

	idMap, err := s.resolveUserIDs(req.UserUUIDs)
	if err != nil {
		return nil, err
	}

	if _, err := s.teamRepo.BulkSoftDeleteMembers(projectID, collectIDs(idMap)); err != nil {
		return nil, err
	}

	return s.teamRepo.GetMembersByProjectID(projectID)
}

func (s *TeamService) GetMembers(projectUUID string) ([]dto.ProjectMemberResponse, error) {
	projectID, err := s.resolveProjectID(projectUUID)
	if err != nil {
		return nil, err
	}
	return s.teamRepo.GetMembersByProjectID(projectID)
}
