package model

import (
	"context"
	"time"
)

const (
	AdminUser = "admin"
	BaseUser  = "user"
)

type (
	User struct {
		ID        int       `json:"id"`
		FirstName string    `json:"firstName"`
		Surename  string    `json:"surename"`
		Email     string    `json:"email"`
		Password  string    `json:"password,omitempty"`
		Role      string    `json:"role"`
		ApiKey    string    `json:"apiKey"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	UserStore interface {
		Create(ctx context.Context, u *User) (int, error)
		List(ctx context.Context, limit, offset int) ([]*User, error)
		Get(ctx context.Context, id int) (*User, error)
		Update(ctx context.Context, u *User) error
		Delete(ctx context.Context, id int) error
		SearchApiKey(ctx context.Context, key string) (*User, error)
	}

	UserUsecase interface {
		Create(ctx context.Context, u *User) (*User, error)
		List(ctx context.Context, limit, offset int) ([]*User, error)
		Get(ctx context.Context, id int) (*User, error)
		Update(ctx context.Context, u *User) (*User, error)
		Delete(ctx context.Context, id int) error
	}
)
