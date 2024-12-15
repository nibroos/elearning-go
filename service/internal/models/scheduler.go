package models

import (
	"time"
)

type Scheduler struct {
	ID          uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"column:name"`
	Description string     `json:"description" gorm:"column:description"`
	Cron        string     `json:"cron" gorm:"column:cron"`
	Payload     string     `json:"payload" gorm:"column:payload"`
	Status      string     `json:"status" gorm:"column:status"`
	EntryID     int        `json:"entry_id" gorm:"column:entry_id"`
	StartAt     time.Time  `json:"start_at" gorm:"column:start_at"`
	EndAt       *time.Time `json:"end_at" gorm:"column:end_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}
