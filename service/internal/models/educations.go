package models

import (
	"gorm.io/gorm"
)

type Education struct {
	gorm.Model
	ID            uint   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ModuleID      uint   `json:"module_id" gorm:"column:module_id"`
	NoUrut        uint   `json:"no_urut" gorm:"column:no_urut"`
	Name          string `json:"name" gorm:"column:name"`
	Description   string `json:"desc" gorm:"column:description"`
	TextMateri    string `json:"text_materi" gorm:"column:text_materi"`
	VideoURL      string `json:"video_url" gorm:"column:video_url"`
	ThumbnailURL  string `json:"thumbnail_url" gorm:"column:thumbnail_url"`
	AttachmentURL string `json:"attachment_url" gorm:"column:attachment_url"`
	CreatedByID   *uint  `json:"created_by_id" gorm:"column:created_by_id"`
	UpdatedByID   *uint  `json:"updated_by_id" gorm:"column:updated_by_id"`
	CreatedAt     string `json:"created_at" gorm:"column:created_at"`
}
