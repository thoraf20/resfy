package main

import (
	"fmt"
	// "log"
	"net/http"

	"github.com/go-chi/chi/v5"
	// "github.com/hibiken/asynq"
	"github.com/thoraf20/resfy/internal/db"
	"github.com/thoraf20/resfy/internal/task"
	// "github.com/thoraf20/resfy/internal/worker"
)

func main() {

	// redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379"}

	// taskProcessor := worker.NewTaskProcessor()
	// srv := asynq.NewServer(redisOpt, asynq.Config{Concurrency: 10})
	// mux := asynq.NewServeMux()
	// mux.HandleFunc(worker.TaskSendReminder, taskProcessor.HandleReminderTask)

	// go func() {
	// 	if err := srv.Run(mux); err != nil {
	// 		log.Fatalf("Could not run asynq server: %v", err)
	// 	}
	// }()

	r := chi.NewRouter()

	// Task setup
	dbConn := db.Connect()
	taskService := task.NewTaskService(dbConn)
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
