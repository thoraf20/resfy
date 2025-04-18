package task

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hibiken/asynq"
	"github.com/thoraf20/resfy/internal/worker"
)

type Handler struct {
	Service     *TaskService // Updated to use GORM-based Service
	Distributor *worker.TaskDistributor
}

func NewHandler(service *TaskService) *Handler {
	return &Handler{
		Service:     service,
		Distributor: worker.NewTaskDistributor(asynq.RedisClientOpt{Addr: "localhost:6379"}),
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"` // as ISO string
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(string)

	task, err := h.Service.Create(userID, body.Title, body.Description, body.DueDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Schedule reminder
	_ = h.Distributor.ScheduleReminder(task.ID, task.DueDate)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	tasks, err := h.Service.GetAllByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("userID").(string)

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
		Completed   bool   `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.Service.Update(userID, id, body.Title, body.Description, body.DueDate, body.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("userID").(string)

	err := h.Service.Delete(userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
