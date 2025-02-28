package store

import (
	"database/sql"
	"context"

)

type UsersStore struct{
	db *sql.DB
}

type User struct{
	ID int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt string `json:"created_at"`
}

func (s *UsersStore) Create(ctx context.Context, user *User) error{
	return nil
}