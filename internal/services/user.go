package services

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	mastersreo "capecom-pm/internal/repositories/masters"
	"capecom-pm/internal/utils"
	"net/http"
)

type UserService struct {
	userRepo       *repositories.UserRepo
	masterDataRepo *mastersreo.MasterDataRepo
}

func NewUserService(
	userRepo *repositories.UserRepo,
	masterDataRepo *mastersreo.MasterDataRepo,
) *UserService {
	return &UserService{
		userRepo:       userRepo,
		masterDataRepo: masterDataRepo,
	}
}

func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	idS, err := s.masterDataRepo.GetUserRelatedIDs(
		req.GroupID,
		req.DesignationID,
		req.DepartmentID,
		req.RoleIDs,
	)
	if err != nil {
		return nil, domainerrors.NewWithCode(http.StatusBadRequest, err.Error(), "service", "check_fg")
	}
	finalPassword := req.Password

	if finalPassword == "" {
		randomPass, err := utils.RandomPassword(8)
		if err != nil {
			return nil, domainerrors.Internal("service", "password")
		}
		finalPassword = utils.HashPassword(randomPass)

	} else {
		finalPassword = utils.HashPassword(req.Password)
	}
	user := &models.User{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		CountryCode:   req.CountryCode,
		PasswordHash:  finalPassword,
		EmployeeID:    req.EmployeeID,
		GroupID:       idS.GroupID,
		DesignationID: idS.DesignationID,
		DepartmentID:  idS.DepartmentID,
		BaseModel:     models.NewBase(nil),
	}
	usr, err := s.userRepo.CreateUserWithRoles(user, idS.RoleIDs)
	if err != nil {

		return nil, err
	}
	return usr, nil
}

func (s *UserService) GetUserByID(uuid string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrUserNotFound.Error(), "service", "GetUserByID")
	}
	return user, nil
}
func (s *UserService) GetUser(pagination common.Pagination) (*dto.ListWithMeta, error) {
	user, err := s.userRepo.GetUsers(pagination)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrUserNotFound.Error(), "service", "GetUserByID")
	}
	return user, nil
}
