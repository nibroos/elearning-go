package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/scheduler"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type SchedulerController struct {
	Cron  *cron.Cron
	DB    *gorm.DB
	SqlDB *sqlx.DB
}

type ScheduleRequest struct {
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at,omitempty"`
	Cron        string     `json:"cron"`
	Name        string     `json:"name"`
	Action      string     `json:"action"`
	Description string     `json:"description"`
}

type ScheduleStopRequest struct {
	Name string `json:"name"`
}

var availableProcesses = map[string]func(){
	"generate_random_string": scheduler.GenerateRandomString,
	"generate_random_number": scheduler.GenerateRandomNumber,
}

func NewSchedulerController(cron *cron.Cron, db *gorm.DB, sqlDB *sqlx.DB) *SchedulerController {
	return &SchedulerController{Cron: cron, DB: db, SqlDB: sqlDB}
}

func (sc *SchedulerController) Schedule(c *fiber.Ctx) error {
	var req ScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": "Invalid request"})
	}

	switch req.Action {
	case "start":
		return sc.startTask(c, req)
	case "stop":
		return sc.stopTask(c, req)
	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid action", "message": "Action must be 'start' or 'stop'"})
	}
}

func (sc *SchedulerController) startTask(c *fiber.Ctx, req ScheduleRequest) error {
	// Validate process name
	task, exists := availableProcesses[req.Name]
	if !exists {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid process name", "message": "Process not found"})
	}

	// Check if the process is already running
	var scheduler models.Scheduler
	if err := sc.SqlDB.Get(&scheduler, "SELECT * FROM schedulers WHERE name = $1 AND status = $2", req.Name, "running"); err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Process is already running", "message": "Process is already running"})
	}

	// Schedule the task
	entryID, err := sc.Cron.AddFunc(req.Cron, func() {
		task()
		if req.EndAt != nil && time.Now().After(*req.EndAt) {
			sc.StopCron(req.Name)
		}
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "Failed to schedule task"})
	}

	payload := map[string]interface{}{
		"name":     req.Name,
		"cron":     req.Cron,
		"action":   req.Action,
		"start_at": req.StartAt,
		"end_at":   req.EndAt,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "Failed to marshal payload"})
	}

	// Save the scheduler to the database
	newScheduler := models.Scheduler{
		Name:        req.Name,
		Cron:        req.Cron,
		Status:      "running",
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		Description: req.Description,
		Payload:     string(payloadJSON),
		EntryID:     int(entryID), // Store the EntryID
	}

	if err := sc.DB.Create(&newScheduler).Error; err != nil {
		sc.Cron.Remove(entryID)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "Failed to save scheduler"})
	}

	return c.JSON(fiber.Map{"message": "Task scheduled successfully"})
}

// func (sc *SchedulerController) Schedule(c *fiber.Ctx) error {
// 	var req ScheduleRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": "Invalid request"})
// 	}

// 	// Validate process name
// 	task, exists := availableProcesses[req.Name]
// 	if !exists {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid process name", "message": "Process not found"})
// 	}

// 	// Check if the process is already running
// 	var scheduler models.Scheduler
// 	if err := sc.SqlDB.Get(&scheduler, "SELECT * FROM schedulers WHERE name = $1 AND status = $2", req.Name, "running"); err == nil {
// 		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Process is already running", "message": "Process is already running"})
// 	}

// 	// Schedule the task
// 	entryID, err := sc.Cron.AddFunc(req.Cron, func() {
// 		task()
// 		if req.EndAt != nil && time.Now().After(*req.EndAt) {
// 			sc.StopCron(req.Name)
// 		}
// 	})
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "Failed to schedule task"})
// 	}

// 	payload := map[string]interface{}{
// 		"name":     req.Name,
// 		"cron":     req.Cron,
// 		"action":   req.Action,
// 		"start_at": req.StartAt,
// 		"end_at":   req.EndAt,
// 	}

// 	payloadJSON, err := json.Marshal(payload)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "Failed to marshal payload"})
// 	}

