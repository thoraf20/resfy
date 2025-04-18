package main

import (
	
	// "github.com/hibiken/asynq"
	"github.com/thoraf20/resfy/internal/db"
	"github.com/thoraf20/resfy/internal/task"
	"github.com/gin-gonic/gin"
	"github.com/thoraf20/resfy/internal/user"
	"github.com/thoraf20/resfy/internal/auth"
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

	r := gin.Default()

	// Connect to DB
	dbConn := db.Connect()

	// Initialize Services
	userService := user.NewService(dbConn)
	taskService := task.NewTaskService(dbConn)

	// Auth Routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", func(c *gin.Context) {
			var body struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			if err := c.BindJSON(&body); err != nil {
				c.JSON(400, gin.H{"error": "invalid request"})
				return
			}

			user, err := userService.Register(body.Email, body.Password)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			token, _ := auth.GenerateToken(user.ID)
			c.JSON(200, gin.H{"token": token})
		})

		authGroup.POST("/login", func(c *gin.Context) {
			var body struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			if err := c.BindJSON(&body); err != nil {
				c.JSON(400, gin.H{"error": "invalid request"})
				return
			}

			user, err := userService.Authenticate(body.Email, body.Password)
			if err != nil {
				c.JSON(401, gin.H{"error": err.Error()})
				return
			}

			token, _ := auth.GenerateToken(user.ID)
			c.JSON(200, gin.H{"token": token})
		})
	}

	// Protected Task Routes
	taskGroup := r.Group("/tasks")
	taskGroup.Use(auth.Middleware())
	{
		taskGroup.GET("/", func(c *gin.Context) {
			userID := c.GetString("userID")
			tasks, err := taskService.GetAllByUser(userID)
			if err != nil {
				c.JSON(400, gin.H{"error": "invalid input"})
				return
			}
			c.JSON(200, tasks)
		})

		taskGroup.POST("/", func(c *gin.Context) {
			userID := c.GetString("userID")
			var body struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				DueDate     string `json:"due_date"` // ISO8601 format expected
			}
			if err := c.BindJSON(&body); err != nil {
				c.JSON(400, gin.H{"error": "invalid input"})
				return
			}

			newTask, err := taskService.Create(userID, body.Title, body.Description, body.DueDate)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(201, newTask)
		})

		taskGroup.PUT("/:id", func(c *gin.Context) {
			userID := c.GetString("userID")
			taskID := c.Param("id")

			var body struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				DueDate     string `json:"due_date"`
				Completed   bool   `json:"completed"`
			}
			if err := c.BindJSON(&body); err != nil {
				c.JSON(400, gin.H{"error": "invalid input"})
				return
			}

			updatedTask, err := taskService.Update(userID, taskID, body.Title, body.Description, body.DueDate, body.Completed)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, updatedTask)
		})

		taskGroup.DELETE("/:id", func(c *gin.Context) {
			userID := c.GetString("userID")
			taskID := c.Param("id")

			err := taskService.Delete(userID, taskID)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			c.Status(204)
		})
	}

	r.Run(":8080")
}

