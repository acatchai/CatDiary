package model

import (
	"time"

	"gorm.io/gorm"
)

// Diary 日记实体模型
type Diary struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID   uint   `gorm:"index;not null;comment:关联的用户ID" json:"user_id"`
	Title    string `gorm:"type:varchar(100);not null;comment:日记标题" json:"title"`
	Content  string `gorm:"type:longtext;not null;comment:Markdown正文" json:"content"`
	Mood     string `gorm:"type:varchar(20);comment:当天心情(如: happy, sad)" json:"mood"`
	Weather  string `gorm:"type:varchar(20);comment:天气情况(如: sunny, rainy)" json:"weather"`
	Location string `gorm:"type:varchar(100);comment:记录地点" json:"location"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Diary) TableName() string {
	return "diary"
}
