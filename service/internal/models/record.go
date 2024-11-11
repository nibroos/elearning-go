package models

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	ID          uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	EducationID uint       `json:"education_id" gorm:"column:education_id"`
	UserID      uint       `json:"user_id" gorm:"column:user_id"`
	TimeSpent   string     `json:"ref_num" gorm:"column:time_spent"`
	LastSeen    *time.Time `json:"last_seen" gorm:"column:last_seen"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}
