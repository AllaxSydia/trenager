package handler

import (
	"encoding/json"
	"net/http"

	"api-gateway/internal/client"

	authpb "github.com/AllaxSydia/trenager/proto/auth"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	clients *client.GRPCClients
}

func NewAuthHandler(clients *client.GRPCClients) *AuthHandler {
	return &AuthHandler{clients: clients}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.clients.Auth.Register(r.Context(), &authpb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			writeError(w, http.StatusBadRequest, st.Message())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.clients.Auth.Login(r.Context(), &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			writeError(w, http.StatusUnauthorized, st.Message())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
