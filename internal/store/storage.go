package store
import(
	"context"
	"database/sql"
	"errors"
	"time"
)

var(
	ErrNotFound= errors.New("Resource not found")
	ErrConflict= errors.New("Resource already exists")
	QueryTimeoutDuration= time.Second*5

)
type Storage struct{
	Posts interface{
		GetByID(context.Context,int64) (*Post,error)
		Create(context.Context, *Post) error
		Delete(context.Context,int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64,PaginatedFeedQuery) ([]PostWithMetadata,error)
	}

	Users interface{
		GetByID(context.Context, int64) (*User,error)
		Create(context.Context, *User) error
		CreateandInvite(ctx context.Context, user *User, token string) error
	}

	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, userID, followerID int64) error
		Unfollow(ctx context.Context, followerID, userID int64) error
	}
	
}

func NewStorage(db *sql.DB) Storage{
	return Storage{
		Posts: &PostStore{db},
		Users: &UsersStore{db},
		Comments: &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}