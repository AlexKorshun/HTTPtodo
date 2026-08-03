package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/AlexKorshun/HTTPtodo/internal/model"
)

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.List(r.Context())
	if err != nil {
		log.Printf("GET /tasks: %v", err)
		respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	respondJSON(w, http.StatusOK, tasks)
}

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	type bodyJSON struct {
		Text string `json:"text"`
	}
	body := bodyJSON{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "не удалось получить JSON")
		return
	}
	task, err := h.taskService.Add(r.Context(), body.Text)
	if err != nil {
		if errors.Is(err, model.ErrEmptyText) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("POST /tasks: %v", err)
		respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	respondJSON(w, http.StatusCreated, task)
}

func (h *Handler) PatchHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "некорректный ввод номера задачи")
		return
	}

	task, err := h.taskService.Change(r.Context(), id)

	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("PATCH /tasks: %v", err)
		respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	respondJSON(w, http.StatusOK, task)

}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "некорректный ввод номера задачи")
		return
	}

	if err = h.taskService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("DELETE /tasks: %v", err)
		respondError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
