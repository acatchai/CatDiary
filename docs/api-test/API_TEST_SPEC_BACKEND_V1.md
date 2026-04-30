<!-- markdownlint-disable MD013 MD024 -->

# API 测试规范（Backend / API v1）

- 文档版本：v1
- 代码同步基线：[router.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/router/router.go)
- 覆盖范围：当前路由文件中已注册的全部接口（/healthz、/readyz、/api/v1/**）
- 约定：除特别说明外，响应体均为 JSON；认证通过 `Authorization: Bearer <token>` 进行

## 通用约定

### 基础信息

| 项 | 值 |
| --- | --- |
| Base Path | `/api/v1` |
| Content-Type（JSON 请求） | `application/json; charset=utf-8` |
| Accept（建议） | `application/json` |
| 认证（需要登录的接口） | `Authorization: Bearer <JWT>` |

### JWT 认证规则（与源码一致）

| 项 | 规则 |
| --- | --- |
| Header | `Authorization: Bearer <token>` |
| 缺少 Header | `401`，`{"error":"缺少Authorization头部"}` |
| Header 格式错误 | `401`，`{"error":"Authorization 格式错误，应为 Bearer <token>"}` |
| Token 无效/过期 | `401`，`{"error":"Token 无效或过期"}` |
| Token 内容 | JWT（HS256），claims：`userID`、`exp`（7 天） |

### 通用错误响应（约定）

多数接口错误响应形态为：

```json
{
  "error": "错误信息"
}
```

> 注意：部分接口还会返回额外字段（例如草稿版本冲突时返回 `current_version`）。

#### JSON Schema（ErrorResponse）

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://catdiary.local/schema/ErrorResponse.json",
  "type": "object",
  "additionalProperties": true,
  "required": ["error"],
  "properties": {
    "error": {
      "type": "string",
      "minLength": 1,
      "examples": ["未登录", "id 参数不合法", "服务器内部错误"]
    }
  }
}
```

### 通用字段图例（Legend）

| 图例 | 含义 |
| --- | --- |
| required | 必填 |
| optional | 选填 |
| enum | 枚举值集合 |
| min/max | 长度或数值边界 |
| format | JSON Schema format（如 `email`、`uri`、`date-time`） |

### 覆盖与验收要求（全局）

- 必填缺失：对所有 `required` 字段覆盖缺失场景（期望 `400` 或符合业务返回）
- 类型错误：字符串/数字/对象类型错误（期望 `400` 或符合业务返回）
- 长度超限：对 `min/max` 覆盖边界与越界（期望 `400` 或符合业务返回）
- 枚举值外：若存在枚举（本版本多为自由字符串），需覆盖枚举外输入（期望 `400` 或符合业务返回）
- 权限不足：所有需要登录接口覆盖未带/错误 token（期望 `401`）
- 并发冲突：草稿更新接口覆盖版本冲突（期望 `409`）
- 幂等性：PUT/DELETE 等覆盖重复提交（期望符合当前实现：重复 DELETE 多数会 `404`）
- 分页边界：`page`、`page_size` 覆盖 0、负数、超大、非数字（本版本实现：会回退默认并 clamp）
- 性能基准：在稳定环境下执行压测或基准测试，建议以 P95 作为断言
  - 读接口：P95 RT ≤ 300ms
  - 写接口：P95 RT ≤ 500ms
- 文档静态检查：markdownlint（建议 `markdownlint-cli2`）零错误
- Schema 校验：对文档内 JSON Schema 进行 JSON 语法校验并用 Schema 工具加载（例如 AJV）

## Components（复用 Schema）

### Schema：User（model.User）

源码依据：[user.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/model/user.go)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://catdiary.local/schema/User.json",
  "type": "object",
  "additionalProperties": true,
  "required": ["id", "username", "email", "avatar", "created_at", "updated_at"],
  "properties": {
    "id": { "type": "integer", "minimum": 1, "examples": [1] },
    "username": { "type": "string", "minLength": 1, "examples": ["alice"] },
    "email": { "type": "string", "examples": ["alice@example.com"] },
    "avatar": { "type": "string", "examples": ["https://example.com/a.png"] },
    "created_at": { "type": "string", "format": "date-time", "examples": ["2026-01-01T00:00:00Z"] },
    "updated_at": { "type": "string", "format": "date-time", "examples": ["2026-01-01T00:00:00Z"] }
  }
}
```

### Schema：Diary（model.Diary）

源码依据：[diary.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/model/diary.go)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://catdiary.local/schema/Diary.json",
  "type": "object",
  "additionalProperties": true,
  "required": ["id", "user_id", "title", "content", "mood", "weather", "location", "created_at", "updated_at"],
  "properties": {
    "id": { "type": "integer", "minimum": 1, "examples": [101] },
    "user_id": { "type": "integer", "minimum": 1, "examples": [1] },
    "title": { "type": "string", "minLength": 1, "maxLength": 100, "examples": ["今天的猫"] },
    "content": { "type": "string", "minLength": 1, "examples": ["# 标题\\n内容"] },
    "mood": { "type": "string", "maxLength": 20, "examples": ["happy"] },
    "weather": { "type": "string", "maxLength": 20, "examples": ["sunny"] },
    "location": { "type": "string", "maxLength": 100, "examples": ["Shanghai"] },
    "created_at": { "type": "string", "format": "date-time", "examples": ["2026-01-01T00:00:00Z"] },
    "updated_at": { "type": "string", "format": "date-time", "examples": ["2026-01-01T00:00:00Z"] }
  }
}
```

### Schema：DraftDiary（service.DraftDiary）

源码依据：[draft.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/service/draft.go#L27-L38)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://catdiary.local/schema/DraftDiary.json",
  "type": "object",
  "additionalProperties": true,
  "required": ["id", "user_id", "title", "content", "mood", "weather", "location", "version", "created_at", "updated_at"],
  "properties": {
    "id": { "type": "integer", "minimum": 1, "examples": [9007199254740991] },
    "user_id": { "type": "integer", "minimum": 1, "examples": [1] },
    "title": { "type": "string", "minLength": 1, "maxLength": 100, "examples": ["草稿标题"] },
    "content": { "type": "string", "minLength": 1, "examples": ["草稿内容"] },
    "mood": { "type": "string", "maxLength": 20, "examples": ["calm"] },
    "weather": { "type": "string", "maxLength": 20, "examples": ["cloudy"] },
    "location": { "type": "string", "maxLength": 100, "examples": ["Beijing"] },
    "version": { "type": "integer", "minimum": 1, "examples": [1] },
    "created_at": { "type": "integer", "minimum": 0, "description": "Unix 毫秒时间戳", "examples": [1767225600000] },
    "updated_at": { "type": "integer", "minimum": 0, "description": "Unix 毫秒时间戳", "examples": [1767225600000] }
  }
}
```

### Schema：DraftConflictResponse

源码依据：[draft.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/draft.go#L174-L180)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://catdiary.local/schema/DraftConflictResponse.json",
  "type": "object",
  "additionalProperties": true,
  "required": ["error", "current_version"],
  "properties": {
    "error": { "type": "string", "examples": ["草稿已被更新"] },
    "current_version": { "type": "integer", "minimum": 1, "examples": [3] }
  }
}
```

---

## 接口：GET /healthz

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/healthz` |
| Method | `GET` |
| 功能 | 健康检查（返回固定字符串） |

### 2) 请求定义

| 类别 | 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- | --- |
| Header | - | - | - | - | 无 |
| Path | - | - | - | - | 无 |
| Query | - | - | - | - | 无 |
| Body | - | - | - | - | 无 |

### 3) 响应定义

#### 成功响应

| 状态码 | 响应头 | 响应体 |
| --- | --- | --- |
| 200 | `Content-Type: text/plain; charset=utf-8` | `ok` |

#### 错误响应

无业务错误分支（按当前实现）。

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| HZ-001 | 服务启动 | 无 | 200，返回 `ok` | body 等于 `ok`，RT P95 ≤ 300ms |

### 5) 覆盖要求补充

- 性能：稳定环境连续 100 次请求，P95 RT ≤ 300ms

---

## 接口：GET /readyz

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/readyz` |
| Method | `GET` |
| 功能 | 就绪检查（返回固定字符串） |

### 2) 请求定义

| 类别 | 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- | --- |
| Header | - | - | - | - | 无 |
| Path | - | - | - | - | 无 |
| Query | - | - | - | - | 无 |
| Body | - | - | - | - | 无 |

### 3) 响应定义

#### 成功响应

| 状态码 | 响应头 | 响应体 |
| --- | --- | --- |
| 200 | `Content-Type: text/plain; charset=utf-8` | `ok` |

#### 错误响应

无业务错误分支（按当前实现）。

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| RZ-001 | 服务启动 | 无 | 200，返回 `ok` | body 等于 `ok`，RT P95 ≤ 300ms |

### 5) 覆盖要求补充

- 性能：稳定环境连续 100 次请求，P95 RT ≤ 300ms

---

## 接口：POST /api/v1/auth/register

源码依据：[auth.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/auth.go#L12-L49)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/auth/register` |
| Method | `POST` |
| 功能 | 注册用户 |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Content-Type | 是 | `application/json` | JSON 请求体 |

#### Body（RegisterReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| username | string | 是 | min=3, max=50 | 用户名唯一（冲突返回 409） |
| password | string | 是 | min=6, max=50 | 明文提交（服务端 bcrypt） |
| email | string | 否 | format=email | 允许缺省或空；不保证唯一性校验发生在本接口 |

#### JSON Schema（RegisterReq）

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["username", "password"],
  "properties": {
    "username": { "type": "string", "minLength": 3, "maxLength": 50 },
    "password": { "type": "string", "minLength": 6, "maxLength": 50 },
    "email": { "type": "string", "format": "email" }
  }
}
```

### 3) 响应定义

#### 成功响应（200）

响应示例：

```json
{
  "message": "注册成功",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com"
  }
}
```

JSON Schema（AuthRegisterResponse）：

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": true,
  "required": ["message", "data"],
  "properties": {
    "message": { "type": "string", "examples": ["注册成功"] },
    "data": {
      "type": "object",
      "additionalProperties": true,
      "required": ["id", "username", "email"],
      "properties": {
        "id": { "type": "integer", "minimum": 1, "examples": [1] },
        "username": { "type": "string", "examples": ["alice"] },
        "email": { "type": "string", "examples": ["alice@example.com"] }
      }
    }
  }
}
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败 | ErrorResponse |
| 409 | 用户名已存在 | `{"error":"用户名已存在"}` |
| 500 | 注册失败 | `{"error":"注册失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| REG-001 | 无同名用户 | username=alice, password=123456, email=`alice@example.com` | 200 | message=注册成功；data.id>0 |
| REG-002 | 无同名用户 | username 长度=3 | 200 | 允许最小边界 |
| REG-003 | 无同名用户 | username 长度=50 | 200 | 允许最大边界 |
| REG-004 | 无同名用户 | password 长度=6 | 200 | 允许最小边界 |
| REG-005 | 无同名用户 | password 长度=50 | 200 | 允许最大边界 |
| REG-006 | 无同名用户 | email 缺省 | 200 | data.email 可能为空字符串或空值（按实际返回断言） |
| REG-101 | - | 缺少 username | 400 | error 存在；提示包含 required |
| REG-102 | - | username 类型=number | 400 | error 存在 |
| REG-103 | - | username 长度=2 | 400 | error 存在；min=3 |
| REG-104 | - | username 长度=51 | 400 | error 存在；max=50 |
| REG-105 | - | 缺少 password | 400 | error 存在 |
| REG-106 | - | password 长度=5 | 400 | error 存在；min=6 |
| REG-107 | - | password 长度=51 | 400 | error 存在；max=50 |
| REG-108 | - | email=not-an-email | 400 | error 存在；email 格式错误 |
| REG-201 | 已存在用户名 alice | username=alice | 409 | error=用户名已存在 |
| REG-301 | 压测环境 | 合法请求并发 20 | 200/409 | 不出现 5xx；RT P95 ≤ 500ms |

### 5) 覆盖要求补充

- 幂等性：注册非幂等；重复相同 username 期望 `409`

---

## 接口：POST /api/v1/auth/login

源码依据：[auth.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/auth.go#L50-L89)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/auth/login` |
| Method | `POST` |
| 功能 | 用户登录（返回 token） |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Content-Type | 是 | `application/json` | JSON 请求体 |

#### Body（LoginReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| username | string | 是 | 非空 | 用户名必须存在 |
| password | string | 是 | 非空 | 密码必须匹配 |

#### JSON Schema（LoginReq）

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["username", "password"],
  "properties": {
    "username": { "type": "string", "minLength": 1 },
    "password": { "type": "string", "minLength": 1 }
  }
}
```

### 3) 响应定义

#### 成功响应（200）

响应示例：

```json
{
  "message": "登录成功",
  "token": "<jwt>",
  "user": { "id": 1, "username": "alice" }
}
```

JSON Schema（AuthLoginResponse）：

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": true,
  "required": ["message", "token", "user"],
  "properties": {
    "message": { "type": "string", "examples": ["登录成功"] },
    "token": { "type": "string", "minLength": 1, "examples": ["eyJhbGciOiJIUzI1NiIs..."] },
    "user": {
      "type": "object",
      "additionalProperties": true,
      "required": ["id", "username"],
      "properties": {
        "id": { "type": "integer", "minimum": 1, "examples": [1] },
        "username": { "type": "string", "examples": ["alice"] }
      }
    }
  }
}
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败 | ErrorResponse |
| 401 | 用户名或密码错误 | `{"error":"用户名或密码错误"}` |
| 500 | 登录失败 | `{"error":"登录失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| LOGIN-001 | 用户存在且密码正确 | username=alice, password=123456 | 200 | token 非空；user.id>0 |
| LOGIN-101 | - | 缺少 username | 400 | error 存在 |
| LOGIN-102 | - | username="" | 400 | error 存在 |
| LOGIN-103 | - | 缺少 password | 400 | error 存在 |
| LOGIN-104 | - | password="" | 400 | error 存在 |
| LOGIN-201 | 用户不存在 | username=not-exist, password=xx | 401 | error=用户名或密码错误 |
| LOGIN-202 | 用户存在 | password 错误 | 401 | error=用户名或密码错误 |
| LOGIN-301 | 压测环境 | 合法请求并发 20 | 200 | token 不为空；RT P95 ≤ 500ms |

### 5) 覆盖要求补充

- 安全：确认响应不包含 `password_hash` 等敏感字段

---

## 接口：POST /api/v1/auth/logout

源码依据：[auth.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/auth.go#L91-L99)、
[auth.go middleware](file:///d:/WorkSpace/Go/CatDiary/backend/internal/middleware/auth.go)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/auth/logout` |
| Method | `POST` |
| 功能 | 用户退出登录（当前实现：JWT 无状态，直接返回成功） |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |

### 3) 响应定义

#### 成功响应（200）

```json
{ "message": "退出成功" }
```

JSON Schema（MessageResponse）：

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": true,
  "required": ["message"],
  "properties": {
    "message": { "type": "string", "minLength": 1, "examples": ["退出成功"] }
  }
}
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 401 | 未登录/Token 无效 | ErrorResponse |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| LOGOUT-001 | 已登录 | Authorization 合法 | 200 | message=退出成功 |
| LOGOUT-101 | - | 缺少 Authorization | 401 | error=缺少Authorization头部 |
| LOGOUT-102 | - | Authorization=`Token xxx` | 401 | error 包含 Bearer 格式 |
| LOGOUT-103 | - | Authorization=`Bearer invalid` | 401 | error=Token 无效或过期 |
| LOGOUT-301 | 压测环境 | 合法请求并发 50 | 200 | RT P95 ≤ 300ms |

### 5) 覆盖要求补充

- 幂等性：重复调用应始终成功（只要 token 仍有效）

---

## 接口：GET /api/v1/auth/me

源码依据：[auth.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/auth.go#L101-L119)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/auth/me` |
| Method | `GET` |
| 功能 | 获取当前登录用户信息 |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |

### 3) 响应定义

#### 成功响应（200）

```json
{ "data": { "id": 1, "username": "alice", "email": "alice@example.com", "avatar": "", "created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-01T00:00:00Z" } }
```

JSON Schema：

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": true,
  "required": ["data"],
  "properties": {
    "data": { "$ref": "https://catdiary.local/schema/User.json" }
  }
}
```

字段结构图：

```text
data (User)
├─ id (int)
├─ username (string)
├─ email (string)
├─ avatar (string)
├─ created_at (date-time)
└─ updated_at (date-time)
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 500 | 获取失败 | `{"error":"获取用户信息失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| MEA-001 | 已登录 | Authorization 合法 | 200 | data.id=token userID；不包含 password_hash |
| MEA-101 | - | 缺少 Authorization | 401 | error=缺少Authorization头部 |
| MEA-102 | - | Authorization 格式错误 | 401 | error 包含 Bearer 格式 |
| MEA-103 | - | token 无效 | 401 | error=Token 无效或过期 |
| MEA-301 | 压测环境 | 合法请求并发 50 | 200 | RT P95 ≤ 300ms |

---

## 接口：GET /api/v1/users/me

源码依据：[user.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/user.go#L36-L78)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/users/me` |
| Method | `GET` |
| 功能 | 获取当前用户信息（与 /auth/me 类似，但包含 404 分支） |

### 2) 请求定义

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |

### 3) 响应定义

#### 成功响应（200）

同 `/api/v1/auth/me`。

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 用户不存在 | `{"error":"用户不存在"}` |
| 500 | 获取失败 | `{"error":"获取用户信息失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| UME-001 | 已登录且用户存在 | Authorization 合法 | 200 | data.id=token userID |
| UME-101 | - | 缺少 Authorization | 401 | error=缺少Authorization头部 |
| UME-201 | token 对应用户被删除 | Authorization 合法 | 404 | error=用户不存在 |
| UME-301 | 压测环境 | 合法请求并发 50 | 200 | RT P95 ≤ 300ms |

---

## 接口：PATCH /api/v1/users/me

源码依据：[user.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/user.go#L80-L135)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/users/me` |
| Method | `PATCH` |
| 功能 | 更新当前用户资料（部分更新） |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |
| Content-Type | 是 | `application/json` | JSON 请求体 |

#### Body（UserPatchMeReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| username | string | 否 | min=3, max=50 | trim 后不可为空；若变更则需全局唯一 |
| email | string | 否 | format=email | trim 后不可为空；若变更则需全局唯一 |
| avatar | string | 否 | format=url | trim |

#### JSON Schema（UserPatchMeReq）

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "username": { "type": "string", "minLength": 3, "maxLength": 50 },
    "email": { "type": "string", "format": "email" },
    "avatar": { "type": "string", "format": "uri" }
  }
}
```

### 3) 响应定义

#### 成功响应（200）

```json
{
  "message": "更新成功",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "avatar": "https://example.com/a.png",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-02T00:00:00Z"
  }
}
```

JSON Schema：

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": true,
  "required": ["message", "data"],
  "properties": {
    "message": { "type": "string", "examples": ["更新成功"] },
    "data": { "$ref": "https://catdiary.local/schema/User.json" }
  }
}
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败/业务校验失败 | ErrorResponse |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 用户不存在 | `{"error":"用户不存在"}` |
| 409 | 用户名或邮箱已存在 | `{"error":"用户名已存在"}` / `{"error":"邮箱已存在"}` |
| 500 | 更新失败 | `{"error":"更新用户信息失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| UPM-001 | 已登录 | {"username":"bob"} | 200 | message=更新成功；data.username=bob |
| UPM-002 | 已登录 | `{"email":"bob@example.com"}` | 200 | data.email 更新 |
| UPM-003 | 已登录 | `{"avatar":"https://example.com/a.png"}` | 200 | data.avatar 更新 |
| UPM-004 | 已登录 | {} | 200 | 允许空更新；返回仍为 User |
| UPM-011 | 已登录 | username 长度=3 | 200 | 边界通过 |
| UPM-012 | 已登录 | username 长度=50 | 200 | 边界通过 |
| UPM-101 | 已登录 | {"username":"ab"} | 400 | error 存在；min=3 |
| UPM-102 | 已登录 | {"username":""} | 400 | error 存在（validator 或业务 trim 为空） |
| UPM-103 | 已登录 | {"email":"not-email"} | 400 | error 存在 |
| UPM-104 | 已登录 | {"avatar":"not-url"} | 400 | error 存在 |
| UPM-201 | 已登录 | {"username":"existing"} | 409 | error=用户名已存在 |
| UPM-202 | 已登录 | `{"email":"exists@example.com"}` | 409 | error=邮箱已存在 |
| UPM-301 | 压测环境 | 并发 20 修改 avatar | 200/409 | 不出现 5xx；RT P95 ≤ 500ms |

### 5) 覆盖要求补充

- 并发冲突：本接口无版本控制，冲突以“最后写入”为准；需覆盖并发写入一致性与最终态断言

---

## 接口：PATCH /api/v1/users/me/password

源码依据：[user.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/user.go#L137-L148)、
[user.go service](file:///d:/WorkSpace/Go/CatDiary/backend/internal/service/user.go#L72-L89)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/users/me/password` |
| Method | `PATCH` |
| 功能 | 修改当前用户密码 |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |
| Content-Type | 是 | `application/json` | JSON 请求体 |

#### Body（UserPatchPasswordReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| old_password | string | 是 | 非空 | 必须匹配当前密码 |
| new_password | string | 是 | min=6, max=50 | bcrypt 存储；不返回 token |

#### JSON Schema（UserPatchPasswordReq）

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["old_password", "new_password"],
  "properties": {
    "old_password": { "type": "string", "minLength": 1 },
    "new_password": { "type": "string", "minLength": 6, "maxLength": 50 }
  }
}
```

### 3) 响应定义

#### 成功响应（200）

```json
{ "message": "密码修改成功" }
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败 | ErrorResponse |
| 401 | 未登录/旧密码错误 | `{"error":"缺少Authorization头部"}` / `{"error":"旧密码错误"}` |
| 404 | 用户不存在 | `{"error":"用户不存在"}` |
| 500 | 修改失败 | `{"error":"修改密码失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| PWD-001 | 已登录且旧密码正确 | old_password=123456, new_password=abcdef | 200 | message=密码修改成功 |
| PWD-002 | 已登录 | new_password 长度=6 | 200 | 边界通过 |
| PWD-003 | 已登录 | new_password 长度=50 | 200 | 边界通过 |
| PWD-101 | - | 缺少 old_password | 400 | error 存在 |
| PWD-102 | - | 缺少 new_password | 400 | error 存在 |
| PWD-103 | 已登录 | new_password 长度=5 | 400 | error 存在 |
| PWD-104 | 已登录 | new_password 长度=51 | 400 | error 存在 |
| PWD-201 | 已登录 | old_password 错误 | 401 | error=旧密码错误 |
| PWD-301 | 压测环境 | 并发 10 修改密码 | 200/401 | 不出现 5xx；RT P95 ≤ 500ms |

### 5) 覆盖要求补充

- 安全：修改成功后，旧密码登录应失败；新密码登录应成功

---

## 接口：POST /api/v1/diaries

源码依据：[diary.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/diary.go#L91-L117)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/diaries` |
| Method | `POST` |
| 功能 | 创建日记 |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |
| Content-Type | 是 | `application/json` | JSON 请求体 |

#### Body（DiaryCreateReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| title | string | 是 | min=1, max=100 | 标题 trim 后入库 |
| content | string | 是 | 非空 | Markdown 正文 |
| mood | string | 否 | max=20 | trim |
| weather | string | 否 | max=20 | trim |
| location | string | 否 | max=100 | trim |

#### JSON Schema（DiaryCreateReq）

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["title", "content"],
  "properties": {
    "title": { "type": "string", "minLength": 1, "maxLength": 100 },
    "content": { "type": "string", "minLength": 1 },
    "mood": { "type": "string", "maxLength": 20 },
    "weather": { "type": "string", "maxLength": 20 },
    "location": { "type": "string", "maxLength": 100 }
  }
}
```

### 3) 响应定义

#### 成功响应（201）

```json
{ "data": { "id": 101, "user_id": 1, "title": "今天的猫", "content": "# 标题", "mood": "happy", "weather": "sunny", "location": "Shanghai", "created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-01T00:00:00Z" } }
```

JSON Schema：

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": true,
  "required": ["data"],
  "properties": {
    "data": { "$ref": "https://catdiary.local/schema/Diary.json" }
  }
}
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败 | ErrorResponse |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 500 | 服务器错误 | `{"error":"服务器内部错误"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DC-001 | 已登录 | title="t", content="c" | 201 | data.id>0；data.user_id=token userID |
| DC-002 | 已登录 | title 长度=1 | 201 | 边界通过 |
| DC-003 | 已登录 | title 长度=100 | 201 | 边界通过 |
| DC-004 | 已登录 | mood 长度=20 | 201 | 边界通过 |
| DC-005 | 已登录 | location 长度=100 | 201 | 边界通过 |
| DC-101 | 已登录 | 缺少 title | 400 | error 存在 |
| DC-102 | 已登录 | title="" | 400 | error 存在 |
| DC-103 | 已登录 | title 长度=101 | 400 | error 存在 |
| DC-104 | 已登录 | 缺少 content | 400 | error 存在 |
| DC-105 | 已登录 | content 类型=object | 400 | error 存在 |
| DC-106 | 已登录 | mood 长度=21 | 400 | error 存在 |
| DC-107 | 已登录 | weather 长度=21 | 400 | error 存在 |
| DC-108 | 已登录 | location 长度=101 | 400 | error 存在 |
| DC-201 | 未登录 | 无 Authorization | 401 | error=缺少Authorization头部 |
| DC-301 | 压测环境 | 并发 20 创建 | 201 | 不出现 5xx；RT P95 ≤ 500ms |

---

## 接口：GET /api/v1/diaries

源码依据：[diary.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/diary.go#L119-L143)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/diaries` |
| Method | `GET` |
| 功能 | 获取日记列表（分页） |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |

#### Query

| 参数 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| page | int | 否 | clamp 到 [1, 1000000] | 非数字时回退默认 1 |
| page_size | int | 否 | clamp 到 [1, 100] | 非数字时回退默认 20 |

> 注意：本接口对非法 query 不返回 400，而是回退默认值并 clamp。

### 3) 响应定义

#### 成功响应（200）

```json
{
  "items": [],
  "page": 1,
  "page_size": 20,
  "total": 0
}
```

JSON Schema：

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": true,
  "required": ["items", "page", "page_size", "total"],
  "properties": {
    "items": { "type": "array", "items": { "$ref": "https://catdiary.local/schema/Diary.json" } },
    "page": { "type": "integer", "minimum": 1 },
    "page_size": { "type": "integer", "minimum": 1, "maximum": 100 },
    "total": { "type": "integer", "minimum": 0 }
  }
}
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 500 | 服务器错误 | `{"error":"服务器内部错误"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DL-001 | 已登录 | page=1&page_size=20 | 200 | items 为数组；page/page_size 回显 |
| DL-002 | 已登录且有 21 条 | page=2&page_size=20 | 200 | items 数量=1 |
| DL-003 | 已登录 | page=0 | 200 | page 被 clamp 为 1 |
| DL-004 | 已登录 | page=-1 | 200 | page 被 clamp 为 1（或回退默认） |
| DL-005 | 已登录 | page=abc | 200 | page 回退默认 1 |
| DL-006 | 已登录 | page_size=0 | 200 | page_size 被 clamp 为 1 |
| DL-007 | 已登录 | page_size=101 | 200 | page_size 被 clamp 为 100 |
| DL-008 | 已登录 | page_size=abc | 200 | page_size 回退默认 20 |
| DL-201 | 未登录 | 无 Authorization | 401 | error=缺少Authorization头部 |
| DL-301 | 压测环境 | 并发 50 列表 | 200 | RT P95 ≤ 300ms |

### 5) 覆盖要求补充

- 分页边界：覆盖 total=0、total>0、最后一页、超大页码（返回空 items）

---

## 接口：GET /api/v1/diaries/:id

源码依据：[diary.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/diary.go#L145-L171)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/diaries/:id` |
| Method | `GET` |
| 功能 | 获取日记详情 |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |

#### Path

| 参数 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| id | uint | 是 | 正整数（>0） | 非法返回 400 |

### 3) 响应定义

#### 成功响应（200）

```json
{ "data": { "id": 101, "user_id": 1, "title": "x", "content": "y", "mood": "", "weather": "", "location": "", "created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-01T00:00:00Z" } }
```

JSON Schema：`{data: Diary}`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": true,
  "required": ["data"],
  "properties": {
    "data": { "$ref": "https://catdiary.local/schema/Diary.json" }
  }
}
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | id 参数不合法 | `{"error":"id 参数不合法"}` |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 日记不存在 | `{"error":"日记不存在"}` |
| 500 | 服务器错误 | `{"error":"服务器内部错误"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DG-001 | 已登录且该 id 存在 | id=101 | 200 | data.id=101；data.user_id=token userID |
| DG-101 | 已登录 | id=0 | 400 | error=id 参数不合法 |
| DG-102 | 已登录 | id=abc | 400 | error=id 参数不合法 |
| DG-201 | 已登录 | id 不存在 | 404 | error=日记不存在 |
| DG-301 | 压测环境 | 并发 50 获取 | 200/404 | RT P95 ≤ 300ms |

---

## 接口：PUT /api/v1/diaries/:id

源码依据：[diary.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/diary.go#L173-L207)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/diaries/:id` |
| Method | `PUT` |
| 功能 | 全量更新日记 |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |
| Content-Type | 是 | `application/json` | JSON 请求体 |

#### Path

| 参数 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| id | uint | 是 | 正整数（>0） | 非法返回 400 |

#### Body（DiaryPutReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| title | string | 是 | min=1, max=100 | 标题 trim 后入库 |
| content | string | 是 | 非空 | Markdown 正文 |
| mood | string | 否 | max=20 | trim |
| weather | string | 否 | max=20 | trim |
| location | string | 否 | max=100 | trim |

#### JSON Schema（DiaryPutReq）

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["title", "content"],
  "properties": {
    "title": { "type": "string", "minLength": 1, "maxLength": 100 },
    "content": { "type": "string", "minLength": 1 },
    "mood": { "type": "string", "maxLength": 20 },
    "weather": { "type": "string", "maxLength": 20 },
    "location": { "type": "string", "maxLength": 100 }
  }
}
```

### 3) 响应定义

#### 成功响应（200）

`{ "data": Diary }`

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败或 id 不合法 | ErrorResponse / `{"error":"id 参数不合法"}` |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 日记不存在 | `{"error":"日记不存在"}` |
| 500 | 服务器错误 | `{"error":"服务器内部错误"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DUP-001 | 已登录且日记存在 | 合法全量 body | 200 | data.title/content 等于输入 |
| DUP-002 | 已登录且日记存在 | title 长度=1,100 | 200 | 边界通过 |
| DUP-101 | 已登录 | id=abc | 400 | error=id 参数不合法 |
| DUP-102 | 已登录且日记存在 | 缺少 title | 400 | error 存在 |
| DUP-103 | 已登录且日记存在 | title 长度=0 | 400 | error 存在 |
| DUP-104 | 已登录且日记存在 | title 长度=101 | 400 | error 存在 |
| DUP-105 | 已登录且日记存在 | 缺少 content | 400 | error 存在 |
| DUP-201 | 已登录 | id 不存在 | 404 | error=日记不存在 |
| DUP-301 | 已登录且日记存在 | 相同 body 重复 PUT 两次 | 200 | 两次响应 data 等价（幂等） |
| DUP-302 | 压测环境 | 并发 20 PUT | 200/404 | RT P95 ≤ 500ms |

---

## 接口：PATCH /api/v1/diaries/:id

源码依据：[diary.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/diary.go#L209-L241)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/diaries/:id` |
| Method | `PATCH` |
| 功能 | 部分更新日记 |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 需要登录 |
| Content-Type | 是 | `application/json` | JSON 请求体 |

#### Path

| 参数 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| id | uint | 是 | 正整数（>0） | 非法返回 400 |

#### Body（DiaryPatchReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| title | string | 否 | min=1, max=100 | 传空串会触发校验失败 |
| content | string | 否 | 任意 | 允许空串（validator: omitempty） |
| mood | string | 否 | max=20 | - |
| weather | string | 否 | max=20 | - |
| location | string | 否 | max=100 | - |

> 注意：若 body 不包含任何可更新字段，将返回 `400` 且 `error=没有可更新的字段`。

### 3) 响应定义

#### 成功响应（200）

`{ "data": Diary }`

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败/无可更新字段/id 不合法 | ErrorResponse |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 日记不存在 | `{"error":"日记不存在"}` |
| 500 | 服务器错误 | `{"error":"服务器内部错误"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DPA-001 | 已登录且日记存在 | {"title":"new"} | 200 | data.title=new |
| DPA-002 | 已登录且日记存在 | {"content":""} | 200 | 允许空串 content（按实际返回断言） |
| DPA-003 | 已登录且日记存在 | {"mood":"m"*20} | 200 | 边界通过 |
| DPA-101 | 已登录且日记存在 | {} | 400 | error=没有可更新的字段 |
| DPA-102 | 已登录且日记存在 | {"title":""} | 400 | error 存在 |
| DPA-103 | 已登录且日记存在 | {"title":"t"*101} | 400 | error 存在 |
| DPA-104 | 已登录且日记存在 | {"mood":"m"*21} | 400 | error 存在 |
| DPA-201 | 已登录 | id 不存在 | 404 | error=日记不存在 |
| DPA-301 | 压测环境 | 并发 20 PATCH | 200/404 | RT P95 ≤ 500ms |

---

## 接口：DELETE /api/v1/diaries/:id

源码依据：[diary.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/diary.go#L243-L266)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/diaries/:id` |
| Method | `DELETE` |
| 功能 | 删除日记 |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 |
| --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` |

#### Path

| 参数 | 类型 | 必填 | 取值范围 |
| --- | --- | --- | --- |
| id | uint | 是 | 正整数（>0） |

### 3) 响应定义

#### 成功响应（200）

```json
{ "message": "删除成功" }
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | id 参数不合法 | `{"error":"id 参数不合法"}` |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 日记不存在 | `{"error":"日记不存在"}` |
| 500 | 服务器错误 | `{"error":"服务器内部错误"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DDL-001 | 已登录且日记存在 | id=101 | 200 | message=删除成功 |
| DDL-101 | 已登录 | id=abc | 400 | error=id 参数不合法 |
| DDL-201 | 已登录 | id 不存在 | 404 | error=日记不存在 |
| DDL-301 | 已登录 | 连续 DELETE 同一 id 两次 | 200 + 404 | 第一次删除成功；第二次日记不存在（非幂等实现） |
| DDL-302 | 压测环境 | 并发 20 DELETE | 200/404 | RT P95 ≤ 500ms |

---

## 接口：POST /api/v1/drafts

源码依据：[draft.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/draft.go#L50-L76)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/drafts` |
| Method | `POST` |
| 功能 | 创建草稿日记（Redis） |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 |
| --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` |
| Content-Type | 是 | `application/json` |

#### Body（DraftCreateReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| title | string | 是 | min=1, max=100 | trim |
| content | string | 是 | 非空 | - |
| mood | string | 否 | max=20 | trim |
| weather | string | 否 | max=20 | trim |
| location | string | 否 | max=100 | trim |

#### JSON Schema（DraftCreateReq）

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["title", "content"],
  "properties": {
    "title": { "type": "string", "minLength": 1, "maxLength": 100 },
    "content": { "type": "string", "minLength": 1 },
    "mood": { "type": "string", "maxLength": 20 },
    "weather": { "type": "string", "maxLength": 20 },
    "location": { "type": "string", "maxLength": 100 }
  }
}
```

### 3) 响应定义

#### 成功响应（201）

`{ "data": DraftDiary }`

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败 | ErrorResponse |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 500 | 创建失败 | `{"error":"创建草稿失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DFC-001 | 已登录 | title="t", content="c" | 201 | data.id>0；data.version>=1 |
| DFC-002 | 已登录 | title 长度=1/100 | 201 | 边界通过 |
| DFC-101 | 已登录 | 缺少 title | 400 | error 存在 |
| DFC-102 | 已登录 | title 长度=101 | 400 | error 存在 |
| DFC-103 | 已登录 | 缺少 content | 400 | error 存在 |
| DFC-104 | 已登录 | mood 长度=21 | 400 | error 存在 |
| DFC-201 | 未登录 | 无 Authorization | 401 | error=缺少Authorization头部 |
| DFC-301 | 压测环境 | 并发 20 创建 | 201 | RT P95 ≤ 500ms |

---

## 接口：GET /api/v1/drafts

源码依据：[draft.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/draft.go#L78-L105)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/drafts` |
| Method | `GET` |
| 功能 | 获取草稿列表（分页，返回 data 包裹） |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 |
| --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` |

#### Query

同 `/api/v1/diaries`：`page`、`page_size` 默认并 clamp。

### 3) 响应定义

#### 成功响应（200）

```json
{
  "data": {
    "items": [],
    "page": 1,
    "page_size": 20,
    "total": 0
  }
}
```

JSON Schema：

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": true,
  "required": ["data"],
  "properties": {
    "data": {
      "type": "object",
      "additionalProperties": true,
      "required": ["items", "page", "page_size", "total"],
      "properties": {
        "items": { "type": "array", "items": { "$ref": "https://catdiary.local/schema/DraftDiary.json" } },
        "page": { "type": "integer", "minimum": 1 },
        "page_size": { "type": "integer", "minimum": 1, "maximum": 100 },
        "total": { "type": "integer", "minimum": 0 }
      }
    }
  }
}
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 500 | 获取失败 | `{"error":"获取草稿列表失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DFL-001 | 已登录 | page=1&page_size=20 | 200 | data.items 为数组；data.page/page_size 回显 |
| DFL-002 | 已登录且有 21 条 | page=2&page_size=20 | 200 | data.items 数量=1 |
| DFL-003 | 已登录 | page=abc | 200 | data.page=1（默认） |
| DFL-004 | 已登录 | page_size=101 | 200 | data.page_size=100（clamp） |
| DFL-201 | 未登录 | 无 Authorization | 401 | error=缺少Authorization头部 |
| DFL-301 | 压测环境 | 并发 50 列表 | 200 | RT P95 ≤ 300ms |

---

## 接口：GET /api/v1/drafts/:id

源码依据：[draft.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/draft.go#L107-L140)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/drafts/:id` |
| Method | `GET` |
| 功能 | 获取草稿详情 |

### 2) 请求定义

| Header | 必填 | 示例 |
| --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` |

| Path 参数 | 类型 | 必填 | 取值范围 |
| --- | --- | --- | --- |
| id | uint64 | 是 | 正整数（>0） |

### 3) 响应定义

#### 成功响应（200）

`{ "data": DraftDiary }`

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | id 参数不合法 | `{"error":"id 参数不合法"}` |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 草稿不存在 | `{"error":"草稿不存在"}` |
| 500 | 获取失败 | `{"error":"获取草稿失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DFG-001 | 已登录且草稿存在 | id=1 | 200 | data.id=id；data.version>=1 |
| DFG-101 | 已登录 | id=0 | 400 | error=id 参数不合法 |
| DFG-102 | 已登录 | id=abc | 400 | error=id 参数不合法 |
| DFG-201 | 已登录 | id 不存在 | 404 | error=草稿不存在 |
| DFG-301 | 压测环境 | 并发 50 获取 | 200/404 | RT P95 ≤ 300ms |

---

## 接口：PUT /api/v1/drafts/:id

源码依据：[draft.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/draft.go#L142-L190)、
[draft.go service](file:///d:/WorkSpace/Go/CatDiary/backend/internal/service/draft.go#L83-L107)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/drafts/:id` |
| Method | `PUT` |
| 功能 | 全量更新草稿（日记草稿支持乐观锁） |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 |
| --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` |
| Content-Type | 是 | `application/json` |

#### Path

| 参数 | 类型 | 必填 | 取值范围 |
| --- | --- | --- | --- |
| id | uint64 | 是 | 正整数（>0） |

#### Body（DraftPutReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| expected_version | uint64 | 否 | >=1 | 提供时启用乐观锁；版本不匹配返回 409 |
| title | string | 是 | min=1, max=100 | trim |
| content | string | 是 | 非空 | - |
| mood | string | 否 | max=20 | trim |
| weather | string | 否 | max=20 | trim |
| location | string | 否 | max=100 | trim |

#### JSON Schema（DraftPutReq）

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["title", "content"],
  "properties": {
    "expected_version": { "type": "integer", "minimum": 1 },
    "title": { "type": "string", "minLength": 1, "maxLength": 100 },
    "content": { "type": "string", "minLength": 1 },
    "mood": { "type": "string", "maxLength": 20 },
    "weather": { "type": "string", "maxLength": 20 },
    "location": { "type": "string", "maxLength": 100 }
  }
}
```

### 3) 响应定义

#### 成功响应（200）

`{ "data": DraftDiary }`

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败/id 不合法 | ErrorResponse |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 草稿不存在 | `{"error":"草稿不存在"}` |
| 409 | 版本冲突 | DraftConflictResponse |
| 500 | 更新失败 | `{"error":"更新草稿失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DFU-001 | 已登录且草稿存在 | 合法全量 body | 200 | data.title/content 等于输入；version 自增或变更 |
| DFU-002 | 已登录且草稿存在 | expected_version=当前版本 | 200 | 成功；version 更新 |
| DFU-101 | 已登录 | 缺少 title | 400 | error 存在 |
| DFU-102 | 已登录 | content="" | 400 | error 存在 |
| DFU-103 | 已登录 | expected_version=0 | 400/200 | 若 validator 不拦截则按实际；建议断言不接受 0 |
| DFU-201 | 已登录 | id 不存在 | 404 | error=草稿不存在 |
| DFU-301 | 已登录且草稿存在 | 两客户端并发：A 读 v=1，B PUT 成功变 v=2，A PUT expected_version=1 | 409 | error=草稿已被更新；current_version=2 |
| DFU-302 | 已登录且草稿存在 | 相同 body + expected_version 未提供，重复 PUT 两次 | 200 | 允许；需确认字段一致，version 可能变化 |
| DFU-303 | 压测环境 | 并发 20 PUT | 200/409/404 | RT P95 ≤ 500ms |

---

## 接口：PATCH /api/v1/drafts/:id

源码依据：[draft.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/draft.go#L192-L247)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/drafts/:id` |
| Method | `PATCH` |
| 功能 | 部分更新草稿（支持 expected_version 乐观锁） |

### 2) 请求定义

#### Headers

| Header | 必填 | 示例 |
| --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` |
| Content-Type | 是 | `application/json` |

#### Path

| 参数 | 类型 | 必填 | 取值范围 |
| --- | --- | --- | --- |
| id | uint64 | 是 | 正整数（>0） |

#### Body（DraftPatchReq）

| 字段 | 类型 | 必填 | 取值范围 | 业务规则 |
| --- | --- | --- | --- | --- |
| expected_version | uint64 | 否 | >=1 | 版本不匹配返回 409 |
| title | string | 否 | min=1, max=100 | - |
| content | string | 否 | min=1 | 此接口 content 不允许空串（validate:min=1） |
| mood | string | 否 | max=20 | - |
| weather | string | 否 | max=20 | - |
| location | string | 否 | max=100 | - |

> 注意：若 body 不包含任何可更新字段，将返回 `400` 且 `error=没有可更新的字段`。

### 3) 响应定义

#### 成功响应（200）

`{ "data": DraftDiary }`

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | 参数校验失败/无可更新字段/id 不合法 | ErrorResponse |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 草稿不存在 | `{"error":"草稿不存在"}` |
| 409 | 版本冲突 | DraftConflictResponse |
| 500 | 更新失败 | `{"error":"更新草稿失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DFP-001 | 已登录且草稿存在 | {"title":"new"} | 200 | data.title=new |
| DFP-002 | 已登录且草稿存在 | {"content":"c"} | 200 | data.content=c |
| DFP-003 | 已登录且草稿存在 | {"expected_version":1,"title":"x"} | 200/409 | 若版本匹配则 200，否则 409 |
| DFP-101 | 已登录且草稿存在 | {} | 400 | error=没有可更新的字段 |
| DFP-102 | 已登录且草稿存在 | {"content":""} | 400 | error 存在（min=1） |
| DFP-103 | 已登录 | id=abc | 400 | error=id 参数不合法 |
| DFP-201 | 已登录 | id 不存在 | 404 | error=草稿不存在 |
| DFP-301 | 已登录且草稿存在 | 并发：A 读 v=1，B PATCH 成功变 v=2，A PATCH expected_version=1 | 409 | current_version=2 |
| DFP-302 | 压测环境 | 并发 20 PATCH | 200/409/404 | RT P95 ≤ 500ms |

---

## 接口：DELETE /api/v1/drafts/:id

源码依据：[draft.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/draft.go#L249-L280)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/drafts/:id` |
| Method | `DELETE` |
| 功能 | 删除草稿 |

### 2) 请求定义

| Header | 必填 | 示例 |
| --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` |

| Path 参数 | 类型 | 必填 | 取值范围 |
| --- | --- | --- | --- |
| id | uint64 | 是 | 正整数（>0） |

### 3) 响应定义

#### 成功响应（200）

```json
{ "message": "删除成功" }
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | id 参数不合法 | `{"error":"id 参数不合法"}` |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 草稿不存在 | `{"error":"草稿不存在"}` |
| 500 | 删除失败 | `{"error":"删除草稿失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| DFD-001 | 已登录且草稿存在 | id=1 | 200 | message=删除成功 |
| DFD-101 | 已登录 | id=0 | 400 | error=id 参数不合法 |
| DFD-201 | 已登录 | id 不存在 | 404 | error=草稿不存在 |
| DFD-301 | 已登录 | 连续 DELETE 两次 | 200 + 404 | 第二次草稿不存在（非幂等实现） |
| DFD-302 | 压测环境 | 并发 20 DELETE | 200/404 | RT P95 ≤ 500ms |

---

## 接口：POST /api/v1/drafts/:id/flush

源码依据：[draft.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/draft.go#L282-L312)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/drafts/:id/flush` |
| Method | `POST` |
| 功能 | 触发草稿落库（异步写入快照） |

### 2) 请求定义

| Header | 必填 | 示例 |
| --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` |

| Path 参数 | 类型 | 必填 | 取值范围 |
| --- | --- | --- | --- |
| id | uint64 | 是 | 正整数（>0） |

### 3) 响应定义

#### 成功响应（200）

```json
{ "message": "触发落库成功" }
```

#### 错误响应

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 400 | id 参数不合法 | `{"error":"id 参数不合法"}` |
| 401 | 未登录/Token 无效 | ErrorResponse |
| 404 | 草稿不存在 | `{"error":"草稿不存在"}` |
| 500 | 触发失败 | `{"error":"触发落库失败"}` |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| FLH-001 | 已登录且草稿存在 | id=1 | 200 | message=触发落库成功 |
| FLH-101 | 已登录 | id=abc | 400 | error=id 参数不合法 |
| FLH-201 | 已登录 | id 不存在 | 404 | error=草稿不存在 |
| FLH-301 | 压测环境 | 并发 20 flush | 200/404 | RT P95 ≤ 300ms |

### 5) 覆盖要求补充

- 幂等性：连续触发 flush 可重复成功（不保证同步落库完成）
- 并发：在持续 PATCH 草稿的同时触发 flush，验证快照版本单调递增（以 DB 数据为准）

---

## 接口：POST /api/v1/uploads

源码依据：[upload.go](file:///d:/WorkSpace/Go/CatDiary/backend/internal/handler/upload.go)

### 1) 基本信息

| 项 | 值 |
| --- | --- |
| Path | `/api/v1/uploads` |
| Method | `POST` |
| 功能 | 创建上传（当前实现：未实现，直接返回 501） |

### 2) 请求定义

> 当前实现未定义请求体；建议未来明确是 `multipart/form-data` 还是 JSON。

#### Headers

| Header | 必填 | 示例 | 说明 |
| --- | --- | --- | --- |
| Authorization | 是 | `Bearer <jwt>` | 路由层要求登录 |

### 3) 响应定义

#### 当前响应（501）

| 状态码 | 响应头 | 响应体 |
| --- | --- | --- |
| 501 | `Content-Type: text/plain; charset=utf-8` | `TODO` |

#### 错误响应（来自认证中间件）

| 状态码 | 场景 | 响应体 |
| --- | --- | --- |
| 401 | 未登录/Token 无效 | ErrorResponse |

### 4) 字段测试用例

| 用例编号 | 前置条件 | 输入数据 | 预期结果 | 断言要点 |
| --- | --- | --- | --- | --- |
| UP-001 | 已登录 | 任意请求体 | 501 | body=TODO |
| UP-101 | 未登录 | 无 Authorization | 401 | error=缺少Authorization头部 |
| UP-301 | 压测环境 | 并发 50 | 501 | RT P95 ≤ 300ms |

---

## 静态检查与交付清单

### markdownlint（建议命令）

```bash
npx -y markdownlint-cli2 "docs/api-test/**/*.md"
```

### Schema 校验（建议命令）

将文档中 JSON Schema 代码块抽取为 `.json` 文件后，用 AJV 校验：

```bash
npx -y ajv-cli validate -s ./schema.json -d ./example.json
```

### 团队评审检查点

| 检查项 | 要求 |
| --- | --- |
| 与源码同步 | 路由、字段、状态码与 handler/service 实现一致 |
| 零歧义 | 规则表述可直接转换为自动化断言 |
| 可维护性 | 组件化 Schema + 章节化用例编号 |
