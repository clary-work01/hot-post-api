package like

import (
	"time"

	"github.com/google/uuid"
)

type PostLike struct {
	PostID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}
