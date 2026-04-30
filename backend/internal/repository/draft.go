package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrDraftNotFoundRedis = errors.New("draft not found")
	ErrDraftConflictRedis = errors.New("draft conflict")
)

type DraftDiaryRedis struct {
	ID          uint64
	UserID      uint
	Title       string
	Content     string
	Mood        string
	Weather     string
	Location    string
	Version     uint64
	CreatedAtMs int64
	UpdatedAtMs int64
	Deleted     bool
}

const (
	draftDiaryNextIDKey    = "draft:diary:next_id"
	draftDiaryDirtyZSetKey = "draft:diary:dirty"
)

// DraftDiaryKey 获取草稿日记的key
func DraftDiaryKey(userID uint, draftID uint64) string {
	return fmt.Sprintf("draft:diary:%d:%d", userID, draftID)
}

// DraftDiaryIndexKey 获取草稿日记的索引key
func DraftDiaryIndexKey(userID uint) string {
	return fmt.Sprintf("draft:diary:index:%d", userID)
}

// DraftDiaryLockKey 获取草稿日记的锁key
func DraftDiaryLockKey(userID uint, draftID uint64) string {
	return fmt.Sprintf("draft:diary:lock:%d:%d", userID, draftID)
}

// ParseDraftDiaryKey 解析草稿日记的key
func ParseDraftDiaryKey(key string) (uint, uint64, bool) {
	parts := strings.Split(key, ":")
	if len(parts) != 4 {
		return 0, 0, false
	}
	if parts[0] != "draft" || parts[1] != "diary" {
		return 0, 0, false
	}
	uid64, err1 := strconv.ParseUint(parts[2], 10, 64)
	did, err2 := strconv.ParseUint(parts[3], 10, 64)
	if err1 != nil || err2 != nil || uid64 == 0 || did == 0 {
		return 0, 0, false
	}
	return uint(uid64), did, true
}

// DraftDirtyKey 获取草稿日记的脏key
func DraftDirtyKey() string {
	return draftDiaryDirtyZSetKey
}

