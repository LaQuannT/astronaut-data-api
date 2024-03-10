package store

import (
	"context"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type astronautStore struct {
	db *pgxpool.Pool
}

func NewAstronautStore(db *pgxpool.Pool) *astronautStore {
	return &astronautStore{
		db: db,
	}
}

func (s *astronautStore) Create(ctx context.Context, a *model.Astronaut) (*model.Astronaut, error) {
	query := `INSERT INTO astronaut
  (name, year, "group", status, birth_date, birth_place, gender, alma_mater, undergraduate_major,
  graduate_major, military_rank, military_branch, space_flights, space_flight_hrs, space_walks,
  space_walk_hrs, missions, death_date, death_mission)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
  RETURNING id;`

	err := s.db.QueryRow(ctx, query, a.Name, a.Year, a.Group, a.Status, a.BirthDate, a.BirthPlace,
		a.Gender, pq.Array(a.AlmaMater), pq.Array(a.UndergraduateMajor), pq.Array(a.GraduateMajor), a.MilitaryRank, a.MilitaryBranch, a.SpaceFlights,
		a.SpaceFlightHours, a.SpaceWalks, a.SpaceWalkHours, pq.Array(a.Missions), a.DeathDate, a.DeathMission).Scan(&a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *astronautStore) List(ctx context.Context, limit, offset int) ([]*model.Astronaut, error) {
	var astronauts []*model.Astronaut

	query := `SELECT * FROM astronaut ORDER BY name ASC LIMIT $1 OFFSET $2;`
	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Next()

	for rows.Next() {
		a, err := fromRowToAstronaut(rows)
		if err != nil {
			return nil, err
		}
		astronauts = append(astronauts, a)
	}
	return astronauts, nil
}

func (s *astronautStore) Get(ctx context.Context, id int) (*model.Astronaut, error) {
	query := `SELECT * FROM astronaut WHERE id=$1;`
	rows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return fromRowToAstronaut(rows)
	}
	return nil, nil
}

func (s *astronautStore) Update(ctx context.Context, a *model.Astronaut) error {
	query := `UPDATE astronaut SET name=$1, year=$2, group=$3, status=$4, birth_date=$5, birth_place=$6, gender=$7, alma_mater=$8, undergraduate_major=$9,
  graduate_major=$10, military_rank=$11, military_branch=$12, space_flights=$13, space_flight_hrs=$14, space_walks=$15, space_walk_hrs=$16, missions=$17,
  death_date=$18, death_mission=$19 WHERE id=$20;`
	_, err := s.db.Exec(ctx, query, a.Name, a.Year, a.Group, a.Status, a.BirthDate, a.BirthPlace,
		a.Gender, pq.Array(a.AlmaMater), pq.Array(a.UndergraduateMajor), pq.Array(a.GraduateMajor), a.MilitaryRank, a.MilitaryBranch, a.SpaceFlights,
		a.SpaceFlightHours, a.SpaceWalks, a.SpaceWalkHours, pq.Array(a.Missions), a.DeathDate, a.DeathMission, a.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *astronautStore) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM astronaut WHERE id=$1;`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func fromRowToAstronaut(r pgx.Rows) (*model.Astronaut, error) {
	a := new(model.Astronaut)
	err := r.Scan(&a.ID, &a.Name, &a.Year, &a.Group, &a.Status, &a.BirthDate, &a.BirthPlace,
		&a.Gender, &a.AlmaMater, &a.UndergraduateMajor, &a.GraduateMajor, &a.MilitaryRank, &a.MilitaryBranch, &a.SpaceFlights,
		&a.SpaceFlightHours, &a.SpaceWalks, &a.SpaceWalkHours, &a.Missions, &a.DeathDate, &a.DeathMission)
	if err != nil {
		return nil, err
	}

	return a, nil
}
