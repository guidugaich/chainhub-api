
package db

import (
	"database/sql"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			version, dirty, versionErr := m.Version()
			if versionErr == nil {
				filename := migrationFilename(migrationsPath, version)
				if filename != "" {
					log.Printf("migrations: no new migrations to run (current: %s, dirty=%t)", filename, dirty)
				} else {
					log.Printf("migrations: no new migrations to run (current version: %d, dirty=%t)", version, dirty)
				}
			} else {
				log.Printf("migrations: no new migrations to run")
			}
			return nil
		}
		return err
	}

	version, dirty, versionErr := m.Version()
	if versionErr == nil {
		filename := migrationFilename(migrationsPath, version)
		if filename != "" {
			log.Printf("migrations: applied new migrations (current: %s, dirty=%t)", filename, dirty)
		} else {
			log.Printf("migrations: applied new migrations (current version: %d, dirty=%t)", version, dirty)
		}
	} else {
		log.Printf("migrations: applied new migrations")
	}
	return nil
}

func migrationFilename(migrationsPath string, version uint) string {
	dir := strings.TrimPrefix(migrationsPath, "file://")
	if dir == "" {
		return ""
	}

	matches, err := filepath.Glob(filepath.Join(dir, "*_*.up.sql"))
	if err != nil {
		return ""
	}

	for _, match := range matches {
		base := filepath.Base(match)
		underscore := strings.Index(base, "_")
		if underscore <= 0 {
			continue
		}

		parsed, err := strconv.ParseUint(base[:underscore], 10, 64)
		if err != nil {
			continue
		}

		if uint(parsed) == version {
			return base
		}
	}

	return ""
}
