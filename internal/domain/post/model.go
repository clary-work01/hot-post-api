package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"default:uuidv7()"`
	AuthorID  uuid.UUID
	Title     string
	Content   string
	LikeCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Detail struct {
	Post
	AuthorName string
	IsLiked    bool
}

// 不回content內容 回excerpt摘要
type Summary struct {
	ID         uuid.UUID
	AuthorID   uuid.UUID
	AuthorName string
	Title      string
	Excerpt    string
	LikeCount  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsLiked    bool
}

type PostLike struct {
	PostID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}
