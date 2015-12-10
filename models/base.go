package models

import (
	"strconv"
	"strings"
)

// Base struct for models that whan to user complex query schema
type Base struct {
	Query map[string][]string `sql:"-" json:",omitempty"`
}

var (
	queryIdentifiers = map[string]string{"gte": ">=", "gt": ">", "lte": "<=", "lt": "<", "eq": "="}
	paramDelimiter   = "|"
)

func (b *Base) BuildQuery() string {
	var query string

	// identifierKey = gte, identifierValue = >=
	for identifierKey, identifierValue := range queryIdentifiers {
		// val = [quantity|200]
		if queryWithKeyValues, ok := b.Query[identifierKey]; ok {

			// currQueryValue  = quantity|200
			for _, currQueryValue := range queryWithKeyValues {
				// splitted = [quantity 200]
				splitted := strings.Split(currQueryValue, paramDelimiter)
				if len(splitted) != 2 {
					continue
				}

				if query != "" {
					query += " and "
				}

				if isNumber(splitted[1]) {
					query += splitted[0] + identifierValue + splitted[1]
				} else {
					query += splitted[0] + identifierValue + "'" + splitted[1] + "'"
				}
			}
		}
	}

	return query
}

func isNumber(s string) bool {
	_, err1 := strconv.Atoi(s)
	_, err2 := strconv.ParseFloat(s, 64)
	if err1 == nil && err2 == nil {
		return true
	}
	return false
}
