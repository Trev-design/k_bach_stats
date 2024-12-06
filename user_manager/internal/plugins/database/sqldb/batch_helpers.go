package sqldb

import (
	"fmt"
	"strings"
)

func makeBatchQuery(selectStmt, placeholder string, quantity int) string {
	placeholders := make([]string, quantity)

	for i := 0; i < quantity; i++ {
		placeholders[i] = placeholder
	}

	return fmt.Sprintf(selectStmt, strings.Join(placeholders, ","))
}
