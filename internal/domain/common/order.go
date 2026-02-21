package common

import (
	"fmt"
	"strings"
)

type OrderBy struct {
	Key   string `json:"orderby[key]" form:"orderby[key]" query:"orderby[key]"`
	Value string `json:"orderby[value]" form:"orderby[value]" query:"orderby[value]" binding:"omitempty,oneof=asc desc"`
}

func (o *OrderBy) GetQuery(allowedColumns []string) string {
	if o.Key == "" {
		return ""
	}

	column := strings.ToLower(o.Key)
	if !isAllowedColumn(column, allowedColumns) {
		return ""
	}

	direction := "ASC"
	if strings.ToLower(o.Value) == "desc" {
		direction = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", column, direction)
}

func BuildOrderByQuery(orderBy OrderBy, allowedColumns []string) string {
	return orderBy.GetQuery(allowedColumns)
}

func isAllowedColumn(column string, allowedColumns []string) bool {
	for _, allowed := range allowedColumns {
		if strings.ToLower(allowed) == column {
			return true
		}
	}
	return false
}
