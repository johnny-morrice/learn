package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dbDriverName = "postgres"

func openDb(databaseURL string) (*sql.DB, error) {
	err := validateDatabaseParam()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(dbDriverName, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening postgres connection: %w", err)
	}

	return db, nil
}

func openGorm(postgresDSN string) (*gorm.DB, error) {
	const retries = 5
	for i := 0; i < retries; i++ {
		db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
		if err != nil {
			log.Printf("error opening gorm connection, retrying: %s", err.Error())
			time.Sleep(time.Second)
			continue
		}

		return db, nil
	}
	return nil, errors.New("too many errors connecting to postgres")
}

func openMigration(databaseURL, migrationsPath string) (*migrate.Migrate, error) {
	db, err := openDb(databaseURL)
	if err != nil {
		return nil, err
	}

	driver, err := pgmigrate.WithInstance(db, &pgmigrate.Config{})
	if err != nil {
		return nil, fmt.Errorf("error creating migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, dbDriverName, driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func migrateDbUp(databaseURL, migrationsPath string) error {
	migration, err := openMigration(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer migration.Close()
	err = migration.Up()
	if err != nil {
		return fmt.Errorf("error migrating up: %w", err)
	}
	log.Println("migrated UP ok")
	return nil
}

func migrateDbDown(databaseURL, migrationsPath string) error {
	migration, err := openMigration(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer migration.Close()
	err = migration.Down()
	if err != nil {
		return fmt.Errorf("error migrating down: %w", err)
	}
	log.Println("migrated DOWN ok")
	return nil
}
