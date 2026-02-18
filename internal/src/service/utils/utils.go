package utilsservice

import (
	utilsrepository "github.com/akashefrath/capecom-pm/internal/src/repository/utils"
)

type Utils struct {
	UtilsRepo *utilsrepository.Utils
}

func NewUtils(utilsRepo *utilsrepository.Utils) *Utils {
	return &Utils{UtilsRepo: utilsRepo}
}
