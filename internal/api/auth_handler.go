package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/clary-work01/hot-post-api/internal/domain/auth"
)

type AuthHandler struct {
	svc *auth.Service
}

func NewAuthHandler(svc *auth.Service) *AuthHandler {
	return &AuthHandler{svc}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.register)
}

// ---- DTO ----
type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *registerReq) validate() error {
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if l := len(req.Username); l < 3 || l > 50 {
		return errors.New("使用者名稱需為 3-50 字元")
	}
	if l := len(req.Password); l < 8 || l > 20 {
		return errors.New("密碼需為 8-20 字元")
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
