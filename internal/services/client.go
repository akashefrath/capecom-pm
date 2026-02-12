package services

import (
	"capecom-pm/internal/domain/common"
	"capecom-pm/internal/domain/dto"
	domainerrors "capecom-pm/internal/domain/error"
	"capecom-pm/internal/domain/models"
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	"net/http"
)

type ClientService struct {
	clientRepo *repositories.ClientRepo
	userRepo   *repositories.UserRepo
	redisRepo  *cacherepo.RedisRepo
}

func NewClientService(clientRepo *repositories.ClientRepo, userRepo *repositories.UserRepo, redisRepo *cacherepo.RedisRepo) *ClientService {
	return &ClientService{clientRepo: clientRepo, userRepo: userRepo, redisRepo: redisRepo}
}

func (s *ClientService) Create(req dto.CreateClientRequest, userUUID string) (*dto.ClientResponse, error) {
	var createdBy *uint64

	if userUUID != "" {
		userID, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
		if err != nil || userID == nil {
			return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "client_service", "get_user_id")
		}
		uid := uint64(*userID)
		createdBy = &uid
	}

	client := &models.Client{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		BaseModel: models.NewBase(createdBy),
	}

	result, err := s.clientRepo.Create(client)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ClientService) GetByUUID(uuid string) (*dto.ClientResponse, error) {
	client, err := s.clientRepo.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, domainerrors.NewWithCode(http.StatusNotFound, domainerrors.ErrClientNotFound.Error(), "service", "GetClientByUUID")
	}
	return client, nil
}

func (s *ClientService) GetClients(pagination common.Pagination) (*dto.ListWithMeta, error) {
	return s.clientRepo.GetClients(pagination)
}

func (s *ClientService) Update(uuid string, req dto.UpdateClientRequest) (*dto.ClientResponse, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}

	if len(updates) == 0 {
		return s.clientRepo.FindByUUID(uuid)
	}

	return s.clientRepo.Update(uuid, updates)
}
