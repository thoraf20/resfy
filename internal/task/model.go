package task

import (
	"time"
)

type Task struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Completed   bool      `json:"completed"`
	UserID      string    `gorm:"type:uuid" json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
