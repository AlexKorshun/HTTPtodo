package httpapi

import "net/http"

func NewRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", handler.GetHandler)
	mux.HandleFunc("POST /tasks", handler.PostHandler)
	mux.HandleFunc("PATCH /tasks/{id}", handler.PatchHandler)
	mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteHandler)
	return mux
}
