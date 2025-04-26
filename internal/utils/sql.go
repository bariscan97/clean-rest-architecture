package utils

import (
	"fmt"
	"strings"
)

func BuildUpdateQueryMap(table string, fields, conditions map[string]interface{}) (string, []interface{}) {
	setClauses := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields)+len(conditions))
	i := 1

	for col, val := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s=$%d", col, i))
		args = append(args, val)
		i++
	}

	query := fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(setClauses, ", "))

	if len(conditions) > 0 {
		condClauses := make([]string, 0, len(conditions))
		for col, val := range conditions {
			condClauses = append(condClauses, fmt.Sprintf("%s=$%d", col, i))
			args = append(args, val)
			i++
		}
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(condClauses, " AND "))
	}

	return query, args
}
