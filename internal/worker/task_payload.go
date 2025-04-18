package worker

const TaskSendReminder = "task:send_reminder"

type PayloadSendReminder struct {
	TaskID string `json:"task_id"`
}