// 	// Save the scheduler to the database
// 	newScheduler := models.Scheduler{
// 		Name:        req.Name,
// 		Cron:        req.Cron,
// 		Status:      "running",
// 		StartAt:     req.StartAt,
// 		EndAt:       req.EndAt,
// 		Description: req.Description,
// 		Payload:     string(payloadJSON),
// 		EntryID:     int(entryID), // Store the EntryID
// 	}

// 	if err := sc.DB.Create(&newScheduler).Error; err != nil {
// 		sc.Cron.Remove(entryID)
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "Failed to save scheduler"})
// 	}

// 	return c.JSON(fiber.Map{"message": "Task scheduled successfully"})
// }

func (sc *SchedulerController) StopCron(name string) error {
	var scheduler models.Scheduler
	if err := sc.DB.Where("name = ? AND status = ?", name, "running").First(&scheduler).Error; err != nil {
		log.Printf("Failed to find scheduler: %v", err)
		return err
	}

	// Use the stored EntryID to remove the cron job
	sc.Cron.Remove(cron.EntryID(scheduler.EntryID))
	scheduler.Status = "stopped"
	scheduler.UpdatedAt = time.Now()
	if err := sc.DB.Save(&scheduler).Error; err != nil {
		return err
	}

	return nil
}
func (sc *SchedulerController) ReloadSchedules() error {
	var schedulers []models.Scheduler
	if err := sc.DB.Where("status = ?", "running").Find(&schedulers).Error; err != nil {
		return err
	}

	var newSchedulers []models.Scheduler
	var deletedSchedulerIDs []uint
	for _, scheduler := range schedulers {
		// task, exists := availableProcesses[scheduler.Name]
		// if !exists {
		// 	continue
		// }

		// entryID, err := sc.Cron.AddFunc(scheduler.Cron, task)
		// if err != nil {
		// 	return err
		// }

		// add to deletedSchedulerIDs
		deletedSchedulerIDs = append(deletedSchedulerIDs, scheduler.ID)

		// nullify the ID
		scheduler.ID = 0

		// add to newSchedulers
		newSchedulers = append(newSchedulers, scheduler)

		// if err := sc.DB.Save(&scheduler).Error; err != nil {
		// 	return err
		// }
	}
	log.Printf("newSchedulers: %v", newSchedulers)
	log.Printf("deletedSchedulerIDs: %v", deletedSchedulerIDs)

	// Delete the old schedulers
	if len(deletedSchedulerIDs) > 0 {
		query, args, err := sqlx.In("DELETE FROM schedulers WHERE id IN (?)", deletedSchedulerIDs)
		if err != nil {
			log.Printf("Failed to construct delete query: %v", err)
			return err
		}
		query = sc.SqlDB.Rebind(query)
		if _, err := sc.SqlDB.Exec(query, args...); err != nil {
			log.Printf("Failed to delete old schedulers: %v", err)
			return err
		}
	}

	// Save the new schedulers
	if len(newSchedulers) > 0 {
		for _, scheduler := range newSchedulers {
			task, exists := availableProcesses[scheduler.Name]
			if !exists {
				continue
			}

			entryID, err := sc.Cron.AddFunc(scheduler.Cron, func() {
				task()
				if scheduler.EndAt != nil && time.Now().After(*scheduler.EndAt) {
					sc.StopCron(scheduler.Name)
				}
			})
			if err != nil {
				log.Printf("Failed to add task: %v", err)
				continue
			}

			scheduler.EntryID = int(entryID)

			if err := sc.DB.Create(&scheduler).Error; err != nil {
				log.Printf("Failed to save new schedulers: %v", err)
			}
		}
	}

	return nil
}

func (sc *SchedulerController) stopTask(c *fiber.Ctx, req ScheduleRequest) error {
	if err := sc.StopCron(req.Name); err != nil {
		log.Printf("Failed to stop task: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "Failed to stop task"})
	}

	return c.JSON(fiber.Map{"message": "Task stopped successfully"})
}

// list of schedules
func (sc *SchedulerController) ListSchedules(c *fiber.Ctx) error {
	entries := sc.Cron.Entries()
	return c.JSON(fiber.Map{"schedules": entries})
}
