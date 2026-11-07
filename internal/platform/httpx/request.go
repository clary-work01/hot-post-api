package httpx

import (
	"net/http"

	"github.com/clary-work01/hot-post-api/internal/platform/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxBodyBytes = 1 << 20 // 1MB
func DecodeJSON(c *gin.Context, dst any) error {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(maxBodyBytes))

	return c.ShouldBindJSON(dst) // 原樣回傳，交給 InvalidInput 判斷
}

// 讀 path 參數（路由要寫成 /users/:id）
type idURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func ParseID(c *gin.Context) (uuid.UUID, error) {
	var u idURI
	if err := c.ShouldBindUri(&u); err != nil {
		return uuid.Nil, err // validator.ValidationErrors，details 會是 {"id": "uuid"}
	}

	return uuid.Parse(u.ID) // 格式已經驗證過，這裡不會失敗
}

// 讀 query 分頁參數
type Page struct {
	Limit  int `form:"limit,default=10" binding:"min=1,max=100"`
	Offset int `form:"offset,default=0" binding:"min=0"`
}

func BindPage(c *gin.Context) (Page, error) {
	var p Page
	if err := c.ShouldBindQuery(&p); err != nil {
		return Page{}, err
	}
	return p, nil
}

// 用在一定要登入的路由；拿不到通常是路由忘了掛 Auth，保守回 401
func CurrentUser(c *gin.Context) (uuid.UUID, bool) {
	id, ok := reqctx.UserID(c.Request.Context())
	if !ok {
		Error(c, http.StatusUnauthorized, ErrorBody{Code: "unauthorized", Message: "未授權"})
	}
	return id, ok
}

// 用在掛 OptionalAuth 的路由；未登入回 uuid.Nil，不寫 response
func ViewerID(c *gin.Context) uuid.UUID {
	id, _ := reqctx.UserID(c.Request.Context())
	return id
}
