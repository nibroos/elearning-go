package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/nibroos/elearning-go/service/internal/tasks"
	"github.com/robfig/cron/v3"
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

	// Create a new Asynq client.
	client := asynq.NewClient(redisConnection)
	defer client.Close()

	// Create a new cron job scheduler.
	c := cron.New()

	// Schedule a job to enqueue tasks every minute.
	c.AddFunc("@every 1m", func() {
		// Enqueue a welcome email task with a 10-second delay.
		welcomeTask := tasks.NewWelcomeEmailTask(42)
		info, err := client.Enqueue(welcomeTask, asynq.Queue("default"), asynq.ProcessIn(10*time.Second))
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		log.Printf("Enqueued welcome email task: id=%s queue=%s", info.ID, info.Queue)

		// Enqueue a reminder email task with a 10-second delay.
		reminderTask := tasks.NewReminderEmailTask(42, time.Now().Add(24*time.Hour))
		info, err = client.Enqueue(reminderTask, asynq.Queue("default"), asynq.ProcessIn(10*time.Second))
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		log.Printf("Enqueued reminder email task: id=%s queue=%s", info.ID, info.Queue)

		log.Println("Tasks enqueued successfully")
	})

	// Start the cron scheduler.
	c.Start()

	// Keep the application running.
	select {}
}
