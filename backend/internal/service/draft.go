package service

import (
	"errors"
	"strings"
	"time"

	"github.com/acatchai/catdiary/backend/internal/repository"
)

var (
	ErrDraftNotFound  = errors.New("draft not found")
	ErrNoDraftUpdates = errors.New("no draft updates")
	defaultDraftTTL   = 30 * 24 * time.Hour
	defaultDebounce   = 2 * time.Second
	defaultDeleteTTL  = 1 * time.Hour
)

type DraftConflictError struct {
	CurrentVersion uint64
}

func (e *DraftConflictError) Error() string {
	return "draft version conflict"
}

type DraftDiary struct {
	ID        uint64 `json:"id"`
	UserID    uint   `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Mood      string `json:"mood"`
	Weather   string `json:"weather"`
	Location  string `json:"location"`
	Version   uint64 `json:"version"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// CreateDraftDiary 创建草稿日记
func CreateDraftDiary(userID uint, title, content, mood, weather, location string) (*DraftDiary, error) {
	d, err := repository.CreateDraftDiaryRedis(userID,
		strings.TrimSpace(title),
		content,
		strings.TrimSpace(mood),
		strings.TrimSpace(weather),
		strings.TrimSpace(location),
		defaultDraftTTL,
		defaultDebounce,
	)
	if err != nil {
		return nil, err
	}
	return mapDraft(d), nil
}

// GetDraftDiary 获取草稿日记
func GetDraftDiary(userID uint, draftID uint64) (*DraftDiary, error) {
	d, err := repository.GetDraftDiaryRedis(userID, draftID)
	if err != nil {
		if err == repository.ErrDraftNotFoundRedis {
			return nil, ErrDraftNotFound
		}
		return nil, err
	}
	return mapDraft(d), nil
}

// mapDraft 映射草稿日记
func mapDraft(d *repository.DraftDiaryRedis) *DraftDiary {
	return &DraftDiary{
		ID:        d.ID,
		UserID:    d.UserID,
		Title:     d.Title,
		Content:   d.Content,
		Mood:      d.Mood,
		Weather:   d.Weather,
		Location:  d.Location,
		Version:   d.Version,
		CreatedAt: d.CreatedAtMs,
		UpdatedAt: d.UpdatedAtMs,
	}
}
