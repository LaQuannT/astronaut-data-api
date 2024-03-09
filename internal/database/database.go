package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	db *pgxpool.Pool
}

func NewPostgresDB(connStr string) *PostgresDB {
	ctx := context.Background()

	log.Println("Connection to database")

	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("unable to connect to database %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		log.Fatalf("unable to verify database connection status: %v", err)
	}

	log.Println("Connected to database")
	return &PostgresDB{
		db: dbpool,
	}
}

func (p *PostgresDB) Init() (*pgxpool.Pool, error) {
	log.Println("Attempting to build tables")

	if err := p.createUserTable(); err != nil {
		return nil, fmt.Errorf("failed to build 'USER' table: %v", err)
	}

	if err := p.createAstronautTable(); err != nil {
		return nil, fmt.Errorf("failed to build 'ASTRONAUT' table: %v", err)
	}

	// TODO - create a function that populates DB with csv data if DB is empty

	return p.db, nil
}

func (p *PostgresDB) createUserTable() error {
	query := `
DROP TYPE IF EXISTS user_role;
  CREATE TYPE user_role AS ENUM('user', 'admin');

  CREATE TABLE IF NOT EXISTS "user" (
  id SERIAL NOT NULL PRIMARY KEY,
  first_name VARCHAR(50) NOT NULL,
  surname VARCHAR(50) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  api_key TEXT NOT NULL UNIQUE,
  role user_role NOT NULL DEFAULT 'user',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := p.db.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) createAstronautTable() error {
	query := `
  DROP TYPE IF EXISTS status;
  CREATE TYPE status AS  ENUM('retired', 'active', 'management', 'deceased');

  DROP TYPE IF EXISTS gender;
  CREATE TYPE gender AS ENUM('male', 'female');


  CREATE TABLE IF NOT EXISTS astronaut (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  year INT NOT NULL,
  "group" INT NOT NULL,
  status status NOT NULL,
  birth_date VARCHAR(20) NOT NULL,
  birth_place VARCHAR(50) NOT NULL,
  gender gender NOT NULL,
  alma_mater VARCHAR(50)[],
  undergraduate_major VARCHAR(50)[],
  graduate_major VARCHAR(50)[],
  military_rank VARCHAR(50),
  military_branch VARCHAR(50),
  space_flights INT NOT NULL,
  space_flight_hrs INT NOT NULL,
  space_walks INT NOT NULL,
  space_walk_hrs INT NOT NULL,
  missions VARCHAR(50)[],
  death_date VARCHAR(20),
  death_misson VARCHAR(50)
);`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := p.db.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
