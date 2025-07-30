package sqlboiler

import "database/sql"

type SQLBoiler interface {
	ConnectDB() *sql.DB
}
