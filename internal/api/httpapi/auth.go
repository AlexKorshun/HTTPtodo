package httpapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SSOClient — то, что нужно HTTP-слою от сервиса авторизации.
// Реализуется клиентом internal/clients/sso/grpc.
type SSOClient interface {
	Register(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string, appID int32) (string, error)
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c credentials) valid() bool {
	return c.Email != "" && c.Password != ""
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	creds := credentials{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil || !creds.valid() {
		respondError(w, http.StatusBadRequest, "нужны email и пароль")
		return
	}

	userID, err := h.sso.Register(r.Context(), creds.Email, creds.Password)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			respondError(w, http.StatusConflict, "пользователь уже существует")
			return
		}
		log.Printf("POST /auth/register: %v", err)
		respondError(w, http.StatusInternalServerError, "не удалось зарегистрировать пользователя")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]int64{"user_id": userID})
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	creds := credentials{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil || !creds.valid() {
		respondError(w, http.StatusBadRequest, "нужны email и пароль")
		return
	}

	token, err := h.sso.Login(r.Context(), creds.Email, creds.Password, int32(h.appID))
	if err != nil {
		// sso отвечает InvalidArgument и на неверный пароль, и на несуществующего
		// пользователя — намеренно, чтобы не подсказывать, какая часть неверна.
		if status.Code(err) == codes.InvalidArgument {
			respondError(w, http.StatusUnauthorized, "неверный email или пароль")
			return
		}
		log.Printf("POST /auth/login: %v", err)
		respondError(w, http.StatusInternalServerError, "не удалось войти")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"token": token})
}
