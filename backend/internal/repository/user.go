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
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := DB.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 用户不存在
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserByID 更新用户信息
func UpdateUserByID(id uint, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	return DB.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func UpdateUserPasswordHash(id uint, passwordHash string) error {
	return DB.Model(&model.User{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
}

