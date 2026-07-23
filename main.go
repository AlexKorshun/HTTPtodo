package main

import (
	"log"
	"net/http"
)

func main() {

	handler := Handler{&TaskService{storage: &FileStorage{fileName: "todos.json"}}}
	http.HandleFunc("GET /tasks", handler.getHandler)

	http.HandleFunc("POST /tasks", handler.postHandler)

	http.HandleFunc("PATCH /tasks/{id}", handler.patchHandler)

	http.HandleFunc("DELETE /tasks/{id}", handler.deleteHandler)

	log.Println("сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
