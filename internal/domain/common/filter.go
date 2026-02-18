package common

import (
	"fmt"
	"strings"

	utilsdto "github.com/akashefrath/capecom-pm/internal/domain/dto/utils"
)

type FilterWithKeys struct {
	Key    string
	Values []string
}

func BuildFilterQuery(filters []FilterWithKeys) (string, []interface{}) {
	var parts []string
	var args []interface{}

	for _, filter := range filters {
		// Create a slice of '?' based on how many values we have
		placeholders := make([]string, len(filter.Values))
		for i := range filter.Values {
			placeholders[i] = "?"
			args = append(args, filter.Values[i])
		}

		// Build the string: "status IN (?, ?)"
		parts = append(parts, fmt.Sprintf("%s IN (%s)", filter.Key, strings.Join(placeholders, ",")))
	}

	return strings.Join(parts, " AND "), args
}

func FilterApplied(filter []FilterWithKeys) []utilsdto.Filter {
	var filters []utilsdto.Filter
	for _, filter := range filter {
		for _, value := range filter.Values {
			filters = append(filters, utilsdto.Filter{Key: filter.Key, Value: value})
		}

	}
	return filters

}
