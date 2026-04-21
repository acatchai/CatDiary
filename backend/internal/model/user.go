package model

import "time"

// User 用户实体模型
type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null;comment:密码哈希值" json:"-"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;comment:邮箱" json:"email"`
	Avatar       string `gorm:"type:varchar(255);comment:头像URL" json:"avatar"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"-"`
}

// TableName 显式指定表名，防止GORM默认加复数变成users
func (User) TableName() string {
	return "user"
}
