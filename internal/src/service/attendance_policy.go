package service

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
)

type AttendancePolicy struct {
	AttendancePolicyRepo *repository.AttendancePolicy
}

func NewAttendancePolicy(attendancePolicyRepo *repository.AttendancePolicy) *AttendancePolicy {
	return &AttendancePolicy{AttendancePolicyRepo: attendancePolicyRepo}
}

func (s *AttendancePolicy) Create(req dto.CreateAttendancePolicyRequest) (*dto.AttendancePolicyResponse, error) {
	id, err := s.AttendancePolicyRepo.Create(req)
	if err != nil {
		return nil, err
	}

	return s.AttendancePolicyRepo.GetByID(*id)
}

func (s *AttendancePolicy) Update(uuid string, req dto.CreateAttendancePolicyRequest) (*dto.AttendancePolicyResponse, error) {
	err := s.AttendancePolicyRepo.Update(uuid, req)
	if err != nil {
		return nil, err
	}
	return s.AttendancePolicyRepo.GetByUUID(uuid)
}

func (s *AttendancePolicy) Delete(uuid string) error {
	return s.AttendancePolicyRepo.Delete(uuid)
}

func (s *AttendancePolicy) GetAll() ([]dto.AttendancePolicyResponse, error) {
	return s.AttendancePolicyRepo.GetAll()
}
func (s *AttendancePolicy) GetAllUtils() ([]dto.AttendancePolicyResponse, error) {
	return s.AttendancePolicyRepo.GetAllUtils()
}
func (s *AttendancePolicy) GetByUUID(uuid string) (*dto.AttendancePolicyResponse, error) {
	return s.AttendancePolicyRepo.GetByUUID(uuid)
}

func (s *AttendancePolicy) SetDefault(uuid string) (*dto.AttendancePolicyResponse, error) {
	err := s.AttendancePolicyRepo.SetDefault(uuid)
	if err != nil {
		return nil, err
	}
	return s.AttendancePolicyRepo.GetByUUID(uuid)
}
