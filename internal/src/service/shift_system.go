package service

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
)

type ShiftSystem struct {
	ShiftSystemRepo *repository.ShiftSystem
}

func NewShiftSystem(shiftSystemRepo *repository.ShiftSystem) *ShiftSystem {
	return &ShiftSystem{ShiftSystemRepo: shiftSystemRepo}
}

func (s *ShiftSystem) Create(req dto.CreateShiftSystemRequest) (*dto.ShiftSystemResponse, error) {
	id, err := s.ShiftSystemRepo.Create(req)
	if err != nil {
		return nil, err
	}

	return s.ShiftSystemRepo.GetByID(*id)
}

func (s *ShiftSystem) Update(uuid string, req dto.UpdateShiftSystemRequest) (*dto.ShiftSystemResponse, error) {
	err := s.ShiftSystemRepo.Update(uuid, req)
	if err != nil {
		return nil, err
	}
	return s.ShiftSystemRepo.GetByUUID(uuid)
}

func (s *ShiftSystem) Delete(uuid string) error {
	return s.ShiftSystemRepo.Delete(uuid)
}

func (s *ShiftSystem) GetAll() ([]dto.ShiftSystemResponse, error) {
	return s.ShiftSystemRepo.GetAll()
}

func (s *ShiftSystem) GetAllUtils() ([]dto.ShiftSystemResponse, error) {
	return s.ShiftSystemRepo.GetAllUtils()
}

func (s *ShiftSystem) GetByUUID(uuid string) (*dto.ShiftSystemResponse, error) {
	return s.ShiftSystemRepo.GetByUUID(uuid)
}

func (s *ShiftSystem) SetDefault(uuid string) (*dto.ShiftSystemResponse, error) {
	err := s.ShiftSystemRepo.SetDefault(uuid)
	if err != nil {
		return nil, err
	}
	return s.ShiftSystemRepo.GetByUUID(uuid)
}
