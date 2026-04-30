package service

import (
	"errors"
	"strings"
	"time"

	"github.com/acatchai/catdiary/backend/internal/model"
	"github.com/acatchai/catdiary/backend/internal/repository"
)

var (
	ErrDiaryNotFound  = errors.New("日记不存在")
	ErrNoDiaryUpdates = errors.New("没有可更新的字段")
)

func normalizeDiaryOccurredAt(d *model.Diary) {
	if d == nil {
		return
	}
	if d.OccurredAt == nil {
		d.OccurredAt = &d.CreatedAt
	}
}

// CreateDiary 创建日记
func CreateDiary(userID uint, occurredAt *time.Time, title, content, mood, weather, location string) (*model.Diary, error) {
	if occurredAt == nil {
		t := time.Now()
		occurredAt = &t
	}
	diary := &model.Diary{
		UserID:     userID,
		Title:    strings.TrimSpace(title),
		Content:  content,
		Mood:     strings.TrimSpace(mood),
		Weather:  strings.TrimSpace(weather),
		Location: strings.TrimSpace(location),
		OccurredAt: occurredAt,
	}
	if err := repository.CreateDiary(diary); err != nil {
		return nil, err
	}
	normalizeDiaryOccurredAt(diary)
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
	items, total, err := repository.ListDiariesByUser(userID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		if items[i].OccurredAt == nil {
			items[i].OccurredAt = &items[i].CreatedAt
		}
	}
	return items, total, nil
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
	normalizeDiaryOccurredAt(diary)
	return diary, nil
}

// PutDiary 更新日记
func PutDiary(userID, diaryID uint, occurredAt *time.Time, title, content, mood, weather, location string) (*model.Diary, error) {
	fields := map[string]any{
		"title":    strings.TrimSpace(title),
		"content":  content,
		"mood":     strings.TrimSpace(mood),
		"weather":  strings.TrimSpace(weather),
		"location": strings.TrimSpace(location),
	}
	if occurredAt != nil {
		fields["occurred_at"] = *occurredAt
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
func PatchDiary(userID, diaryID uint, occurredAt *time.Time, title, content, mood, weather, location *string) (*model.Diary, error) {
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
	if occurredAt != nil {
		fields["occurred_at"] = *occurredAt
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
