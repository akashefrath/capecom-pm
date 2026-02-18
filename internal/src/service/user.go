package service

import (
	"github.com/akashefrath/capecom-pm/internal/database"
	"github.com/akashefrath/capecom-pm/internal/domain/common"
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	utilsdto "github.com/akashefrath/capecom-pm/internal/domain/dto/utils"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
	"github.com/akashefrath/capecom-pm/internal/utils"
)

type User struct {
	UserRepo *repository.User
}

func NewUser(userRepo *repository.User) *User {
	return &User{UserRepo: userRepo}
}

func (s User) Create(createdBy int64, req dto.CreateUserRequest) (*dto.User, error) {

	if req.Password == "" {
		req.Password = utils.HashPassword("Capecom@2025")
	} else {
		req.Password = utils.HashPassword(req.Password)
	}
	userId, err := s.UserRepo.Create(createdBy, req)

	if err != nil {
		d, isError := database.ParseDuplicate(err)
		if isError {
			return nil, d
		}
		return nil, err
	}
	user, err := s.UserRepo.FindUserByID(*userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s User) GetAll(pg common.Pagination, filter []common.FilterWithKeys) (*utilsdto.ListWithMeta, error) {

	return s.UserRepo.GetAll(pg, filter)
}

func (s User) GetByID(uuid string) (*dto.User, error) {
	return s.UserRepo.GetByID(uuid)
}
