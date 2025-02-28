package store

import (
	"database/sql"
	"context"

)

type UsersStore struct{
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context) error{
	return nil
}