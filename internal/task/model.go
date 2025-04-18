package task

import "time"

type Task struct {
	ID          string    `gorm:"primaryKey"`
	Title       string
	Description string
	Completed   bool
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
