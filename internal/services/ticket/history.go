package ticketsvc

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	ticketrepo "capecom-pm/internal/repositories/ticket"
	"net/http"
)

type HistoryService struct {
	historyRepo *ticketrepo.HistoryRepo
	ticketRepo  *ticketrepo.TicketRepo
}

func NewHistoryService(
	historyRepo *ticketrepo.HistoryRepo,
	ticketRepo *ticketrepo.TicketRepo,
) *HistoryService {
	return &HistoryService{
		historyRepo: historyRepo,
		ticketRepo:  ticketRepo,
	}
}

func (s *HistoryService) GetAllByTicket(ticketUUID string, pg *common.Pagination) (*dto.ListWithMeta, error) {
	ticketID, err := s.ticketRepo.GetTicketInternalIDByUUID(ticketUUID)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrTicketNotFound.Error(), "history_service", "resolve_ticket")
	}

	total, err := s.historyRepo.CountByTicketID(ticketID)
	if err != nil {
		return nil, err
	}

	history, err := s.historyRepo.GetAllByTicketID(ticketID, pg.BuildPaginationQuery())
	if err != nil {
		return nil, err
	}

	return &dto.ListWithMeta{
		Data: history,
		Meta: dto.PaginationMeta{
			Page:    pg.Page,
			Limit:   pg.Limit,
			Total:   total,
			HasMore: pg.HasMore(total),
		},
	}, nil
}
