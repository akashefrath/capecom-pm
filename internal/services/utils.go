package services

import (
	"capecom-pm/internal/domain/dto"
	"capecom-pm/internal/repositories"
)

type UtilsService struct {
	utilsRepo *repositories.UtilsRepo
}

func NewUtilsService(utilsRepo *repositories.UtilsRepo) *UtilsService {
	return &UtilsService{utilsRepo: utilsRepo}
}

func (s *UtilsService) GetAll() (*dto.UtilsResponse, error) {
	roles, err := s.utilsRepo.GetRoles()
	if err != nil {
		return nil, err
	}
	userGroups, err := s.utilsRepo.GetUserGroups()
	if err != nil {
		return nil, err
	}
	designations, err := s.utilsRepo.GetDesignations()
	if err != nil {
		return nil, err
	}
	departments, err := s.utilsRepo.GetDepartments()
	if err != nil {
		return nil, err
	}
	clients, err := s.utilsRepo.GetClients()
	if err != nil {
		return nil, err
	}
	ticketTypes, err := s.utilsRepo.GetTicketTypes()
	if err != nil {
		return nil, err
	}
	users, err := s.utilsRepo.GetUsers()
	if err != nil {
		return nil, err
	}

	return &dto.UtilsResponse{
		Roles:        roles,
		UserGroups:   userGroups,
		Designations: designations,
		Departments:  departments,
		Clients:      clients,
		TicketTypes:  ticketTypes,
		Users:        users,
	}, nil
}

func (s *UtilsService) GetRoles() ([]dto.UtilOption, error) {
	return s.utilsRepo.GetRoles()
}

func (s *UtilsService) GetUserGroups() ([]dto.UtilOption, error) {
	return s.utilsRepo.GetUserGroups()
}

func (s *UtilsService) GetDesignations() ([]dto.UtilOption, error) {
	return s.utilsRepo.GetDesignations()
}

func (s *UtilsService) GetDepartments() ([]dto.UtilOption, error) {
	return s.utilsRepo.GetDepartments()
}

func (s *UtilsService) GetClients() ([]dto.UtilOption, error) {
	return s.utilsRepo.GetClients()
}

func (s *UtilsService) GetTicketTypes() ([]dto.UtilOption, error) {
	return s.utilsRepo.GetTicketTypes()
}

func (s *UtilsService) GetUsers() ([]dto.UtilOption, error) {
	return s.utilsRepo.GetUsers()
}
