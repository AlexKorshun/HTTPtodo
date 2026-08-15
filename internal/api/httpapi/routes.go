package httpapi

import (
	"io/fs"
	"net/http"
)

// NewRouter собирает маршруты: статика и вход открыты, задачи — только с токеном.
func NewRouter(handler *Handler, secret string, appID int, static fs.FS) http.Handler {
	mux := http.NewServeMux()

	auth := AuthMiddleware(secret, appID)
	protected := func(h http.HandlerFunc) http.Handler { return auth(h) }

	mux.Handle("GET /tasks", protected(handler.GetHandler))
	mux.Handle("POST /tasks", protected(handler.PostHandler))
	mux.Handle("PATCH /tasks/{id}", protected(handler.PatchHandler))
	mux.Handle("DELETE /tasks/{id}", protected(handler.DeleteHandler))

	mux.HandleFunc("POST /auth/register", handler.RegisterHandler)
	mux.HandleFunc("POST /auth/login", handler.LoginHandler)

	mux.Handle("GET /", http.FileServerFS(static))

	return mux
}
