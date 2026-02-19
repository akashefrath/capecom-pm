package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	domainerrors "github.com/akashefrath/capecom-pm/internal/domain/error"
	"github.com/go-sql-driver/mysql"
)

var (
	DuplicateCode = uint16(1062)
)

type DuplicateError struct {
	Key   string
	Value string
}

func SQLErrorFinder(err error) error {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) {

		if mysqlErr.Number == DuplicateCode {
			matches := dupRegex.FindStringSubmatch(mysqlErr.Message)
			if len(matches) != 3 {
				return err
			}
			msg := fmt.Sprintf("Duplicate entry '%s' for key '%s'", matches[1], matches[2])
			errorInvalid := domainerrors.NewWithCodeNoTr(http.StatusConflict, msg, "", "duplicate error")
			return errorInvalid
		}
	} else if errors.Is(err, sql.ErrNoRows) {

		return domainerrors.NewWithCode(http.StatusNotFound, "item_not_found", "", "")

	}

	return err
}

func SQLErrorParser(err error, code uint16) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == code {
			return true
		}
	}

	return false
}

var dupRegex = regexp.MustCompile(`Duplicate entry '(.+)' for key '(.+)'`)

func ParseDuplicate(err error) (error, bool) {
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) || mysqlErr.Number != DuplicateCode {
		return nil, false
	}

	matches := dupRegex.FindStringSubmatch(mysqlErr.Message)
	if len(matches) != 3 {
		return nil, false
	}
	msg := fmt.Sprintf("Duplicate entry '%s' for key '%s'", matches[1], matches[2])
	errorInvalid := domainerrors.NewWithCodeNoTr(http.StatusConflict, msg, "", "duplicate error")
	return errorInvalid, true
}
