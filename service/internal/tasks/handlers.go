package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// HandleWelcomeEmailTask handler for welcome email task.
func HandleWelcomeEmailTask(c context.Context, t *asynq.Task) error {
	// Get user ID from given task.
	var payload struct {
		UserID int `json:"user_id"`
	}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	id := payload.UserID

	// Dummy message to the worker's output.
	fmt.Printf("Send Welcome Email to User ID %d\n", id)
	log.Printf("Processed welcome email task for User ID %d", id)

	return nil
}

// HandleReminderEmailTask for reminder email task.
func HandleReminderEmailTask(c context.Context, t *asynq.Task) error {
	// Get int with the user ID from the given task.
	var payload struct {
		UserID int    `json:"user_id"`
		SentIn string `json:"sent_in"`
	}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	id := payload.UserID
	time := payload.SentIn

	// Dummy message to the worker's output.
	fmt.Printf("Send Reminder Email to User ID %d\n", id)
	fmt.Printf("Reason: time is up (%v)\n", time)
	log.Printf("Processed reminder email task for User ID %d", id)

	return nil
}
