package worker

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type TaskDistributor struct {
	client *asynq.Client
}

func NewTaskDistributor(redisOpt asynq.RedisClientOpt) *TaskDistributor {
	return &TaskDistributor{
		client: asynq.NewClient(redisOpt),
	}
}

func (d *TaskDistributor) ScheduleReminder(taskID string, when time.Time) error {
	payload, err := json.Marshal(PayloadSendReminder{TaskID: taskID})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskSendReminder, payload)

	// Schedule the task to run 30 minutes before due date
	_, err = d.client.Enqueue(task, asynq.ProcessAt(when.Add(-30*time.Minute)))
	return err
}
