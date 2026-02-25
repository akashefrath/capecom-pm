package service

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
	"github.com/akashefrath/capecom-pm/internal/utils"
)

type TimeClock struct {
	TimeClockRepo *repository.TimeClock
}

func NewTimeClock(timeClockRepo *repository.TimeClock) *TimeClock {
	return &TimeClock{TimeClockRepo: timeClockRepo}
}

func (s *TimeClock) GetTodayDetails(userID *int64) (*dto.AttendanceDetailsListWithSummary, error) {
	data, err := s.TimeClockRepo.GetTodayDetails(userID)

	return data, err

}
func (s *TimeClock) AdvancePunch(employeeID int64, req dto.TimeClockRequest, punchType string) (*dto.TimeClockResponse, error) {

	lastLog, err := s.TimeClockRepo.GetUsersLastLog(employeeID)
	if err != nil && punchType != models.LogIn {
		return nil, err

	}

	currentPunch := punchType
	validatePunch := validatePunch(lastLog, currentPunch)
	if validatePunch != nil {
		return nil, validatePunch
	}
	canPunch := canPunch(lastLog, currentPunch)
	if canPunch {
		id, err := s.TimeClockRepo.TimePunch(employeeID, req, currentPunch)
		if err != nil {
			return nil, err
		}
		return s.TimeClockRepo.GetByID(*id)
	} else {
		return nil, domainerrors.CantPerformThis
	}
}

func (s *TimeClock) TimeOut(req *dto.TimeClockTimeOutRequest, employeeID int64) error {
	pendingSummary, err := s.TimeClockRepo.AttendanceSummary.GetCurrentPendingSummaryWithID(employeeID)

	if pendingSummary == nil || err != nil {
		return err
	}
	baseDate := pendingSummary.CreatedAt

	finalTime, err := utils.CombineDateTime(*baseDate, *req.Time)
	if err != nil {

		return err
	}

	_, err = s.TimeClockRepo.TimeOutPunch(employeeID, dto.TimeClockRequest{
		Source:    req.Source,
		DeviceID:  req.DeviceID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Remarks:   req.Remarks,
	}, models.TimeOut, finalTime)
	if err != nil {
		return err
	}

	return nil

}

func canPunch(prev *string, next string) bool {

	if prev == nil {
		return next == models.LogIn
	}

	switch *prev {

	case models.LogIn:
		return next == models.BrakeIn || next == models.LogOut

	case models.BrakeIn:
		return next == models.BrakeOut

	case models.BrakeOut:
		return next == models.BrakeIn || next == models.LogOut

	case models.LogOut, models.TimeOut:
		return next == models.LogIn
	}

	return false
}
func validatePunch(prev *string, next string) error {

	// First punch of the day
	if prev == nil {
		if next == models.LogIn {
			return nil
		}
		return domainerrors.MustClockInFirst
	}

	// same button double tap
	if *prev == next {
		switch next {
		case models.LogIn:
			return domainerrors.AlreadyClockedIn
		case models.BrakeIn:
			return domainerrors.AlreadyOnBreak
		case models.BrakeOut:
			return domainerrors.NotOnBreak
		case models.LogOut:
			return domainerrors.AlreadyClockedOut
		}
		return domainerrors.DuplicatePunch
	}

	switch *prev {

	// ===== WORKING =====
	case models.LogIn:
		if next == models.BrakeIn || next == models.LogOut {
			return nil
		}
		return domainerrors.InvalidPunchOrder

	// ===== ON BREAK =====
	case models.BrakeIn:
		if next == models.BrakeOut {
			return nil
		}
		return domainerrors.MustBreakOutFirst

	// ===== RESUMED WORK =====
	case models.BrakeOut:
		if next == models.BrakeIn || next == models.LogOut {
			return nil
		}
		return domainerrors.InvalidPunchOrder

	// ===== DAY CLOSED =====
	case models.LogOut, models.TimeOut:
		// Important: block re-login same day
		if next == models.LogIn {
			return domainerrors.AlreadyClockedOut
		}
		return domainerrors.CantPerformThis
	}

	return domainerrors.InvalidPunchOrder
}
