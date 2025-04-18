package task

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service *TaskService
}

func NewHandler(service *TaskService) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"` // as ISO string
	}

	if err := json.NewDecoder(r.Body).Decode(&body); 
	err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dueDate, err := time.Parse(time.RFC3339, body.DueDate)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	task := h.Service.Create(body.Title, body.Description, dueDate)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Service.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}


func (h *Handler) MarkTaskComplete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := h.Service.MarkAsCompleted(id)
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}