func DraftNextID(ctx context.Context) (uint64, error) {
	if RDB == nil {
		return 0, errors.New("redis not initialized")
	}
	id, err := RDB.Incr(ctx, draftDiaryNextIDKey).Result()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

// draftCreateScript 创建草稿日记的脚本
var draftCreateScript = redis.NewScript(`
local draftKey = KEYS[1]
local indexKey = KEYS[2]
local dirtyKey = KEYS[3]

local userID = ARGV[1]
local draftID = ARGV[2]
local ttlMs = tonumber(ARGV[3])
local debounceMs = tonumber(ARGV[4])
local nowMs = tonumber(ARGV[5])

redis.call('HSET', draftKey,
	'id', draftID,
	'user_id', userID,
	'title', ARGV[6],
	'content', ARGV[7],
	'mood', ARGV[8],
	'weather', ARGV[9],
	'location', ARGV[10],
	'version', 1,
	'created_at', nowMs,
	'updated_at', nowMs,
	'deleted', 0
)

redis.call("PEXPIRE", draftKey, ttlMs)
redis.call('ZADD', indexKey, nowMs, draftID)
redis.call('ZADD', dirtyKey, nowMs+debounceMs, draftKey)

return {1, nowMs}
`)

// draftPutScript 更新草稿日记的脚本
var draftPutScript = redis.NewScript(`
local draftKey = KEYS[1]
local indexKey = KEYS[2]
local dirtyKey = KEYS[3]

local expectedStr = ARGV[1]
local ttlMs = tonumber(ARGV[2])
local debounceMs = tonumber(ARGV[3])
local nowMs = tonumber(ARGV[4])
local draftID = ARGV[5]

if redis.call('EXISTS', draftKey) == 0 then
	return {'NOT_FOUND'}
end

local deleted = redis.call('HGET', draftKey, 'deleted')
if deleted == '1' then
	return {'NOT_FOUND'}
end

local cur = tonumber(redis.call('HGET',draftKey,
'version') or '0')
if expectedStr ~= '' then
	local expected = tonumber(expectedStr)
	if expected ~= cur then
		return {'CONFLICT',tostring(cur)}
	end
end

redis.call('HSET', draftKey,
	'title', ARGV[6],
	'content', ARGV[7],
	'mood', ARGV[8],
	'weather', ARGV[9],
	'location', ARGV[10],
	'updated_at', nowMs,
	'deleted', 0
)

local newv = redis.call('HINCRBY', draftKey, 'version' , 1)
redis.call('PEXPIRE', draftKey, ttlMs)
redis.call('ZADD', indexKey, nowMs, draftID)
redis.call('ZADD', dirtyKey, nowMs + debounceMs, draftKey)

return {'OK', tostring(newv), tostring(nowMs)}
`)

// draftPatchScript 更新草稿日记的脚本
var draftPatchScript = redis.NewScript(`
local draftKey = KEYS[1]
local indexKey = KEYS[2]
local dirtyKey = KEYS[3]

local expectedStr = ARGV[1]
local ttlMs = tonumber(ARGV[2])
local debounceMs = tonumber(ARGV[3])
local nowMs = tonumber(ARGV[4])
local draftID = ARGV[5]
local n = tonumber(ARGV[6])

if redis.call('EXISTS', draftKey) == 0 then
	return {'NOT_FOUND'}
end

local deleted = redis.call('HGET', draftKey, 'deleted')
if deleted == '1' then
	return {'NOT_FOUND'}
end

local cur = tonumber(redis.call('HGET', draftKey, 'version') or '0')
if expectedStr ~= '' then
	local expected = tonumber(expectedStr)
	if expected ~= cur then
		return {'CONFLICT',tostring(cur)}
	end
end

if n <= 0 then
	return {'NO_UPDATES'}
end

for i=1,n do
	local k = ARGV[6 + (i-1)*2 + 1]
	local v = ARGV[6 + (i-1)*2 + 2]
	redis.call('HSET', draftKey, k, v)
end

redis.call('HSET', draftKey, 'updated_at', nowMs, 'deleted', 0)
local newv = redis.call('HINCRBY', draftKey, 'version', 1)

redis.call('PEXPIRE',draftKey, ttlMs)
redis.call('ZADD', indexKey, nowMs, draftID)
redis.call('ZADD',dirtyKey, nowMs + debounceMs, draftKey)

return {'OK', tostring(newv), tostring(nowMs)}
`)

// draftDeleteScript 删除草稿日记的脚本
var draftDeleteScript = redis.NewScript(`
local draftKey = KEYS[1]
local indexKey = KEYS[2]
local dirtyKey = KEYS[3]

local ttlDeleteMs = tonumber(ARGV[1])
local nowMs = tonumber(ARGV[2])
local draftID = ARGV[3]

if redis.call('EXISTS', draftKey) == 0 then
	return {'NOT_FOUND'}
end

local deleted = redis.call('HGET', draftKey, 'deleted')
if deleted == '1' then
	return {'NOT_FOUND'}
end

redis.call('HSET', draftKey, 'deleted', 1, 'updated_at', nowMs)
local newv = redis.call('HINCRBY', draftKey, 'version', 1)

redis.call('PEXPIRE', draftKey, ttlDeleteMs)
redis.call('ZREM', indexKey, draftID)
redis.call('ZADD', dirtyKey, nowMs, draftKey)

return {'OK', tostring(newv), tostring(nowMs)}
`)

// CreateDraftDiaryRedis 创建草稿日记的Redis对象
func CreateDraftDiaryRedis(userID uint, title, content, mood, weather, location string, ttl time.Duration, debounce time.Duration) (*DraftDiaryRedis, error) {
	if RDB == nil {
		return nil, errors.New("redis not initialized")
	}
	ctx := context.Background()

	id, err := DraftNextID(ctx)
	if err != nil {
		return nil, err
	}

	nowMs := time.Now().UnixMilli()
	draftKey := DraftDiaryKey(userID, id)

	_, err = draftCreateScript.Run(ctx, RDB, []string{draftKey, DraftDiaryIndexKey(userID), DraftDirtyKey()},
		strconv.FormatUint(uint64(userID), 10),
		strconv.FormatUint(id, 10),
		strconv.FormatInt(ttl.Milliseconds(), 10),
		strconv.FormatInt(debounce.Milliseconds(), 10),
		strconv.FormatInt(nowMs, 10),
		title, content, mood, weather, location,
	).Result()
	if err != nil {
		return nil, err
	}
	return GetDraftDiaryRedis(userID, id)
}

// GetDraftDiaryRedis 获取草稿日记的Redis对象
func GetDraftDiaryRedis(userID uint, draftID uint64) (*DraftDiaryRedis, error) {
	if RDB == nil {
		return nil, errors.New("redis not initialized")
	}
	ctx := context.Background()

	key := DraftDiaryKey(userID, draftID)
	m, err := RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		return nil, ErrDraftNotFoundRedis
	}

	if m["deleted"] == "1" {
		return nil, ErrDraftNotFoundRedis
	}
	return mapToDraftDiaryRedis(m)

}

