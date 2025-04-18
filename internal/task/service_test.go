package task_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/thoraf20/resfy/internal/task"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&task.Task{}))
	return db
}

func createTestTask(service *task.TaskService, userID string) *task.Task {
	task, _ := service.Create("Test Task", "Test Desc", userID, time.Now().Add(24*time.Hour).Local().String())
	return task
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	service := task.NewTaskService(db)

	userID := "user123"
	title := "My Task"
	description := "My Description"
	due := time.Now().Add(48 * time.Hour)

	tk, err := service.Create(title, description, userID, due.Local().String())
	require.NoError(t, err)

	assert.Equal(t, title, tk.Title)
	assert.Equal(t, userID, tk.UserID)
}

func TestGetAllByUser(t *testing.T) {
	db := setupTestDB(t)
	service := task.NewTaskService(db)

	userID := "user1"
	_ = createTestTask(service, userID)
	_ = createTestTask(service, userID)

	tasks, err := service.GetAllByUser(userID)
	require.NoError(t, err)
	assert.Len(t, tasks, 2)
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	service := task.NewTaskService(db)

	original := createTestTask(service, "user1")
	newTitle := "Updated"
	newDesc := "Updated Desc"
	newDue := time.Now().Add(72 * time.Hour)

	updated, ok := service.Update(original.ID.String(), newTitle, newDesc, newDue.Local().String())
	assert.True(t, ok)
	assert.Equal(t, newTitle, updated.Title)
	assert.Equal(t, newDesc, updated.Description)
}

func TestMarkAsCompleted(t *testing.T) {
	db := setupTestDB(t)
	service := task.NewTaskService(db)

	task := createTestTask(service, "user1")
	completed, ok := service.Update(task.ID.String())
	assert.True(t, ok)
	assert.True(t, completed.Completed)
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	service := task.NewTaskService(db)

	task := createTestTask(service, "user1")
	ok := service.Delete(task.ID.String())
	assert.True(t, ok)

	tasks, err := service.GetAllByUser("user1")
	require.NoError(t, err)
	assert.Len(t, tasks, 0)
}
