package app

import (
	"database/sql"
)

type Context struct {
	DB      *sql.DB
	Setting Settings
}
