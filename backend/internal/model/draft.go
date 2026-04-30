package model

import (
	"time"

	"gorm.io/gorm"
)

type DraftDiary struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement:false" json:"id"`
	UserID   uint   `gorm:"index;not null" json:"user_id"`
	Title    string `gorm:"type:varchar(100);not null" json:"title"`
	Content  string `gorm:"type:longtext;not null" json:"content"`
	Mood     string `gorm:"type:varchar(20)" json:"mood"`
	Weather  string `gorm:"type:varchar(20)" json:"weather"`
	Location string `gorm:"type:varchar(100)" json:"location"`

	Version   uint64         `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DraftDiary) TableName() string {
	return "draft_diary"
}
