package task

import (
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	tasks map[string]Task
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[string]Task),
	}
}

func (s *TaskService) Create(title, description string, dueDate time.Time) Task {
	id := uuid.New().String()
	task := Task{
		ID:          id,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	s.tasks[id] = task
	return task
}

func (s *TaskService) GetAll() []Task {
	tasks := []Task{}
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}
