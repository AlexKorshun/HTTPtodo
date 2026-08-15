package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey struct{}

var userIDKey = ctxKey{}

func AuthMiddleware(secret string, appID int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || tokenStr == "" {
				respondError(w, http.StatusUnauthorized, "требуется токен")
				return
			}

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				respondError(w, http.StatusUnauthorized, "невалидный токен")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				respondError(w, http.StatusUnauthorized, "невалидный токен")
				return
			}

			uid, ok := claims["uid"].(float64)
			if !ok {
				respondError(w, http.StatusUnauthorized, "в токене нет uid")
				return
			}

			tokenAppID, ok := claims["app_id"].(float64)
			if !ok || int(tokenAppID) != appID {
				respondError(w, http.StatusUnauthorized, "токен выдан для другого приложения")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, int64(uid))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func userIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}
