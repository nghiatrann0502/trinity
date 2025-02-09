package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type mysqlDB struct {
	db *sql.DB
}

var _ DBEngine = (*mysqlDB)(nil)

func NewMySQLDatabase(url string) (DBEngine, error) {
	db, err := sql.Open("mysql", string(url))
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &mysqlDB{db}, nil
}

func (s *mysqlDB) Health() (bool, map[string]string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return false, stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return true, stats
}

func (s *mysqlDB) GetDB() *sql.DB {
	return s.db
}

func (s *mysqlDB) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *mysqlDB) Migrate() error {
	driver, err := mysql.WithInstance(s.db, &mysql.Config{})
	if err != nil {
		return err
	}

	cur, _ := os.Getwd()

	dir := filepath.Dir(cur + "../../../")

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/%s", dir, "db/migrations"),
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return nil
}
