package auth

import (
	"net/http"
	"time"

	"github.com/clary-work01/hot-post-api/internal/domain/user"
	"github.com/clary-work01/hot-post-api/internal/platform/httpx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ---- DTO ----
type registerReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required"`
}

type userResp struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func toUserResp(u *user.User) userResp {
	return userResp{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}
}

// POST /auth/register
func (h *Handler) register(c *gin.Context) {
	var req registerReq

	err := httpx.DecodeJSON(c, &req)
	if err != nil {
		httpx.InvalidInputError(c, err)
		return
	}

	u, err := h.svc.register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		writeError(c, err)
		return
	}

	httpx.JSON(c, http.StatusCreated, toUserResp(u))
}
