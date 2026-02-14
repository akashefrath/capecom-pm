package ticketsvc

import (
	"capecom-pm/internal/config"
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	ticketrepo "capecom-pm/internal/repositories/ticket"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type TimeEntryService struct {
	timeEntryRepo *ticketrepo.TimeEntryRepo
	ticketRepo    *ticketrepo.TicketRepo
	historyRepo   *ticketrepo.HistoryRepo
	userRepo      *repositories.UserRepo
	redisRepo     *cacherepo.RedisRepo
}

func NewTimeEntryService(
	timeEntryRepo *ticketrepo.TimeEntryRepo,
	ticketRepo *ticketrepo.TicketRepo,
	historyRepo *ticketrepo.HistoryRepo,
	userRepo *repositories.UserRepo,
	redisRepo *cacherepo.RedisRepo,
) *TimeEntryService {
	return &TimeEntryService{
		timeEntryRepo: timeEntryRepo,
		ticketRepo:    ticketRepo,
		historyRepo:   historyRepo,
		userRepo:      userRepo,
		redisRepo:     redisRepo,
	}
}

func (s *TimeEntryService) Create(ticketUUID string, req dto.CreateTimeEntryRequest, userUUID string) (*dto.TimeEntryResponse, error) {
	// resolve ticket
	ticketID, err := s.ticketRepo.GetTicketInternalIDByUUID(ticketUUID)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrTicketNotFound.Error(), "time_entry_service", "resolve_ticket")
	}

	// resolve user
	var userID *uint64
	var createdBy *uint64
	if userUUID != "" {
		uid, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
		if err != nil || uid == nil {
			return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "time_entry_service", "get_user_id")
		}
		u := uint64(*uid)
		userID = &u
		createdBy = &u
	}

	// verify user is assigned to ticket
	assignedUserID, err := s.timeEntryRepo.GetTicketAssignedUserID(ticketID)

	if err != nil {
		return nil, err
	}

	if assignedUserID == nil || userID == nil || *assignedUserID != int64(*userID) {
		return nil, domainerrors.NewWithCode(http.StatusForbidden, domainerrors.ErrNotTicketAssignee.Error(), "time_entry_service", "check_assignee")
	}
	workDate, _ := time.Parse("2006-01-02", req.WorkDate)
	date, err := s.timeEntryRepo.GetTotalHoursByUserIDByDate(int64(*userID), workDate)
	maxAllowedHours, _ := strconv.ParseFloat(config.GetEnvMust("MAX_BOOK_HOURS_PER_DAY"), 64)
	if err != nil {
		// Fallback to a default if the env is messed up
		maxAllowedHours = 8.0
	}
	if date != nil && ((*date + req.Hours) > maxAllowedHours) {
		return nil, domainerrors.NewWithCode(http.StatusForbidden, domainerrors.ErrHoursExceeded.Error(), "time_entry_service", "check_hours")
	}

	// get project_id from ticket
	ticket, err := s.ticketRepo.FindByUUID(ticketUUID)

	if err != nil {
		return nil, err
	}

	entry := &models.TimeEntry{
		TicketID:    uint64(ticketID),
		ProjectID:   ticket.InternalProjectID,
		UserID:      *userID,
		WorkDate:    workDate,
		Hours:       req.Hours,
		Description: req.Description,
		BaseModel:   models.NewBase(createdBy),
	}

	result, err := s.timeEntryRepo.Create(entry)

	if err != nil {
		return nil, err
	}

	// log history
	s.logHistory(ticketID, *userID, "time_entry_added", nil, &req.WorkDate, nil)

	return result, nil
}

func (s *TimeEntryService) GetByUUID(entryUUID string) (*dto.TimeEntryResponse, error) {
	return s.timeEntryRepo.FindByUUID(entryUUID)
}

func (s *TimeEntryService) GetAllByTicket(ticketUUID string, pg *common.Pagination) (*dto.ListWithMeta, error) {
	ticketID, err := s.ticketRepo.GetTicketInternalIDByUUID(ticketUUID)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrTicketNotFound.Error(), "time_entry_service", "resolve_ticket")
	}

	total, err := s.timeEntryRepo.CountByTicketID(ticketID)
	if err != nil {
		return nil, err
	}

	entries, err := s.timeEntryRepo.GetAllByTicketID(ticketID, pg.BuildPaginationQuery())
	if err != nil {
		return nil, err
	}

	return &dto.ListWithMeta{
		Data: entries,
		Meta: dto.PaginationMeta{
			Page:    pg.Page,
			Limit:   pg.Limit,
			Total:   total,
			HasMore: pg.HasMore(total),
		},
	}, nil
}

func (s *TimeEntryService) Update(entryUUID string, req dto.UpdateTimeEntryRequest, userUUID string) (*dto.TimeEntryResponse, error) {
	// resolve user
	var userID *int64
	if userUUID != "" {
		uid, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
		if err != nil || uid == nil {
			return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "time_entry_service", "get_user_id")
		}
		userID = uid
	}

	// verify user owns this time entry
	ownerID, err := s.timeEntryRepo.GetTimeEntryOwnerID(entryUUID)
	if err != nil {
		return nil, err
	}
	if ownerID == nil || userID == nil || *ownerID != *userID {
		return nil, domainerrors.NewWithCode(http.StatusForbidden, domainerrors.ErrNotTicketAssignee.Error(), "time_entry_service", "check_owner")
	}

	updates := make(map[string]any)

	if req.WorkDate != nil {
		workDate, _ := time.Parse("2006-01-02", *req.WorkDate)
		updates["work_date"] = workDate
	}
	if req.Hours != nil {
		updates["hours"] = *req.Hours
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if len(updates) == 0 {
		return s.timeEntryRepo.FindByUUID(entryUUID)
	}

	return s.timeEntryRepo.Update(entryUUID, updates)
}

func (s *TimeEntryService) Delete(entryUUID string, userUUID string) error {
	// resolve user
	var userID *int64
	if userUUID != "" {
		uid, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
		if err != nil || uid == nil {
			return domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "time_entry_service", "get_user_id")
		}
		userID = uid
	}

	// verify user owns this time entry
	ownerID, err := s.timeEntryRepo.GetTimeEntryOwnerID(entryUUID)
	if err != nil {
		return err
	}
	if ownerID == nil || userID == nil || *ownerID != *userID {
		return domainerrors.NewWithCode(http.StatusForbidden, domainerrors.ErrNotTicketAssignee.Error(), "time_entry_service", "check_owner")
	}

	return s.timeEntryRepo.Delete(entryUUID)
}

func (s *TimeEntryService) logHistory(ticketID int64, userID uint64, fieldName string, oldValue, newValue, note *string) {
	history := &models.TicketHistory{
		UUID:      uuid.New().String(),
		TicketID:  uint64(ticketID),
		ChangedBy: userID,
		FieldName: fieldName,
		OldValue:  oldValue,
		NewValue:  newValue,
		Note:      note,
	}
	err := s.historyRepo.Create(history)
	if err != nil {
		return
	}
}
