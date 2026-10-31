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

type PostDetail struct {
	Post
	AuthorName string
	IsLiked    bool
}
