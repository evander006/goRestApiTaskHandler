package handlers

import (
	"encoding/json"
	"net/http"
	"restApiGo/internal/database"
	"restApiGo/internal/models"
	"strconv"
	"strings"
)

type Handler struct {
	store *database.TaskStore
}

func NewHandler(store *database.TaskStore) *Handler {
	return &Handler{
		store: store,
	}
}
func respondWithJson(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	respondWithJson(w, statusCode, map[string]string{"error": msg})
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting all tasks")
		return
	}
	respondWithJson(w, http.StatusOK, tasks)
}
func (h *Handler) GetTaskById(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload: missing title")
		return
	}
	task, err := h.store.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating task")
		return
	}
	respondWithJson(w, http.StatusCreated, task)
}
func (h *Handler) DeleteTaskById(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}
	if err := h.store.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error deleting task")
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"msg": "success"})
}
func (h *Handler) UpdateTaskById(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID to update")
		return
	}
	var input models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload: missing title")
		return
	}
	task, err := h.store.Update(id, input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error updating task")
		return
	}
	respondWithJson(w, http.StatusOK, task)
}
