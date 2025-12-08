package db

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"
)

var (
	username   = os.Getenv("DB_USERNAME")
	password   = os.Getenv("DB_PASSWORD")
	host       = os.Getenv("DB_HOST")
	port       = os.Getenv("DB_PORT")
	dbName     = os.Getenv("DB_NAME")
	dbSchema   = os.Getenv("DB_SCHEMA")
	dbInstance *DBService
)

type DBService struct {
	DB *sql.DB
}

type DatabaseService interface {
	Open() (*DBService, error)
}

func Open() (*DBService, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, dbName, dbSchema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not reach db connection: %v", err)
	}

	dbInstance = &DBService{
		DB: db,
	}

	return dbInstance, nil
}

func MigrateFS(dbService *DBService, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(dbService.DB, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("error with goose dialect setup: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("could not run gooseUp: %w", err)
	}
	return nil
}
