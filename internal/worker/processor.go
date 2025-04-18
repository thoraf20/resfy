// internal/worker/processor.go

package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

func NewTaskProcessor() *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			Concurrency: 10,
		},
	)
}

func HandleReminderTask(ctx context.Context, t *asynq.Task) error {
	var p PayloadSendReminder
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	// TODO: send email/notification here
	fmt.Printf("📬 Reminder: Task %s is due soon!\n", p.TaskID)
	return nil
}
