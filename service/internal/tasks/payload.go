package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	// TypeWelcomeEmail is a name of the task type
	// for sending a welcome email.
	TypeWelcomeEmail = "email:welcome"

	// TypeReminderEmail is a name of the task type
	// for sending a reminder email.
	TypeReminderEmail = "email:reminder"
)

// NewWelcomeEmailTask task payload for a new welcome email.
func NewWelcomeEmailTask(id int) *asynq.Task {
	// Specify task payload.
	payload := map[string]interface{}{
		"user_id": id, // set user ID
	}

	// Marshal the payload to JSON.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		// Handle error.
		return nil
	}

	// Return a new task with given type and payload.
	return asynq.NewTask(TypeWelcomeEmail, payloadBytes)
}

// NewReminderEmailTask task payload for a reminder email.
func NewReminderEmailTask(id int, ts time.Time) *asynq.Task {
	// Specify task payload.
	payload := map[string]interface{}{
		"user_id": id,          // set user ID
		"sent_in": ts.String(), // set time to sending
	}

	// Marshal the payload to JSON.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		// Handle error.
		return nil
	}

	// Return a new task with given type and payload.
	return asynq.NewTask(TypeReminderEmail, payloadBytes, asynq.MaxRetry(5), asynq.Timeout(1*time.Minute))
}
