package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

var (
	username = os.Getenv("PG_USERNAME")
	pwd      = os.Getenv("PG_PASSWORD")
	host     = os.Getenv("PG_HOST")
	port     = os.Getenv("PG_PORT")
	dbName   = os.Getenv("PG_DATABASE")
	mode     = os.Getenv("PG_SSLMODE")
)

func Connect() (*pgxpool.Pool, error) {
	ctx := context.Background()
	path := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, pwd, host, port, dbName, mode)

	dbpool, err := pgxpool.New(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("unable to verify database connection status: %w", err)
	}

	if err := migrateUp(path); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("unable to complete database migrations: %w", err)
	}

	return dbpool, err
}

func migrateUp(path string) error {
	m, err := migrate.New("file://internal/database/migrations", path)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		return err
	}
	return nil
}
