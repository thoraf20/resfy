package task

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	// "github.com/hibiken/asynq"
	"github.com/thoraf20/resfy/internal/utils"
	// "github.com/thoraf20/resfy/internal/worker"
)

type Handler struct {
	Service *TaskService
}

func NewHandler(service *TaskService) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID      string `json:"user_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	dueDate, err := time.Parse(time.RFC3339, body.DueDate)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid due_date format")
		return
	}

	task, err := h.Service.Create(body.Title, body.Description, body.UserID, dueDate.Local().String())
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Schedule task reminder
	// distributor := worker.NewTaskDistributor(asynq.RedisClientOpt{Addr: "localhost:6379"})
	// _ = distributor.ScheduleReminder(task.ID, task.DueDate)

	utils.JSON(w, http.StatusCreated, task)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	tasks, err := h.Service.GetAllByUser(userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, tasks)
}

// func (h *Handler) MarkTaskComplete(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	task, ok := h.Service.MarkAsCompleted(id)
// 	if !ok {
// 		utils.Error(w, http.StatusNotFound, "Task not found")
// 		return
// 	}
// 	utils.JSON(w, http.StatusOK, task)
// }

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
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	dueDate, err := time.Parse(time.RFC3339, body.DueDate)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid due_date format")
		return
	}

	task, err := h.Service.Update(userID, id, body.Title, body.Description, dueDate.Local().String(), body.Completed)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "Task not found")
		return
	}
	utils.JSON(w, http.StatusOK, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("userID").(string)

	err := h.Service.Delete(userID, id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "Task not found")
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"message": "Task deleted"})
}
