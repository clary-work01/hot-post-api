package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/clary-work01/hot-post-api/internal/domain/auth"
	"github.com/google/uuid"
)

const refreshCookieName = "refresh_token"

type AuthHandler struct {
	svc            *auth.Service
	isCookieSecure bool
}

func NewAuthHandler(svc *auth.Service, isCookieSecure bool) *AuthHandler {
	return &AuthHandler{
		svc:            svc,
		isCookieSecure: isCookieSecure,
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.register)
	mux.HandleFunc("POST /auth/login", h.login)
	mux.HandleFunc("POST /auth/refresh", h.refresh)
}

// ---- DTO ----
type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *registerReq) validate() error {
	req.Username = strings.TrimSpace(req.Username)

	if l := len(req.Username); l < 3 || l > 50 {
		return errors.New("使用者名稱需為 3-50 字元")
	}
	if l := len(req.Password); l < 8 || l > 20 {
		return errors.New("密碼需為 8-20 字元")
	}
	return nil
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *loginReq) validate() error {
	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" || req.Password == "" {
		return errors.New("請輸入帳密")
	}
	return nil
}

type userResp struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type tokenResp struct {
	AccessToken string `json:"access_token"`
}

// POST /auth/register
func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var req registerReq

	err := decodeJSON(w, r, &req)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數錯誤")
		return
	}

	if err := req.validate(); err != nil {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, err.Error())
		return
	}

	u, err := h.svc.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusCreated, "register success", userResp{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	})
}

// POST /auth/login
func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := decodeJSON(w, r, &req); err != nil {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, "參數錯誤")
		return
	}

	if err := req.validate(); err != nil {
		errorJSON(w, http.StatusBadRequest, codeInvalidInput, err.Error())
		return
	}

	res, err := h.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, r, err)
		return
	}

	h.setRefreshCookie(w, res.RefreshToken)

	JSON(w, http.StatusOK, "login success", tokenResp{
		AccessToken: res.AccessToken,
	})
}

// POST /auth/refresh
func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshCookieName)
	if err != nil {
		writeError(w, r, auth.ErrInvalidRefreshToken)
		return
	}

	res, err := h.svc.Refresh(r.Context(), cookie.Value)
	if errors.Is(err, auth.ErrInvalidRefreshToken) {
		h.clearRefreshCookie(w)
	}
	if err != nil {
		writeError(w, r, err)
		return
	}

	h.setRefreshCookie(w, res.RefreshToken)

	JSON(w, http.StatusOK, "refresh success", tokenResp{
		AccessToken: res.AccessToken,
	})
}

func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:  refreshCookieName,
		Value: refreshToken,
		// 只有 /auth/... 開頭的請求會帶上
		Path: "/auth",
		// 存活秒數，7 天後瀏覽器自動刪除
		MaxAge: int(auth.RefreshTTL.Seconds()),
		// JS 讀不到
		HttpOnly: true,
		// 只在 HTTPS 連線時才送出
		Secure: h.isCookieSecure,
		// 只有從你自己網站發出的請求才帶
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/auth",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.isCookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}
