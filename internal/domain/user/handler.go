package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r gin.IRouter, requireAuth, optionalAuth gin.HandlerFunc) {
	g := r.Group("/users")

	// 訪客可看；有登入時會帶 is_following
	pub := g.Group("", optionalAuth)
	pub.GET("/:id/followers", h.followers)
	pub.GET("/:id/followees", h.followees)
	pub.GET("/:id", h.profile)

	// 一定要登入
	priv := g.Group("", requireAuth)
	priv.PUT("/:id/follow", h.follow)
	priv.DELETE("/:id/follow", h.unfollow)

}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		httpx.Error(c, http.StatusNotFound, httpx.ErrorBody{Code: "user_not_found", Message: "用戶不存在"})
	case errors.Is(err, ErrUsernameTaken):
		httpx.Error(c, http.StatusConflict, httpx.ErrorBody{Code: "username_taken", Message: "用戶名已存在"})
	case errors.Is(err, errCannotFollowSelf):
		httpx.Error(c, http.StatusBadRequest, httpx.ErrorBody{Code: "cannot_follow_self", Message: "不能關注自己"})
	default:
		httpx.ServerError(c, err)
	}
}

// ---- DTO ----

type FollowItemResp struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	FollowedAt  time.Time `json:"followed_at"`
	IsFollowing bool      `json:"is_following"`
}

type FollowListResp struct {
	List []FollowItemResp `json:"list"`
}

func toFollowListResp(items []FollowItem) *FollowListResp {
	followItems := make([]FollowItemResp, len(items))
	for i, item := range items {
		f := FollowItemResp{
			ID:          item.ID,
			Username:    item.Username,
			FollowedAt:  item.FollowedAt,
			IsFollowing: item.IsFollowing,
		}
		followItems[i] = f
	}

	return &FollowListResp{List: followItems}
}