// ListDraftDiariesRedis 获取草稿日记的Redis对象列表
func ListDraftDiariesRedis(userID uint, page, pageSize int) ([]DraftDiaryRedis, int64, error) {
	if RDB == nil {
		return nil, 0, errors.New("redis not initialized")
	}
	ctx := context.Background()

	indexKey := DraftDiaryIndexKey(userID)

	total, err := RDB.ZCard(ctx, indexKey).Result()
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []DraftDiaryRedis{}, 0, nil
	}

	start := int64((page - 1) * pageSize)
	stop := start + int64(pageSize) - 1
	ids, err := RDB.ZRevRange(ctx, indexKey, start, stop).Result()
	if err != nil {
		return nil, 0, err
	}
	pipe := RDB.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, 0, len(ids))

	for _, idStr := range ids {
		did, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil || did == 0 {
			continue
		}
		cmds = append(cmds, pipe.HGetAll(ctx, DraftDiaryKey(userID, did)))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, 0, err
	}

	items := make([]DraftDiaryRedis, 0, len(cmds))
	for _, cmd := range cmds {
		m, err := cmd.Result()
		if err != nil {
			continue
		}
		if len(m) == 0 || m["deleted"] == "1" {
			continue
		}
		d, err := mapToDraftDiaryRedis(m)
		if err != nil {
			continue
		}
		items = append(items, *d)
	}
	return items, total, nil
}

func PutDraftDiaryRedis(userID uint, draftID uint64, expectedVersion *uint64, title, content, mood, weather, location string, ttl time.Duration, debounce time.Duration) (*DraftDiaryRedis, uint64, error) {
	if RDB == nil {
		return nil, 0, errors.New("redis not initialized")
	}
	ctx := context.Background()
	expected := ""
	if expectedVersion != nil {
		expected = strconv.FormatUint(*expectedVersion, 10)
	}
	nowMs := time.Now().UnixMilli()
	key := DraftDiaryKey(userID, draftID)

	res, err := draftPutScript.Run(ctx, RDB, []string{key, DraftDiaryIndexKey(userID), DraftDirtyKey()},
		expected,
		strconv.FormatInt(ttl.Milliseconds(), 10),
		strconv.FormatInt(debounce.Milliseconds(), 10),
		strconv.FormatInt(nowMs, 10),
		strconv.FormatUint(draftID, 10),
		title, content, mood, weather, location,
	).Result()
	if err != nil {
		return nil, 0, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) == 0 {
		return nil, 0, errors.New("unexpected redis response")
	}

	switch fmt.Sprint(arr[0]) {
	case "NOT_FOUND":
		return nil, 0, ErrDraftNotFoundRedis
	case "CONFLICT":
		if len(arr) >= 2 {
			cur, _ := strconv.ParseUint(fmt.Sprint(arr[1]), 10, 64)
			return nil, cur, ErrDraftConflictRedis
		}
		return nil, 0, ErrDraftConflictRedis
	case "OK":
		var cur uint64
		if len(arr) >= 2 {
			cur, _ = strconv.ParseUint(fmt.Sprint(arr[1]), 10, 64)
		}
		d, err := GetDraftDiaryRedis(userID, draftID)
		if err != nil {
			return nil, cur, err
		}
		return d, cur, nil
	default:
		return nil, 0, errors.New("unexpected redis response")
	}
}

// PatchDraftDiaryRedis 更新草稿日记的Redis对象
func PatchDraftDiaryRedis(userID uint, draftID uint64, expectedVersion *uint64, fields map[string]string, ttl time.Duration, debounce time.Duration) (*DraftDiaryRedis, uint64, error) {
	if RDB == nil {
		return nil, 0, errors.New("redis not initialized")
	}
	ctx := context.Background()

	expected := ""
	if expectedVersion != nil {
		expected = strconv.FormatUint(*expectedVersion, 10)
	}

	nowMs := time.Now().UnixMilli()
	key := DraftDiaryKey(userID, draftID)

	args := make([]any, 0, 6+len(fields)*2)
	args = append(args,
		expected,
		strconv.FormatInt(ttl.Milliseconds(), 10),
		strconv.FormatInt(debounce.Milliseconds(), 10),
		strconv.FormatInt(nowMs, 10),
		strconv.FormatUint(draftID, 10),
		strconv.Itoa(len(fields)),
	)
	for k, v := range fields {
		args = append(args, k, v)
	}

	res, err := draftPatchScript.Run(ctx, RDB, []string{key, DraftDiaryIndexKey(userID), DraftDirtyKey()}, args...).Result()
	if err != nil {
		return nil, 0, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) == 0 {
		return nil, 0, errors.New("unexpected redis response")
	}

	switch fmt.Sprint(arr[0]) {
	case "NOT_FOUND":
		return nil, 0, ErrDraftNotFoundRedis
	case "NO_UPDATES":
		return nil, 0, errors.New("no updates")
	case "CONFLICT":
		if len(arr) >= 2 {
			cur, _ := strconv.ParseUint(fmt.Sprint(arr[1]), 10, 64)
			return nil, cur, ErrDraftConflictRedis
		}
		return nil, 0, ErrDraftConflictRedis
	case "OK":
		var cur uint64
		if len(arr) >= 2 {
			cur, _ = strconv.ParseUint(fmt.Sprint(arr[1]), 10, 64)
		}
		d, err := GetDraftDiaryRedis(userID, draftID)
		if err != nil {
			return nil, cur, err
		}
		return d, cur, nil
	default:
		return nil, 0, errors.New("unexpected redis response")
	}
}

