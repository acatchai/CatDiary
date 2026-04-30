package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/acatchai/catdiary/backend/internal/repository"
)

func StartDraftFlusher() {
	go runDraftFlusher()
}

func runDraftFlusher() {
	ctx := context.Background()

	batchSize := int64(100)
	lockTTL := 30 * time.Second
	tick := 200 * time.Millisecond

	for {
		time.Sleep(tick)

		nowMs := time.Now().UnixMilli()
		keys, err := repository.DraftDirtyDueKeys(ctx, nowMs, batchSize)
		if err != nil {
			continue
		}
		if len(keys) == 0 {
			continue
		}
		for _, key := range keys {
			userID, draftID, ok := repository.ParseDraftDiaryKey(key)
			if !ok {
				_ = repository.DraftDirtyRemove(ctx, key)
				continue
			}

			token := fmt.Sprintf("%d:%d", time.Now().UnixNano(), draftID)
			locked, err := repository.TryDraftLock(ctx, userID, draftID, token, lockTTL)
			if err != nil {
				_ = repository.DraftDirtyReschedule(ctx, key, time.Now().Add(2*time.Second).UnixMilli())
				continue
			}
			if !locked {
				continue
			}

			func() {
				defer func() {
					_ = repository.UnlockDraft(ctx, userID, draftID, token)
				}()
				d, err := repository.GetDraftDiaryRedis(userID, draftID)
				if err != nil {
					if err == repository.ErrDraftNotFoundRedis {
						_ = repository.DraftDirtyRemove(ctx, key)
						return
					}
					_ = repository.DraftDirtyReschedule(ctx, key, time.Now().Add(2*time.Second).UnixMilli())
					return
				}

				if d.Deleted {
					if err := repository.DeleteDraftDiaryMySQL(ctx, userID, draftID); err != nil {
						_ = repository.DraftDirtyReschedule(ctx, key, time.Now().Add(5*time.Second).UnixMilli())
						return
					}
					_ = repository.DraftDirtyRemove(ctx, key)
					return
				}

				snap := repository.DraftDiarySnapshot{
					ID:          d.ID,
					UserID:      userID,
					Title:       d.Title,
					Content:     d.Content,
					Mood:        d.Mood,
					Weather:     d.Weather,
					Location:    d.Location,
					OccurredAtMs: d.OccurredAtMs,
					Version:     d.Version,
					CreatedAtMs: d.CreatedAtMs,
					UpdatedAtMs: d.UpdatedAtMs,
				}
				if err := repository.UpsertDraftDiarySnapshot(ctx, snap); err != nil {
					_ = repository.DraftDirtyReschedule(ctx, key, time.Now().Add(5*time.Second).UnixMilli())
					return
				}

				_ = repository.DraftDirtyRemove(ctx, key)
			}()
		}
	}
}
