package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type TaskProcessor struct{}

func NewTaskProcessor() *TaskProcessor {
	return &TaskProcessor{}
}

func (p *TaskProcessor) HandleReminderTask(ctx context.Context, t *asynq.Task) error {
	var payload PayloadSendReminder
	if err := json.Unmarshal(t.Payload(), &payload); 
	err != nil {
		return err
	}

	// TODO: send reminder logic (log, email, notification, etc.)
	fmt.Printf("⏰ Reminder: Task %s is due soon!\n", payload.TaskID)

	return nil
}
