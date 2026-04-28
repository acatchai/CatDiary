package repository

import (
	"context"
	"time"

	"github.com/acatchai/catdiary/backend/internal/model"
)

type DraftDiarySnapshot struct {
	ID          uint64
	UserID      uint
	Title       string
	Content     string
	Mood        string
	Weather     string
	Location    string
	CreatedAtMs int64
	UpdatedAtMs int64
}

// UpsertDraftDiarySnapshot 插入或更新草稿日记快照
func UpsertDraftDiarySnapshot(ctx context.Context, s DraftDiarySnapshot) error {
	createdAt := time.UnixMilli(s.CreatedAtMs)
	updatedAt := time.UnixMilli(s.UpdatedAtMs)
	sql := `
INSERT INTO draft_diary
(id, user_id, title, content, mood, weather, location,
version, created_at, updated_at, deleted_at)
VALUES
(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL)
ON DUPLICATE KEY UPDATE
title = IF(VALUES(version) > version, VALUES(title),title),
content = IF(VALUES(version) > version, VALUES(content), content),
mood = IF(VALUES(version) > version, VALUES(mood), mood),
weather = IF(VALUES(version) > version, VALUES(weather), weather),
location = IF(VALUES(version) > version, VALUES(location), location),
updated_at = IF(VALUES(version) > version, VALUES(updated_at), updated_at),
version = GREATEST(version, VALUES(version)),
deleted_at = NULL
`
	return DB.WithContext(ctx).Exec(sql,
		s.ID, s.UserID, s.Title, s.Content, s.Mood, s.Weather, s.Location, createdAt, updatedAt).Error
}

// DeleteDraftDiaryMysql 删除草稿日记快照
func DeleteDraftDiaryMysql(ctx context.Context, userID uint, draftID uint64) error {
	return DB.WithContext(ctx).Where("id = ? AND user_id = ?", draftID, userID).Delete(&model.DraftDiary{}).Error
}
