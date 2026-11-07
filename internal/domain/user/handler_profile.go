package user

import (
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileResp struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	PostCount      int       `json:"post_count"`
	FollowerCount  int       `json:"follower_count"`
	FollowingCount int       `json:"following_count"`
	IsFollowing    bool      `json:"is_following"`
}

func toProfileResp(p *Profile) *ProfileResp {
	return &ProfileResp{
		ID:             p.ID,
		Username:       p.Username,
		PostCount:      p.PostCount,
		FollowerCount:  p.FollowerCount,
		FollowingCount: p.FollowingCount,
		IsFollowing:    p.IsFollowing,
	}
}

// GET /users/:id
func (h *Handler) profile(c *gin.Context) {
	userID, err := httpx.ParseID(c)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	viewerID := httpx.ViewerID(c) // 訪客為 uuid.Nil

	profile, err := h.svc.profile(c.Request.Context(), userID, viewerID)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.JSON(c, http.StatusOK, toProfileResp(profile))
}
