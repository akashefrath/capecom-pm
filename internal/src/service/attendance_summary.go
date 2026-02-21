package service

import (
	"github.com/akashefrath/capecom-pm/internal/domain/common"
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	utilsdto "github.com/akashefrath/capecom-pm/internal/domain/dto/utils"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
)

type AttendanceSummary struct {
	AttendanceSummaryRepo *repository.AttendanceSummary
}

func NewAttendanceSummary(attendanceSummaryRepo *repository.AttendanceSummary) *AttendanceSummary {
	return &AttendanceSummary{AttendanceSummaryRepo: attendanceSummaryRepo}
}

func (s *AttendanceSummary) GetByUUID(uuid string) (*dto.AttendanceSummaryResponse, error) {
	return s.AttendanceSummaryRepo.GetByUUID(uuid)
}

func (s *AttendanceSummary) GetList(pagination common.Pagination, query dto.AttendanceSummaryListQuery) (*utilsdto.ListWithMeta, error) {
	results, err := s.AttendanceSummaryRepo.GetList(pagination, query)
	if err != nil {
		return nil, err
	}

	count, err := s.AttendanceSummaryRepo.GetCount(query)
	if err != nil {
		return nil, err
	}

	return &utilsdto.ListWithMeta{
		Data: results,
		Meta: utilsdto.PaginationMeta{
			Page:    pagination.Page,
			Limit:   pagination.Limit,
			Total:   count,
			HasMore: pagination.HasMore(count),
		},
	}, nil
}
