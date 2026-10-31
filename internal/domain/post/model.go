package post

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        uuid.UUID
	AuthorID  uuid.UUID
	Title     string
	Content   string
	LikeCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.Must(uuid.NewV7())
	}
	return nil
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
