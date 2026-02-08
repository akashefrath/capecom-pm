package services

import (
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	mastersreo "capecom-pm/internal/repositories/masters"
	"capecom-pm/internal/utils"

	"github.com/gin-gonic/gin"
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

func (s *UserService) CreateUser(c *gin.Context, req dto.CreateUserRequest) error {
	idS, err := s.masterDataRepo.GetUserRelatedIDs(
		req.GroupID,
		req.DesignationID,
		req.DepartmentID,
		req.RoleIDs,
	)
	if err != nil {

		return err
	}
	finalPassword := req.Password
	if finalPassword != "" {
		randomPass, err := utils.RandomPassword(8)
		if err != nil {

			return domainerrors.ErrInternal

		}
		finalPassword = utils.HashPassword(randomPass)
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
	err = s.userRepo.CreateUserWithRoles(user, idS.RoleIDs)
	if err != nil {

		return err
	}
	return nil
}
