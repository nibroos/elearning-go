package models

import (
	"time"

	"gorm.io/gorm"
)

// TODO : add pools with educations

type Quiz struct {
	gorm.Model
	ID          uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"column:name"`
	Description string     `json:"description" gorm:"column:description"`
	Threshold   uint       `json:"threshold" gorm:"column:threshold"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}
