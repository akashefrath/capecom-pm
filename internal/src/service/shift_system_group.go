package service

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
)

type ShiftSystemGroup struct {
	ShiftSystemGroupRepo *repository.ShiftSystemGroup
	ShiftSystemRepo      *repository.ShiftSystem
}

func NewShiftSystemGroup(shiftSystemGroupRepo *repository.ShiftSystemGroup, shiftSystemRepo *repository.ShiftSystem) *ShiftSystemGroup {
	return &ShiftSystemGroup{
		ShiftSystemGroupRepo: shiftSystemGroupRepo,
		ShiftSystemRepo:      shiftSystemRepo,
	}
}

func (s *ShiftSystemGroup) Create(req dto.CreateShiftSystemGroupRequest) (*dto.ShiftSystemGroupResponse, error) {
	shiftSystemID, err := s.ShiftSystemRepo.GetIDByUUID(req.ShiftSystemID)
	if err != nil {
		return nil, err
	}

	id, err := s.ShiftSystemGroupRepo.Create(req, *shiftSystemID)
	if err != nil {
		return nil, err
	}
	return s.ShiftSystemGroupRepo.GetByID(*id)
}

func (s *ShiftSystemGroup) Update(uuid string, req dto.UpdateShiftSystemGroupRequest) (*dto.ShiftSystemGroupResponse, error) {
	shiftSystemID, err := s.ShiftSystemRepo.GetIDByUUID(req.ShiftSystemID)
	if err != nil {
		return nil, err
	}

	err = s.ShiftSystemGroupRepo.Update(uuid, req, *shiftSystemID)
	if err != nil {
		return nil, err
	}
	return s.ShiftSystemGroupRepo.GetByUUID(uuid)
}

func (s *ShiftSystemGroup) Delete(uuid string) error {
	return s.ShiftSystemGroupRepo.Delete(uuid)
}

func (s *ShiftSystemGroup) GetAll() ([]dto.ShiftSystemGroupResponse, error) {
	return s.ShiftSystemGroupRepo.GetAll()
}

func (s *ShiftSystemGroup) GetByUUID(uuid string) (*dto.ShiftSystemGroupResponse, error) {
	return s.ShiftSystemGroupRepo.GetByUUID(uuid)
}

func (s *ShiftSystemGroup) AssignUsers(groupUUID string, req dto.AssignUsersToShiftGroupRequest) error {
	groupID, err := s.ShiftSystemGroupRepo.GetIDByUUID(groupUUID)
	if err != nil {
		return err
	}
	return s.ShiftSystemGroupRepo.AssignUsers(*groupID, req.UserUUIDs)
}

func (s *ShiftSystemGroup) RemoveUsers(req dto.RemoveUsersFromShiftGroupRequest) error {
	return s.ShiftSystemGroupRepo.RemoveUsers(req.UserUUIDs)
}

func (s *ShiftSystemGroup) GetUsersInGroup(groupUUID string) ([]dto.UserMinimalResponse, error) {
	return s.ShiftSystemGroupRepo.GetUsersInGroup(groupUUID)
}
