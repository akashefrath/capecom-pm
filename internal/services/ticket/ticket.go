package ticketsvc

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	projectrepo "capecom-pm/internal/repositories/project"
	ticketrepo "capecom-pm/internal/repositories/ticket"
	"capecom-pm/internal/utils"
	"net/http"
)

type TicketService struct {
	ticketRepo  *ticketrepo.TicketRepo
	projectRepo *projectrepo.ProjectRepo
	userRepo    *repositories.UserRepo
	redisRepo   *cacherepo.RedisRepo
}

func NewTicketService(
	ticketRepo *ticketrepo.TicketRepo,
	projectRepo *projectrepo.ProjectRepo,
	userRepo *repositories.UserRepo,
	redisRepo *cacherepo.RedisRepo,
) *TicketService {
	return &TicketService{
		ticketRepo:  ticketRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
		redisRepo:   redisRepo,
	}
}

func (s *TicketService) Create(projectUUID string, req dto.CreateTicketRequest, userUUID string) (*dto.TicketResponse, error) {
	// resolve project
	projectID, err := s.projectRepo.GetInternalIDByUUID(projectUUID)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrProjectNotFound.Error(), "ticket_service", "resolve_project")
	}

	// resolve created_by (the person creating = assigned_by)
	var createdBy *uint64
	var assignedBy *uint64
	if userUUID != "" {
		uid, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
		if err != nil || uid == nil {
			return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "ticket_service", "get_user_id")
		}
		u := uint64(*uid)
		createdBy = &u
		assignedBy = &u
	}

	// resolve ticket_type
	ticketTypeID, err := s.ticketRepo.GetTicketTypeIDByUUID(req.TicketTypeUUID)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrTicketTypeNotFound.Error(), "ticket_service", "resolve_ticket_type")
	}

	// resolve assigned_to — must be a project member
	var assignedTo *uint64
	if req.AssignedToUUID != nil && *req.AssignedToUUID != "" {
		idMap, err := s.userRepo.GetActiveUserIDsByUuids([]string{*req.AssignedToUUID})
		if err != nil || len(idMap) == 0 {
			return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrUserNotFound.Error(), "ticket_service", "resolve_assigned_to")
		}
		uid := idMap[*req.AssignedToUUID]

		isMember, err := s.ticketRepo.IsUserProjectMember(projectID, uid)
		if err != nil || !isMember {
			return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrNotProjectMember.Error(), "ticket_service", "check_membership")
		}
		u := uint64(uid)
		assignedTo = &u
	}

	// resolve parent ticket
	var parentID *uint64
	if req.ParentUUID != nil && *req.ParentUUID != "" {
		pID, err := s.ticketRepo.GetTicketInternalIDByUUID(*req.ParentUUID)
		if err != nil {
			return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrTicketNotFound.Error(), "ticket_service", "resolve_parent")
		}
		u := uint64(pID)
		parentID = &u
	}

	// generate ticket code
	code, err := s.ticketRepo.GenerateCode(projectID)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusInternalServerError, domainerrors.ErrInternal.Error(), "ticket_service", "generate_code")
	}

	priority := "medium"
	if req.Priority != nil {
		priority = *req.Priority
	}

	ticket := &models.Ticket{
		ProjectID:              uint64(projectID),
		Code:                   code,
		Title:                  req.Title,
		Description:            req.Description,
		Branch:                 req.Branch,
		TicketTypeID:           uint64(ticketTypeID),
		AssignedTo:             assignedTo,
		AssignedBy:             assignedBy,
		StartDate:              utils.ParseDate(req.StartDate),
		InternalStartDate:      utils.ParseDate(req.InternalStartDate),
		EndDate:                utils.ParseDate(req.EndDate),
		InternalEndDate:        utils.ParseDate(req.InternalEndDate),
		EstimatedHours:         req.EstimatedHours,
		InternalEstimatedHours: req.InternalEstimatedHours,
		LifecycleStatus:        "todo",
		Priority:               priority,
		ParentID:               parentID,
		DueDate:                utils.ParseDate(req.DueDate),
		BaseModel:              models.NewBase(createdBy),
	}

	return s.ticketRepo.Create(ticket)
}

