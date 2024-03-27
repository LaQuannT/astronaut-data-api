package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/LaQuannT/astronaut-data-api/internal/config"
	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/gocarina/gocsv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type PostgresDB struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPostgresDB(connStr string, l *slog.Logger) *PostgresDB {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		l.Log(ctx, config.LevelTrace, "unable to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		l.Log(ctx, config.LevelTrace, "unable to verify database connection status", slog.Any("error", err))
		os.Exit(1)
	}

	l.Info("Connected to database")
	return &PostgresDB{
		db:     dbpool,
		logger: l,
	}
}

func (p *PostgresDB) Init() (*pgxpool.Pool, error) {
	count, err := p.checkAstronautCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get 'ASTRONAUT' count: %w", err)
	}

	if count == 0 {
		p.logger.Info("Attempting to seed Astronaut table")
		if err := p.populateAstronautTable(); err != nil {
			return nil, fmt.Errorf("faild to seed 'ASTRONAUT' table: %w", err)
		}
		p.logger.Info("'ASTRONAUT' table seeded")
	}

	return p.db, nil
}

func (p *PostgresDB) checkAstronautCount() (int, error) {
	count := 0
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT COUNT(*) FROM astronaut;`

	row := p.db.QueryRow(ctx, query)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p *PostgresDB) populateAstronautTable() error {
	// TODO - refactor to batch insert
	ctx := context.Background()

	file, err := os.Open("astronauts.csv")
	if err != nil {
		return fmt.Errorf("unable to open CSV file: %w", err)
	}
	defer file.Close()

	var astronauts []*model.Astronaut

	if err := gocsv.UnmarshalFile(file, &astronauts); err != nil {
		return fmt.Errorf("unable to unmarshal CSV data to stuct slice: %w", err)
	}

	query := `
  INSERT INTO astronaut
  (name, year, "group", status, birth_date, birth_place, gender, alma_mater, undergraduate_major,
  graduate_major, military_rank, military_branch, space_flights, space_flight_hrs, space_walks,
  space_walk_hrs, missions, death_date, death_mission)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19);`

	for _, a := range astronauts {
		formatStrsToLower(a)

		a.AlmaMater = strings.Split(a.AlmaMaterStr, ";")
		a.UndergraduateMajor = strings.Split(a.UndergraduateMajorStr, ";")
		a.GraduateMajor = strings.Split(a.GraduateMajorStr, ";")
		a.Missions = strings.Split(a.MissionStr, ",")

		_, err := p.db.Exec(ctx, query, a.Name, a.Year, a.Group, a.Status, a.BirthDate, a.BirthPlace, a.Gender, pq.Array(a.AlmaMater),
			pq.Array(a.UndergraduateMajor), pq.Array(a.GraduateMajor), a.MilitaryRank, a.MilitaryBranch, a.SpaceFlights, a.SpaceFlightHours,
			a.SpaceWalks, a.SpaceWalkHours, pq.Array(a.Missions), a.DeathDate, a.DeathMission)
		if err != nil {
			return err
		}
	}

	return nil
}

func formatStrsToLower(a *model.Astronaut) {
	a.Name = strings.ToLower(a.Name)
	a.Status = strings.ToLower(a.Status)
	a.BirthPlace = strings.ToLower(a.BirthPlace)
	a.Gender = strings.ToLower(a.Gender)
	a.AlmaMaterStr = strings.ToLower(a.AlmaMaterStr)
	a.UndergraduateMajorStr = strings.ToLower(a.UndergraduateMajorStr)
	a.GraduateMajorStr = strings.ToLower(a.GraduateMajorStr)
	a.MilitaryBranch = strings.ToLower(a.MilitaryBranch)
	a.MilitaryRank = strings.ToLower(a.MilitaryRank)
	a.DeathMission = strings.ToLower(a.DeathMission)
	a.MissionStr = strings.ToLower(a.MissionStr)
}
