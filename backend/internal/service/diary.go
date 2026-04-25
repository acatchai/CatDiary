package service

import (
	"errors"
	"strings"

	"github.com/acatchai/catdiary/backend/internal/model"
	"github.com/acatchai/catdiary/backend/internal/repository"
)

var (
	ErrDiaryNotFound  = errors.New("日记不存在")
	ErrNoDiaryUpdates = errors.New("没有可更新的字段")
)

// CreateDiary 创建日记
func CreateDiary(userID uint, title, content, mood, weather, location string) (*model.Diary, error) {
	diary := &model.Diary{
		UserID:   userID,
		Title:    strings.TrimSpace(title),
		Content:  content,
		Mood:     strings.TrimSpace(mood),
		Weather:  strings.TrimSpace(weather),
		Location: strings.TrimSpace(location),
	}
	if err := repository.CreateDiary(diary); err != nil {
		return nil, err
	}
	return diary, nil
}

// ListDiaries 获取用户日记列表
func ListDiaries(userID uint, page, pageSize int) ([]model.Diary, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return repository.ListDiariesByUser(userID, offset, pageSize)
}

// GetDiary 获取日记详情
func GetDiary(userID, diaryID uint) (*model.Diary, error) {
	diary, err := repository.GetDiaryByIDAndUser(userID, diaryID)
	if err != nil {
		return nil, err
	}
	if diary == nil {
		return nil, ErrDiaryNotFound
	}
	return diary, nil
}

// PutDiary 更新日记
func PutDiary(userID, diaryID uint, title, content, mood, weather, location string) (*model.Diary, error) {
	fields := map[string]any{
		"title":    strings.TrimSpace(title),
		"content":  content,
		"mood":     strings.TrimSpace(mood),
		"weather":  strings.TrimSpace(weather),
		"location": strings.TrimSpace(location),
	}
	affected, err := repository.UpdateDiaryByIDAndUser(userID, diaryID, fields)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrDiaryNotFound
	}
	return GetDiary(userID, diaryID)
}

// PatchDiary 更新日记内容
func PatchDiary(userID, diaryID uint, title, content, mood, weather, location *string) (*model.Diary, error) {
	fields := make(map[string]any)
	if title != nil {
		fields["title"] = strings.TrimSpace(*title)
	}
	if content != nil {
		fields["content"] = *content
	}
	if mood != nil {
		fields["mood"] = strings.TrimSpace(*mood)
	}
	if weather != nil {
		fields["weather"] = strings.TrimSpace(*weather)
	}
	if location != nil {
		fields["location"] = strings.TrimSpace(*location)
	}

	if len(fields) == 0 {
		return nil, ErrNoDiaryUpdates
	}

	affected, err := repository.UpdateDiaryByIDAndUser(userID, diaryID, fields)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrDiaryNotFound
	}
	return GetDiary(userID, diaryID)
}

func DeleteDiary(userID, diaryID uint) error {
	affected, err := repository.DeleteDiaryByIDAndUser(userID, diaryID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrDiaryNotFound
	}
	return nil
}
