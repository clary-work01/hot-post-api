package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/clary-work01/hot-post-api/internal/domain/auth"
)

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

// POST /auth/register
func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var req registerReq

	err := decodeJSON(w, r, &req)
	if err != nil {
		JSON(w, http.StatusBadRequest, "參數錯誤", nil)
		return
	}

	if err := req.validate(); err != nil {
		JSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	u, err := h.svc.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, r, err)
		return
	}

	JSON(w, http.StatusCreated, "register success", u.Username)
}

// POST /auth/login
func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := decodeJSON(w, r, &req); err != nil {
		JSON(w, http.StatusBadRequest, "參數錯誤", nil)
		return
	}

	if err := req.validate(); err != nil {
		JSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	res, err := h.svc.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, r, fmt.Errorf("login: %w", err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: res.RefreshToken,
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

	JSON(w, http.StatusOK, "login success", map[string]any{
		"access_token": res.AccessToken,
	})
}
