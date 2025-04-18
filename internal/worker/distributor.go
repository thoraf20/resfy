package worker

import (
	"time"
	"github.com/hibiken/asynq"
	// "github.com/thoraf20/resfy/internal/worker/"
)

type TaskDistributor struct {
	client *asynq.Client
}

func NewTaskDistributor(redisOpts asynq.RedisClientOpt) *TaskDistributor {
	return &TaskDistributor{
		client: asynq.NewClient(redisOpts),
	}
}

func (d *TaskDistributor) ScheduleReminder(taskID string, dueDate time.Time, offset time.Duration) error {
	reminderTime := dueDate.Add(-offset)
	if reminderTime.Before(time.Now()) {
		return nil // too late to schedule
	}

	payload := &payload.PayloadSendReminder{TaskID: taskID}
	task := asynq.NewTask("send:reminder", payload.Marshal())

	_, err := d.client.Enqueue(task, asynq.ProcessAt(reminderTime))
	return err
}
