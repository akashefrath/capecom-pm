package container

import (
	"github.com/akashefrath/capecom-pm/internal/config"
	"github.com/akashefrath/capecom-pm/internal/database"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
	utilsrepository "github.com/akashefrath/capecom-pm/internal/src/repository/utils"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	Session               *repository.Session
	User                  *repository.User
	Role                  *utilsrepository.Role
	Redis                 *utilsrepository.Redis
	Utils                 *utilsrepository.Utils
	AttendancePolicy      *repository.AttendancePolicy
	AttendancePolicyGroup *repository.AttendancePolicyGroup
	ShiftSystem           *repository.ShiftSystem
	ShiftSystemGroup      *repository.ShiftSystemGroup
	TimeClock             *repository.TimeClock
	AttendanceSummary     *repository.AttendanceSummary
}

func NewRepository(db *sqlx.DB, config *config.Config, redis *redis.Client, dbTX *database.Database) *Repositories {
	attendanceSummary := repository.NewAttendanceSummary(db, dbTX)
	return &Repositories{
		Session:               repository.NewSession(db, config),
		User:                  repository.NewUser(db, dbTX),
		Role:                  utilsrepository.NewRole(db),
		Redis:                 utilsrepository.NewRedis(redis),
		Utils:                 utilsrepository.NewUtils(db),
		AttendancePolicy:      repository.NewAttendancePolicy(db),
		AttendancePolicyGroup: repository.NewAttendancePolicyGroup(db),
		ShiftSystem:           repository.NewShiftSystem(db),
		ShiftSystemGroup:      repository.NewShiftSystemGroup(db),
		TimeClock:             repository.NewTimeClock(db, dbTX, attendanceSummary),
		AttendanceSummary:     attendanceSummary,
	}
}