// DeleteDraftDiaryRedis 删除草稿日记的Redis对象
func DeleteDraftDiaryRedis(userID uint, draftID uint64, deleteTTL time.Duration) error {
	if RDB == nil {
		return errors.New("redis not initialized")
	}
	ctx := context.Background()

	nowMs := time.Now().UnixMilli()
	key := DraftDiaryKey(userID, draftID)

	res, err := draftDeleteScript.Run(ctx, RDB, []string{key, DraftDiaryIndexKey(userID), DraftDirtyKey()},
		strconv.FormatInt(deleteTTL.Milliseconds(), 10),
		strconv.FormatInt(nowMs, 10),
		strconv.FormatUint(draftID, 10),
	).Result()
	if err != nil {
		return err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) == 0 {
		return errors.New("unexpected redis response")
	}
	if fmt.Sprint(arr[0]) == "NOT_FOUND" {
		return ErrDraftNotFoundRedis
	}
	return nil
}

func FlushDraftDiaryRedis(userID uint, draftID uint64) error {
	if RDB == nil {
		return errors.New("redis not initialized")
	}
	ctx := context.Background()

	key := DraftDiaryKey(userID, draftID)
	exists, err := RDB.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return ErrDraftNotFoundRedis
	}
	if v, err := RDB.HGet(ctx, key, "deleted").Result(); err == nil && v == "1" {
		return ErrDraftNotFoundRedis
	}

	nowMs := time.Now().UnixMilli()
	return RDB.ZAdd(ctx, DraftDirtyKey(), redis.Z{Score: float64(nowMs), Member: key}).Err()
}

// DraftDirtyDueKeys 获取待刷新的草稿日记的Redis对象
func DraftDirtyDueKeys(ctx context.Context, nowMs int64, limit int64) ([]string, error) {
	if RDB == nil {
		return nil, errors.New("redis not initialized")
	}
	return RDB.ZRangeByScore(ctx, DraftDirtyKey(), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatInt(nowMs, 10),
		Offset: 0,
		Count:  limit,
	}).Result()
}

// DraftDirtyRemove 删除待刷新的草稿日记的Redis对象
func DraftDirtyRemove(ctx context.Context, key string) error {
	if RDB == nil {
		return errors.New("redis not initialized")
	}
	return RDB.ZRem(ctx, DraftDirtyKey(), key).Err()
}

// DraftDirtyReschedule 重新调度待刷新的草稿日记的Redis对象
func DraftDirtyReschedule(ctx context.Context, key string, nextMs int64) error {
	if RDB == nil {
		return errors.New("redis not initialized")
	}
	return RDB.ZAdd(ctx, DraftDirtyKey(), redis.Z{Score: float64(nextMs), Member: key}).Err()
}

// TryDraftLock 尝试获取草稿日记的锁
func TryDraftLock(ctx context.Context, userID uint, draftID uint64, token string, ttl time.Duration) (bool, error) {
	if RDB == nil {
		return false, errors.New("redis not initialized")
	}
	return RDB.SetNX(ctx, DraftDiaryLockKey(userID, draftID), token, ttl).Result()
}

var unlockScript = redis.NewScript(`
if redis.call('GET', KEYS[1]) == ARGV[1] then
	return redis.call('DEL',KEYS[1])
end
return 0
`)

// UnlockDraft 解锁草稿日记
func UnlockDraft(ctx context.Context, userID uint, draftID uint64, token string) error {
	if RDB == nil {
		return errors.New("redis not initialized")
	}
	_, err := unlockScript.Run(ctx, RDB, []string{DraftDiaryLockKey(userID, draftID)}, token).Result()
	return err
}

func mapToDraftDiaryRedis(m map[string]string) (*DraftDiaryRedis, error) {
	id, _ := strconv.ParseUint(m["id"], 10, 64)
	uid64, _ := strconv.ParseUint(m["user_id"], 10, 64)
	version, _ := strconv.ParseUint(m["version"], 10, 64)
	createdAtMs, _ := strconv.ParseInt(m["created_at"], 10, 64)
	updatedAtMs, _ := strconv.ParseInt(m["updated_at"], 10, 64)

	return &DraftDiaryRedis{
		ID:          id,
		UserID:      uint(uid64),
		Title:       m["title"],
		Content:     m["content"],
		Mood:        m["mood"],
		Weather:     m["weather"],
		Location:    m["location"],
		Version:     version,
		CreatedAtMs: createdAtMs,
		UpdatedAtMs: updatedAtMs,
		Deleted:     m["deleted"] == "1",
	}, nil
}
