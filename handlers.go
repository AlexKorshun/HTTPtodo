package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	taskService *TaskService
}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.List()
	if err != nil {
		log.Printf("GET /tasks: %v", err)
		respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	respondJSON(w, http.StatusOK, tasks)
}

func (h *Handler) postHandler(w http.ResponseWriter, r *http.Request) {
	type bodyJSON struct {
		Text string `json:"text"`
	}
	body := bodyJSON{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "не удалось получить JSON")
		return
	}
	task, err := h.taskService.Add(body.Text)
	if err != nil {
		if errors.Is(err, ErrEmptyText) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("POST /tasks: %v", err)
		respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	respondJSON(w, http.StatusCreated, task)
}

func (h *Handler) patchHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "некорректный ввод номера задачи")
		return
	}

	task, err := h.taskService.Change(id)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("PATCH /tasks: %v", err)
		respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	respondJSON(w, http.StatusOK, task)

}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "некорректный ввод номера задачи")
		return
	}

	if err = h.taskService.Delete(id); err != nil {
		if errors.Is(err, ErrNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("DELETE /tasks: %v", err)
		respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
