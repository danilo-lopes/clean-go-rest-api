// Clean Architecture - Frameworks & Drivers Layer
// Database migration runner
package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migration struct {
	db            *sql.DB
	locationFiles string
}

func NewMigration(db *sql.DB, migrationFiles string) *Migration {
	return &Migration{db: db, locationFiles: migrationFiles}
}

func (m *Migration) RunMigrations() {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		log.Println("Error creating migration driver:", err)
		return
	}
	migration, err := migrate.NewWithDatabaseInstance(
		m.locationFiles,
		"postgres",
		driver,
	)
	if err != nil {
		log.Println("Error creating new migration instance:", err)
		return
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("Error running migrations:", err)
	}
	_ = driver.Close()
}
