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
	local v = ARGC[6 + (i-1)*2 + 2]
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
