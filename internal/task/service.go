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
	return &TaskService{db: db}
}

func (s *TaskService) Create(userID, title, description, dueDate string) (*Task, error) {
	parsedDueDate, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return nil, err
	}

	task := &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		DueDate:     parsedDueDate,
		Completed:   false,
		UserID:      userID,
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetAllByUser(userID string) ([]Task, error) {
	var tasks []Task
	if err := s.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) Update(userID, taskID, title, description, dueDate string, completed bool) (*Task, error) {
	var task Task
	if err := s.db.First(&task, "id = ? AND user_id = ?", taskID, userID).Error; err != nil {
		return nil, err
	}

	parsedDueDate, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return nil, err
	}

	task.Title = title
	task.Description = description
	task.DueDate = parsedDueDate
	task.Completed = completed

	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) Delete(userID, taskID string) error {
	return s.db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&Task{}).Error
}
