package database

import (
	"errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"backend/internal/logger"
)

func RunMigrations() {
	start := time.Now()

	logger.Info.Println("starting database migrations")
	logger.Info.Println("database driver: sqlite3")
	logger.Info.Println("database file: projects_db.db")
	logger.Info.Println("migrations source: file://migrations")

	m, err := migrate.New(
		"file://migrations",
		"sqlite3://projects_db.db",
	)
	if err != nil {
		logger.Fatal.Println("failed to create migrate instance:", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		logger.Error.Println("failed to get current migration version:", err)
	} else {
		logger.Info.Printf("current migration version: %d, dirty: %v\n", version, dirty)
	}

	err = m.Up()
	switch {
	case err == nil:
		logger.Info.Println("migrations applied successfully")
	case errors.Is(err, migrate.ErrNoChange):
		logger.Info.Println("no new migrations to apply")
	default:
		logger.Fatal.Println("migration failed:", err)
	}

	logger.Info.Printf(
		"migration process finished in %s\n",
		time.Since(start),
	)
}
