package database

import (
	"database/sql"
	"time"
)

type DBEngine interface {
	Health() (bool, map[string]string)
	GetDB() *sql.DB
	Close()
	Migrate(dir string) error
}

type BaseSql struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