func (s *TicketService) GetByUUID(ticketUUID string) (*dto.TicketResponse, error) {
	return s.ticketRepo.FindByUUID(ticketUUID)
}

func (s *TicketService) GetAllByProject(projectUUID string, pg *common.Pagination, userID string) (*dto.ListWithMeta, error) {
	projectID, err := s.projectRepo.GetInternalIDByUUID(projectUUID)

	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrProjectNotFound.Error(), "ticket_service", "resolve_project")
	}

	userId, err := s.redisRepo.GetUserIdByUuid(userID, *s.userRepo)

	if err != nil || userId == nil {
		return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "ticket_service", "resolve_project")
	}

	member, err := s.projectRepo.IsUserProjectMember(*userId, projectUUID)
	isAdmin, err := s.userRepo.IsManagerOrAdmin(userID)

	if err != nil || (!member && !isAdmin) {
		return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "ticket_service", "resolve_project")
	}

	total, err := s.ticketRepo.CountByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	tickets, err := s.ticketRepo.GetAllByProjectID(projectID, pg.BuildPaginationQuery())
	if err != nil {
		return nil, err
	}

	return &dto.ListWithMeta{
		Data: tickets,
		Meta: dto.PaginationMeta{
			Page:    pg.Page,
			Limit:   pg.Limit,
			Total:   total,
			HasMore: pg.HasMore(total),
		},
	}, nil
}

func (s *TicketService) Update(ticketUUID string, req dto.UpdateTicketRequest) (*dto.TicketResponse, error) {
	_, err := s.ticketRepo.FindByUUID(ticketUUID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Branch != nil {
		updates["branch"] = *req.Branch
	}
	if req.TicketTypeUUID != nil {
		typeID, err := s.ticketRepo.GetTicketTypeIDByUUID(*req.TicketTypeUUID)
		if err != nil {
			return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrTicketTypeNotFound.Error(), "ticket_service", "resolve_ticket_type")
		}
		updates["ticket_type_id"] = typeID
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
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.ParentUUID != nil {
		if *req.ParentUUID == "" {
			updates["parent_id"] = nil
		} else {
			pID, err := s.ticketRepo.GetTicketInternalIDByUUID(*req.ParentUUID)
			if err != nil {
				return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrTicketNotFound.Error(), "ticket_service", "resolve_parent")
			}
			updates["parent_id"] = pID
		}
	}
	if req.DueDate != nil {
		updates["due_date"] = utils.ParseDate(req.DueDate)
	}

	if len(updates) == 0 {
		return s.ticketRepo.FindByUUID(ticketUUID)
	}

	return s.ticketRepo.Update(ticketUUID, updates)
}

func (s *TicketService) UpdateLifecycleStatus(ticketUUID string, status string, userID string) (*dto.TicketResponse, error) {
	ticketID, err := s.ticketRepo.IsUserAssignedToTicket(ticketUUID, userID)
	isAdmin, err := s.userRepo.IsManagerOrAdmin(userID)
	if err != nil || (!isAdmin && ticketID == 0) {
		return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "ticket_service", "resolve_project")
	}
	return s.ticketRepo.Update(ticketUUID, map[string]any{"lifecycle_status": status})
}

func (s *TicketService) UpdateAssignee(ticketUUID string, projectUUID string, assignedToUUID string) (*dto.TicketResponse, error) {
	projectID, err := s.projectRepo.GetInternalIDByUUID(projectUUID)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrProjectNotFound.Error(), "ticket_service", "resolve_project")
	}

	idMap, err := s.userRepo.GetActiveUserIDsByUuids([]string{assignedToUUID})
	if err != nil || len(idMap) == 0 {
		return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrUserNotFound.Error(), "ticket_service", "resolve_assignee")
	}
	uid := idMap[assignedToUUID]

	isMember, err := s.ticketRepo.IsUserProjectMember(projectID, uid)
	if err != nil || !isMember {
		return nil, domainerrors.NewWithCode(http.StatusBadRequest, domainerrors.ErrNotProjectMember.Error(), "ticket_service", "check_membership")
	}

	return s.ticketRepo.Update(ticketUUID, map[string]any{"assigned_to": uid})
}

func (s *TicketService) Delete(ticketUUID string) error {
	return s.ticketRepo.Delete(ticketUUID)
}
