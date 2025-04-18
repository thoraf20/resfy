package task

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hibiken/asynq"
	"github.com/thoraf20/resfy/internal/worker"
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

	// Schedule task reminder
	distributor := worker.NewTaskDistributor(asynq.RedisClientOpt{Addr: "localhost:6379"})
	_ = distributor.ScheduleReminder(task.ID, task.DueDate)

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

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); 
	err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dueDate, err := time.Parse(time.RFC3339, body.DueDate)
	if err != nil {
		http.Error(w, "Invalid due_date format", http.StatusBadRequest)
		return
	}

	task, ok := h.Service.Update(id, body.Title, body.Description, dueDate)
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ok := h.Service.Delete(id)
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}