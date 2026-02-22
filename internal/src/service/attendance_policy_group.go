package service

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
)

type AttendancePolicyGroup struct {
	AttendancePolicyGroupRepo *repository.AttendancePolicyGroup
}

func NewAttendancePolicyGroup(attendancePolicyGroupRepo *repository.AttendancePolicyGroup) *AttendancePolicyGroup {
	return &AttendancePolicyGroup{AttendancePolicyGroupRepo: attendancePolicyGroupRepo}
}

func (s *AttendancePolicyGroup) Create(req dto.CreateAttendancePolicyGroupRequest) (*dto.AttendancePolicyGroupSingleResponse, error) {
	policyID, err := s.AttendancePolicyGroupRepo.GetPolicyIDByUUID(req.AttendancePolicyUUID)
	if err != nil {
		return nil, err
	}

	id, err := s.AttendancePolicyGroupRepo.Create(req, *policyID)
	if err != nil {
		return nil, err
	}

	return s.AttendancePolicyGroupRepo.GetByID(*id)
}

func (s *AttendancePolicyGroup) Update(uuid string, req dto.UpdateAttendancePolicyGroupRequest) (*dto.AttendancePolicyGroupSingleResponse, error) {
	policyID, err := s.AttendancePolicyGroupRepo.GetPolicyIDByUUID(req.AttendancePolicyUUID)
	if err != nil {
		return nil, err
	}

	err = s.AttendancePolicyGroupRepo.Update(uuid, req, *policyID)
	if err != nil {
		return nil, err
	}
	return s.AttendancePolicyGroupRepo.GetByUUID(uuid)
}

func (s *AttendancePolicyGroup) Delete(uuid string) error {
	return s.AttendancePolicyGroupRepo.Delete(uuid)
}

func (s *AttendancePolicyGroup) GetAll() ([]dto.AttendancePolicyGroupSingleResponse, error) {
	return s.AttendancePolicyGroupRepo.GetAll()
}

func (s *AttendancePolicyGroup) GetByUUID(uuid string) (*dto.AttendancePolicyGroupSingleResponse, error) {
	return s.AttendancePolicyGroupRepo.GetByUUID(uuid)
}

func (s *AttendancePolicyGroup) AssignUsers(groupUUID string, req dto.AssignUsersToGroupRequest) error {
	groupID, err := s.AttendancePolicyGroupRepo.GetIDByUUID(groupUUID)
	if err != nil {
	 
		return err
	}
	return s.AttendancePolicyGroupRepo.AssignUsers(*groupID, req.UserUUIDs)
}

func (s *AttendancePolicyGroup) RemoveUsers(req dto.RemoveUsersFromGroupRequest) error {
	return s.AttendancePolicyGroupRepo.RemoveUsers(req.UserUUIDs)
}

func (s *AttendancePolicyGroup) GetUsersInGroup(groupUUID string) ([]dto.UserMinimalResponse, error) {
	return s.AttendancePolicyGroupRepo.GetUsersInGroup(groupUUID)
}
