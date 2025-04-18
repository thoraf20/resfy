package task

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	db.AutoMigrate(&Task{})
	return &TaskService{db: db}
}

func (s *TaskService) Create(title, description string, dueDate time.Time) Task {
	task := Task{
		ID:          uuid.NewString(),
		Title:       title,
		Description: description,
		DueDate:     dueDate,
	}
	s.db.Create(&task)
	return task
}

func (s *TaskService) GetAll() []Task {
	var tasks []Task
	s.db.Order("created_at desc").Find(&tasks)
	return tasks
}

func (s *TaskService) MarkAsCompleted(id string) (Task, bool) {
	var task Task
	if err := s.db.First(&task, "id = ?", id).Error; 
	err != nil {
		return Task{}, false
	}
	task.Completed = true
	s.db.Save(&task)
	return task, true
}

func (s *TaskService) Update(id, title, description string, dueDate time.Time) (Task, bool) {
	var task Task
	if err := s.db.First(&task, "id = ?", id).Error; err != nil {
		return Task{}, false
	}
	task.Title = title
	task.Description = description
	task.DueDate = dueDate
	s.db.Save(&task)
	return task, true
}

func (s *TaskService) Delete(id string) bool {
	result := s.db.Delete(&Task{}, "id = ?", id)
	return result.RowsAffected > 0
}
