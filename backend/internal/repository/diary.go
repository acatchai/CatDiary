package repository

import (
	"errors"

	"github.com/acatchai/catdiary/backend/internal/model"
	"gorm.io/gorm"
)

// CreateDiary 创建日记
func CreateDiary(diary *model.Diary) error {
	return DB.Create(diary).Error
}

// ListDiariesByUser 获取用户的所有日记
func ListDiariesByUser(userID uint, offset, limit int) ([]model.Diary, int64, error) {
	var total int64
	if err := DB.Model(&model.Diary{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Diary
	tx := DB.Where("user_id = ?", userID).Order("created_at DESC, id DESC").Offset(offset).Limit(limit).Find(&items)
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	return items, total, nil
}

// GetDiaryByIDAndUser 获取用户指定 ID 的日记
func GetDiaryByIDAndUser(userID, diaryID uint) (*model.Diary, error) {
	var diary model.Diary
	err := DB.Where("id = ? AND user_id = ?", diaryID, userID).First(&diary).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &diary, nil
}

// UpdateDiaryByIDAndUser 更新用户指定 ID 的日记
func UpdateDiaryByIDAndUser(userID, diaryID uint, fields map[string]any) (int64, error) {
	if len(fields) == 0 {
		return 0, nil
	}
	tx := DB.Model(&model.Diary{}).Where("id = ? AND user_id = ?", diaryID, userID).Updates(fields)
	return tx.RowsAffected, tx.Error
}

// DeleteDiaryByIDAndUser 删除用户指定 ID 的日记
func DeleteDiaryByIDAndUser(userID, diaryID uint) (int64, error) {
	tx := DB.Where("id = ? AND user_id = ?", diaryID, userID).Delete(&model.Diary{})
	return tx.RowsAffected, tx.Error
}
