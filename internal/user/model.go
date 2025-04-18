package user

import "time"

type User struct {
	ID        string `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
