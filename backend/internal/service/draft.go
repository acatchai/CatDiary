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

// ListDraftDiaries 获取草稿日记列表
func ListDraftDiaries(userID uint, page, pageSize int) ([]DraftDiary, int64, error) {
	items, total, err := repository.ListDraftDiariesRedis(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]DraftDiary, 0, len(items))
	for i := range items {
		d := items[i]
		out = append(out, *mapDraft(&d))
	}
	return out, total, nil
}

// PutDraftDiary 更新草稿日记
func PutDraftDiary(userID uint, draftID uint64, expectedVersion *uint64, title, content, mood, weather, location string) (*DraftDiary, error) {
	d, cur, err := repository.PutDraftDiaryRedis(
		userID,
		draftID,
		expectedVersion,
		strings.TrimSpace(title),
		content,
		strings.TrimSpace(mood),
		strings.TrimSpace(weather),
		strings.TrimSpace(location),
		defaultDraftTTL,
		defaultDebounce,
	)
	if err != nil {
		if err == repository.ErrDraftNotFoundRedis {
			return nil, ErrDraftNotFound
		}
		if err == repository.ErrDraftConflictRedis {
			return nil, &DraftConflictError{CurrentVersion: cur}
		}
		return nil, err
	}
	return mapDraft(d), nil
}

// PatchDraftDiary 更新草稿日记
func PatchDraftDiary(userID uint, draftID uint64, expectedVersion *uint64, title, content, mood, weather, location *string) (*DraftDiary, error) {
	fields := map[string]string{}
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
		return nil, ErrNoDraftUpdates
	}
	d, cur, err := repository.PatchDraftDiaryRedis(
		userID, draftID, expectedVersion,
		fields,
		defaultDraftTTL,
		defaultDebounce,
	)
	if err != nil {
		if err == repository.ErrDraftNotFoundRedis {
			return nil, ErrDraftNotFound
		}
		if err == repository.ErrDraftConflictRedis {
			return nil, &DraftConflictError{CurrentVersion: cur}
		}
		return nil, err
	}
	return mapDraft(d), nil
}

// DeleteDraftDiary 删除草稿日记
func DeleteDraftDiary(userID uint, draftID uint64) error {
	if err := repository.DeleteDraftDiaryRedis(userID, draftID, defaultDeleteTTL); err != nil {
		if err == repository.ErrDraftNotFoundRedis {
			return ErrDraftNotFound
		}
		return err
	}
	return nil
}

func FlushDraftDiary(userID uint, draftID uint64) error {
	if err := repository.FlushDraftDiaryRedis(userID, draftID); err != nil {
		if err == repository.ErrDraftNotFoundRedis {
			return ErrDraftNotFound
		}
		return err
	}
	return nil
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
