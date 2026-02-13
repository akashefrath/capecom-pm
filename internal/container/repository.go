package container

import (
	"capecom-pm/internal/repositories"
	cacherepo "capecom-pm/internal/repositories/cache"
	mastersreo "capecom-pm/internal/repositories/masters"
	projectrepo "capecom-pm/internal/repositories/project"
	ticketrepo "capecom-pm/internal/repositories/ticket"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	AuthRepo         *repositories.AuthRepo
	UserRepo         *repositories.UserRepo
	SessionRepo      *repositories.SessionRepo
	RoleRepo         *mastersreo.RoleRepo
	GroupRepo        *mastersreo.GroupRepo
	DesignationRepo  *mastersreo.DesignationRepo
	DepartmentRepo   *mastersreo.DepartmentRepo
	MasterDataRepo   *mastersreo.MasterDataRepo
	CacheRepo        *cacherepo.RedisRepo
	FileRepo         *repositories.FileRepo
	ClientRepo       *repositories.ClientRepo
	ProjectRepo      *projectrepo.ProjectRepo
	ProjectAssetRepo *projectrepo.AssetRepo
	ProjectTeamRepo  *projectrepo.TeamRepo
	TicketRepo       *ticketrepo.TicketRepo
	TimeEntryRepo    *ticketrepo.TimeEntryRepo
	HistoryRepo      *ticketrepo.HistoryRepo
	UtilsRepo        *repositories.UtilsRepo
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		AuthRepo:         repositories.NewAuthRepo(db),
		UserRepo:         repositories.NewUserRepo(db),
		SessionRepo:      repositories.NewSessionRepo(db),
		RoleRepo:         mastersreo.NewRoleRepo(db),
		GroupRepo:        mastersreo.NewGroupRepo(db),
		DesignationRepo:  mastersreo.NewDesignationRepo(db),
		DepartmentRepo:   mastersreo.NewDepartmentRepo(db),
		MasterDataRepo:   mastersreo.NewMasterDataRepo(db),
		CacheRepo:        cacherepo.NewRedisRepo(redis),
		FileRepo:         repositories.NewFileRepo(db),
		ClientRepo:       repositories.NewClientRepo(db),
		ProjectRepo:      projectrepo.NewProjectRepo(db),
		ProjectAssetRepo: projectrepo.NewAssetRepo(db),
		ProjectTeamRepo:  projectrepo.NewTeamRepo(db),
		TicketRepo:       ticketrepo.NewTicketRepo(db),
		TimeEntryRepo:    ticketrepo.NewTimeEntryRepo(db),
		HistoryRepo:      ticketrepo.NewHistoryRepo(db),
		UtilsRepo:        repositories.NewUtilsRepo(db),
	}
}
