package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/nibroos/elearning-go/service/internal/tasks"
)

func main() {
	// Create and configure Redis connection.
	redisConnection := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: func() int {
			db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
			if err != nil {
				log.Fatalf("Invalid REDIS_DB value: %v", err)
			}
			return db
		}(),
	}

	// Create and configure Asynq worker server.
	worker := asynq.NewServer(redisConnection, asynq.Config{
		// Specify how many concurrent workers to use.
		Concurrency: 10,
		// Specify multiple queues with different priority.
		Queues: map[string]int{
			"critical": 6, // processed 60% of the time
			"default":  3, // processed 30% of the time
			"low":      1, // processed 10% of the time
		},
	})

	// Create a new task's mux instance.
	mux := asynq.NewServeMux()

	// Define a task handler for the welcome email task.
	mux.HandleFunc(
		tasks.TypeWelcomeEmail,       // task type
		tasks.HandleWelcomeEmailTask, // handler function
	)

	// Define a task handler for the reminder email task.
	mux.HandleFunc(
		tasks.TypeReminderEmail,       // task type
		tasks.HandleReminderEmailTask, // handler function
	)

	// Run worker server.
	if err := worker.Run(mux); err != nil {
		log.Fatal(err)
	}
}
