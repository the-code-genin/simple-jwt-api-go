package errors

import (
	"strings"

	"github.com/jackc/pgx/v5"
)

func IsNoRecordError(err error) bool {
	return strings.Contains(err.Error(), pgx.ErrNoRows.Error())
}
