package internal

import "database/sql"
import _ "github.com/lib/pq"

func NewPostgres() (*sql.DB, error) {
	connStr := "host=postgres user=user password=postgres dbname=test sslmode=disable"
	return sql.Open("postgres", connStr)
}
