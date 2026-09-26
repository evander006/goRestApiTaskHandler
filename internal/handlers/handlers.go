package handlers

import (
	"encoding/json"
	"net/http"
	"restApiGo/internal/database"
	"strconv"
	"strings"
)

type Handlers struct {
	store database.TaskStore
}

func NewHandlers(store database.TaskStore) *Handlers {
	return &Handlers{
		store: store,
	}
}
func respondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	respondWithJson(w, statusCode, map[string]string{"error": msg})
}

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting all tasks")
		return
	}
	respondWithJson(w, http.StatusOK, tasks)
}
func (h *Handlers) GetTaskById(id int, w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	task, err := h.store.GetById(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting task")
		return
	}
	respondWithJson(w, http.StatusOK, task)
}
