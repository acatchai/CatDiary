package repository

import (
	"errors"

	"github.com/acatchai/catdiary/backend/internal/model"
	"gorm.io/gorm"
)

// CreateUser 创建用户
func CreateUser(user *model.User) error {
	return DB.Create(user).Error
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 用户不存在
	}
	return &user, nil
}

func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := DB.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 用户不存在
	}
	return &user, nil
}
