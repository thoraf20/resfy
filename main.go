package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thoraf20/resfy/internal/task"
)

func main() {
	r := chi.NewRouter()

	// Task setup
	taskService := task.NewTaskService()
	taskHandler := task.NewHandler(taskService)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Resfy!"))
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", taskHandler.CreateTask)
		r.Get("/", taskHandler.GetTasks)
		r.Put("/{id}", taskHandler.UpdateTask)
		r.Put("/{id}/complete", taskHandler.MarkTaskComplete)
		r.Delete("/{id}", taskHandler.DeleteTask)
	})

	fmt.Println("Server running at :8080")
	http.ListenAndServe(":8080", r)
}
