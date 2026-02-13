package projectsvc

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	projectrepo "capecom-pm/internal/repositories/project"
	"capecom-pm/internal/utils"
	"net/http"
)

type ProjectService struct {
	projectRepo *projectrepo.ProjectRepo
	clientRepo  *repositories.ClientRepo
	userRepo    *repositories.UserRepo
	redisRepo   *cacherepo.RedisRepo
}

func NewProjectService(
	projectRepo *projectrepo.ProjectRepo,
	clientRepo *repositories.ClientRepo,
	userRepo *repositories.UserRepo,
	redisRepo *cacherepo.RedisRepo,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		clientRepo:  clientRepo,
		userRepo:    userRepo,
		redisRepo:   redisRepo,
	}
}

func (s *ProjectService) Create(req dto.CreateProjectRequest, userUUID string) (*dto.ProjectResponse, error) {
	var createdBy *uint64

	if userUUID != "" {
		userID, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
		if err != nil || userID == nil {
			return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "project_service", "get_user_id")
		}
		uid := uint64(*userID)
		createdBy = &uid
	}

	var clientID *uint64
	var clientNameSnapshot *string

	if req.ClientUUID != nil && *req.ClientUUID != "" {
		cID, err := s.clientRepo.GetInternalIDByUUID(*req.ClientUUID)
		if err != nil {
			return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrClientNotFound.Error(), "project_service", "resolve_client")
		}
		uid := uint64(cID)
		clientID = &uid

		clientResp, err := s.clientRepo.FindByUUID(*req.ClientUUID)
		if err == nil && clientResp != nil {
			clientNameSnapshot = &clientResp.Name
		}
	}

	lifecycleStatus := "todo"
	if req.LifecycleStatus != nil {
		lifecycleStatus = *req.LifecycleStatus
	}

	project := &models.Project{
		ProjectName:            req.ProjectName,
		ProjectCode:            req.ProjectCode,
		ClientID:               clientID,
		ClientNameSnapshot:     clientNameSnapshot,
		LifecycleStatus:        lifecycleStatus,
		StartDate:              utils.ParseDate(req.StartDate),
		InternalStartDate:      utils.ParseDate(req.InternalStartDate),
		EndDate:                utils.ParseDate(req.EndDate),
		InternalEndDate:        utils.ParseDate(req.InternalEndDate),
		EstimatedHours:         req.EstimatedHours,
		InternalEstimatedHours: req.InternalEstimatedHours,
		PrimaryRepoURL:         req.PrimaryRepoURL,
		BaseModel:              models.NewBase(createdBy),
	}

	return s.projectRepo.Create(project)
}

func (s *ProjectService) Update(uuid string, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	_, err := s.projectRepo.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)

	if req.ProjectName != nil {
		updates["project_name"] = *req.ProjectName
	}
	if req.ProjectCode != nil {
		updates["project_code"] = *req.ProjectCode
	}
	if req.LifecycleStatus != nil {
		updates["lifecycle_status"] = *req.LifecycleStatus
	}
	if req.StartDate != nil {
		updates["start_date"] = utils.ParseDate(req.StartDate)
	}
	if req.InternalStartDate != nil {
		updates["internal_start_date"] = utils.ParseDate(req.InternalStartDate)
	}
	if req.EndDate != nil {
		updates["end_date"] = utils.ParseDate(req.EndDate)
	}
	if req.InternalEndDate != nil {
		updates["internal_end_date"] = utils.ParseDate(req.InternalEndDate)
	}
	if req.EstimatedHours != nil {
		updates["estimated_hours"] = *req.EstimatedHours
	}
	if req.InternalEstimatedHours != nil {
		updates["internal_estimated_hours"] = *req.InternalEstimatedHours
	}
	if req.PrimaryRepoURL != nil {
		updates["primary_repo_url"] = *req.PrimaryRepoURL
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.ClientUUID != nil {
		if *req.ClientUUID == "" {
			updates["client_id"] = nil
			updates["client_name_snapshot"] = nil
		} else {
			cID, err := s.clientRepo.GetInternalIDByUUID(*req.ClientUUID)
			if err != nil {
				return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrClientNotFound.Error(), "project_service", "resolve_client")
			}
			updates["client_id"] = cID
			clientResp, err := s.clientRepo.FindByUUID(*req.ClientUUID)
			if err == nil && clientResp != nil {
				updates["client_name_snapshot"] = clientResp.Name
			}
		}
	}

	if len(updates) == 0 {
		return s.projectRepo.FindByUUID(uuid)
	}

	return s.projectRepo.Update(uuid, updates)
}

func (s *ProjectService) UpdateLifecycleStatus(uuid string, status string) (*dto.ProjectResponse, error) {
	return s.projectRepo.Update(uuid, map[string]any{"lifecycle_status": status})
}

func (s *ProjectService) Delete(uuid string) error {
	return s.projectRepo.Delete(uuid)
}

func (s *ProjectService) GetProjects(pg *common.Pagination) (*[]dto.ProjectResponse, error) {
	projects, err := s.projectRepo.GetAll(pg)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
