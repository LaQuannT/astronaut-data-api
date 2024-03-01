package store

import (
	"context"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) Create(ctx context.Context, u *model.User) (int, error) {
	var id int

	query := ` INSERT INTO "user" (first_name, surname, email, password, api_key, role, created_at, updated_at)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`

	err := s.db.QueryRow(ctx, query, u.FirstName, u.Surename, u.Email, u.Password, u.ApiKey, u.Role, u.CreatedAt, u.UpdatedAt).Scan(&id)
	if err != nil {
		// check for postgres unique constraint error (code = 23505)
		return 0, err
	}
	return id, nil
}

func (s *UserStore) List(ctx context.Context, limt, offset int) ([]*model.User, error) {
	users := make([]*model.User, 0)

	query := `SELECT * FROM "user" ORDER BY surname ASC LIMIT $1 OFFSET $2;`
	rows, err := s.db.Query(ctx, query, limt, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u, err := fromRowToUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *UserStore) Get(ctx context.Context, id int) (*model.User, error) {
	u := new(model.User)
	query := `SELECT * FROM "user" WHERE id=$1;`

	if err := s.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.FirstName, &u.Surename, &u.Email, &u.Password, &u.ApiKey, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserStore) Update(ctx context.Context, u *model.User) error {
	query := `UPDATE "user" SET first_name=$1, surname=$2, email=$3, password=$4 api_key=$5, role=$6, updated_at=$7
  WHERE id=$8;`

	_, err := s.db.Exec(ctx, query, u.FirstName, u.Surename, u.Email, u.Password, u.ApiKey, u.Role, u.UpdatedAt, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "user" WHERE id=$1;`

	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func fromRowToUser(r pgx.Rows) (*model.User, error) {
	u := new(model.User)

	err := r.Scan(&u.ID, &u.FirstName, &u.Surename, &u.Email, &u.Password, &u.ApiKey, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
