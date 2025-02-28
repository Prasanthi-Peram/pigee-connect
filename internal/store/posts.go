package store

import (
	"database/sql"
	"context"

	"github.com/lib/pq"
)

type PostStore struct{
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context) error{
	
	return nil
}