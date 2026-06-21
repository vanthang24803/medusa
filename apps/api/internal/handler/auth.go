package handler

import (
	"net/http"

	"ecommerce/apps/api/internal/middleware"
	"ecommerce/modules/auth"
	"ecommerce/packages/httpx"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AuthHandler struct {
	svc auth.Service
	log *zap.Logger
}

func NewAuthHandler(svc auth.Service, log *zap.Logger) *AuthHandler {
	return &AuthHandler{svc: svc, log: log}
}

func (h *AuthHandler) Routes(r chi.Router) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.svc))
		r.Post("/logout", h.Logout)
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterReq
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	resp, err := h.svc.Register(r.Context(), req)
	if err != nil {
		h.log.Error("failed to register user", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginReq
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		h.log.Error("failed to decode login request", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		h.log.Error("failed to login user", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req auth.RefreshReq
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		h.log.Error("failed to decode refresh request", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	resp, err := h.svc.RefreshToken(r.Context(), req)
	if err != nil {
		h.log.Error("failed to refresh token", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authID := middleware.GetAuthIdentityID(r.Context())
	if err := h.svc.Logout(r.Context(), authID); err != nil {
		h.log.Error("failed to logout user", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "logged out"})
}
