package service

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
)

type TimeClock struct {
	TimeClockRepo *repository.TimeClock
}

func NewTimeClock(timeClockRepo *repository.TimeClock) *TimeClock {
	return &TimeClock{TimeClockRepo: timeClockRepo}
}

func (s *TimeClock) GetTodayDetails(userID *int64) (*[]dto.AttendanceDetails, error) {
	data, err := s.TimeClockRepo.GetTodayDetails(userID)

	return data, err

}

func (s *TimeClock) ClockIn(employeeID int64, req dto.TimeClockRequest) (*dto.TimeClockResponse, error) {

	lastLog, err := s.TimeClockRepo.GetUsersLastLog(employeeID)

	if lastLog != nil {
		return nil, domainerrors.CantPerformThis

	}
	id, err := s.TimeClockRepo.TimePunch(employeeID, req, models.LogIn)
	if err != nil {
		return nil, err
	}
	return s.TimeClockRepo.GetByID(*id)
}

func (s *TimeClock) ClockOut(employeeID int64, req dto.TimeClockRequest) (*dto.TimeClockResponse, error) {
	lastLog, err := s.TimeClockRepo.GetUsersLastLog(employeeID)

	if lastLog == nil || (*lastLog == models.LogOut || *lastLog == models.TimeOut) {
		return nil, domainerrors.CantPerformThis

	}

	id, err := s.TimeClockRepo.TimePunch(employeeID, req, models.LogOut)
	if err != nil {
		return nil, err
	}
	return s.TimeClockRepo.GetByID(*id)
}

func (s *TimeClock) BreakIn(employeeID int64, req dto.TimeClockRequest) (*dto.TimeClockResponse, error) {

	lastLog, err := s.TimeClockRepo.GetUsersLastLog(employeeID)

	if lastLog == nil || (*lastLog != models.LogIn && *lastLog != models.BrakeOut) {
		return nil, domainerrors.CantPerformThis

	}
	id, err := s.TimeClockRepo.TimePunch(employeeID, req, models.BrakeIn)
	if err != nil {
		return nil, err
	}
	return s.TimeClockRepo.GetByID(*id)
}

func (s *TimeClock) BreakOut(employeeID int64, req dto.TimeClockRequest) (*dto.TimeClockResponse, error) {
	lastLog, err := s.TimeClockRepo.GetUsersLastLog(employeeID)

	if lastLog == nil || *lastLog != models.BrakeIn {
		return nil, domainerrors.CantPerformThis

	}
	id, err := s.TimeClockRepo.TimePunch(employeeID, req, models.BrakeOut)
	if err != nil {
		return nil, err
	}
	return s.TimeClockRepo.GetByID(*id)
}
