package services

import (
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	"net/http"
	"time"
)

type ProjectService struct {
	projectRepo *repositories.ProjectRepo
	clientRepo  *repositories.ClientRepo
	userRepo    *repositories.UserRepo
	redisRepo   *cacherepo.RedisRepo
}

func NewProjectService(
	projectRepo *repositories.ProjectRepo,
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
		StartDate:              parseDate(req.StartDate),
		InternalStartDate:      parseDate(req.InternalStartDate),
		EndDate:                parseDate(req.EndDate),
		InternalEndDate:        parseDate(req.InternalEndDate),
		EstimatedHours:         req.EstimatedHours,
		InternalEstimatedHours: req.InternalEstimatedHours,
		PrimaryRepoURL:         req.PrimaryRepoURL,
		BaseModel:              models.NewBase(createdBy),
	}

	return s.projectRepo.Create(project)
}

func parseDate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil
	}
	return &t
}
