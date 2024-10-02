package controller

import (
	"betamart/internal/database"
	"database/sql"
)

type ApiConfig struct {
	Query *database.Queries
	DB  *sql.DB
}
