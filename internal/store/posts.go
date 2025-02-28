package store

import (
	"database/sql"
	"context"

	"github.com/lib/pq"
)

type Post struct{
	ID int64 `json:"id"`
	Content string `json:"content"`
	Title string `json:"title"`
	UserID int64 `json:"user_id"`
	Tags []string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

}

type PostStore struct{
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error{
	
	return nil
}