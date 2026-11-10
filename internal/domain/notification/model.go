package notification

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID          uuid.UUID `gorm:"default:uuidv7()"`
	RecipientID uuid.UUID
	ActorID     uuid.UUID
	Type        NotificationType
	PostID      *uuid.UUID
	ReadAt      *time.Time
	CreatedAt   time.Time
}

type NotificationType string

var (
	TypePostLiked    NotificationType = "post_liked"
	TypeNewPost      NotificationType = "new_post"
	TypeUserFollowed NotificationType = "user_followed"
)
