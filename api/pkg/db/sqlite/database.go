package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
)

// The OpenDB() function returns a sql.DB connection pool.
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:database.db")
	if err != nil {
		return nil, err
	}

	// Use Ping to establish a new connection to the database, if the connection couldn't be
	// established successfully this will return an error.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}

// Run migrate scripts to create database if not created before.
func RunMigrateScripts(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating sqlite3 db driver failed %s", err)
	}

	res := bindata.Resource(AssetNames(),
		func(name string) ([]byte, error) {
			return Asset(name)
		})

	d, err := bindata.WithInstance(res)
	m, err := migrate.NewWithInstance("go-bindata", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("initializing db migration failed %s", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrating database failed %s", err)
	}

	// err = m.Down()
	// if err != nil && err != migrate.ErrNoChange {
	// 	return fmt.Errorf("migrating database failed %s", err)
	// }

	return nil
}
