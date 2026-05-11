# CLAUDE.md

本文件为 Claude Code（claude.ai/code）在此仓库中工作时提供指导。

## 命令

### 开发环境
- `make dev` — 启动 Docker 依赖（MySQL + Redis）然后运行 Go API 服务
- `make up` — `docker compose up -d`（仅启动依赖服务）
- `make run` — `cd backend && go run ./cmd/api/main.go`（仅运行 API，不启动 Docker）
- `make down` — 停止并移除 Docker 容器
- `make logs` — 查看 Docker Compose 日志（跟随模式）
- `make restart` — 重启 Docker 服务（stop + up）
- `make stop` — 停止 Docker 服务（不删除容器）

不使用 Makefile：
- `go run ./backend/cmd/api/main.go` — 从项目根目录运行 API
- `cd backend && go run ./cmd/api/main.go` — 从 backend 目录运行

### 前端
- `cd frontend && npm install` — 安装前端依赖
- `cd frontend && npm run dev` — 启动 Vite 开发服务器
- `cd frontend && npm run build` — 生产构建

### 构建
无独立构建脚本，标准 `go build`：在 `backend/` 下执行 `go build ./cmd/api/main.go`。

### 依赖管理
- `cd backend && go mod tidy` — 同步 Go 模块依赖
- `cd backend && go mod download` — 下载所有依赖

### 测试
项目目前没有任何自动化测试文件（`*_test.go`）。手动 API 测试用例文档位于 `docs/api-test/API_TEST_SPEC_BACKEND_V1.md`。

---

## 项目概览

**CatDiary（猫咪日记）** — 全栈日记管理系统，支持沉浸式 Markdown 写作、多级数据持久化和 LLM 扩展。

### 技术栈
- **后端**: Go 1.26, CloudWeGo Hertz（HTTP 框架）, GORM（ORM）, go-redis, JWT 认证
- **前端**: Vue 3 + Vite + Vue Router 4, Tailwind CSS 4 + daisyUI 5
- **基础设施**: Docker Compose（MySQL 8.0 + Redis 7.4）

### 目录结构

```
CatDiary/
├── .env                          # 环境变量（gitignored，仅含 CATDIARY_MYSQL_DSN）
├── Makefile                      # 开发命令快捷方式
├── docker-compose.yml            # MySQL 8.0 + Redis 7.4
├── backend/                      # Go 后端
│   ├── go.mod / go.sum
│   ├── cmd/api/main.go           # 入口：加载配置 → 初始化 DB/Redis → 注册路由 → 启动服务
│   └── internal/
│       ├── config/env.go         # godotenv 加载 .env
│       ├── model/                # GORM 模型（User, Diary, DraftDiary）
│       ├── repository/           # 数据访问层
│       │   ├── db.go             # GORM 初始化 + AutoMigrate
│       │   ├── redis.go          # Redis 客户端初始化
│       │   ├── user.go           # 用户 MySQL CRUD
│       │   ├── diary.go          # 日记 MySQL CRUD（按 user_id 隔离）
│       │   ├── draft.go          # 草稿 Redis CRUD（Lua 脚本原子操作）
│       │   └── draft_mysql.go    # 草稿 MySQL 落盘（upsert/delete）
│       ├── service/              # 业务逻辑层，返回领域错误
│       ├── handler/              # HTTP 处理层（Hertz handler）
│       │   ├── health.go         # 健康检查
│       │   ├── auth.go           # 注册/登录/登出
│       │   ├── user.go           # 用户信息与密码管理
│       │   ├── diary.go          # 日记 CRUD
│       │   ├── draft.go          # 草稿 CRUD
│       │   └── upload.go         # 图片上传
│       ├── middleware/auth.go    # JWT Bearer 认证中间件
│       ├── router/router.go      # 路由注册 + 静态文件服务
│       └── worker/
│           └── draft_flusher.go  # 后台草稿刷新器（Redis → MySQL，200ms 轮询）
│   └── pkg/
│       ├── utils/jwt.go          # JWT 生成/解析（HS256，7 天过期）
│       ├── eino/                 # （预留）LLM 集成
│       ├── errors/               # （预留）错误处理
│       └── logger/               # （预留）结构化日志
├── frontend/                     # Vue 3 前端
│   ├── vite.config.js
│   ├── tailwind.config.js
│   └── src/
│       ├── main.js / App.vue / main.css
│       ├── router/index.js       # 10 条路由，含 auth guard
│       ├── services/
│       │   ├── api.js            # 通用 fetch 封装（自动拼接 /api/v1，自动带 Bearer token）
│       │   └── auth.js           # Token localStorage 存取
│       ├── views/                # 页面组件
│       │   ├── LandingPage.vue   # 公开落地页
│       │   ├── LoginPage.vue / RegisterPage.vue
│       │   └── app/              # 需认证的页面（AppLayout 壳 + 子路由）
│       └── components/           # 公共组件
├── frontend-study/               # 前期前端原型（gitignored）
├── docs/api-test/                # API 测试用例文档
├── deploy/mysql/ / deploy/redis/ # 部署预留（空目录）
└── data/                         # Docker 数据卷（gitignored）
```

### 后端分层数据流

`Handler（请求/响应） → Service（业务逻辑） → Repository（数据访问） → GORM/Redis → MySQL/Redis`

- **Handler**: 解析请求、用 `go-playground/validator/v10` 做结构体校验、调用 Service、返回 JSON 响应
- **Service**: 包含全部业务逻辑，返回领域错误（如 `ErrUserNotFound`、`ErrDiaryNotFound`、`ErrDraftConflict`）
- **Repository**: MySQL 操作通过 GORM，Redis 操作通过 go-redis（草稿操作使用 Lua 脚本保证原子性）
- **Model**: GORM 实体，显式定义 `TableName()` 返回单数表名

### API 路由（全部在 `/api/v1` 下）

| 分组 | 需认证 | 路由 |
|------|--------|------|
| `/healthz`, `/readyz` | 否 | GET（存活/就绪探针） |
| `/auth` | 部分 | POST register/login（公开），POST logout + GET me（需认证） |
| `/users` | 是 | GET me, PATCH me, PATCH me/password |
| `/diaries` | 是 | POST create, GET list（分页 `?page=&page_size=`）, GET/:id, PUT/:id, PATCH/:id, DELETE/:id |
| `/drafts` | 是 | POST create, GET list（分页）, GET/:id, PUT/:id, PATCH/:id, DELETE/:id, POST/:id/flush |
| `/uploads` | 是 | POST create（图片上传），GET 静态文件服务 `/uploads/:filename` |

### 关键设计决策

1. **JWT 认证**: 无状态 token 认证。`middleware.RequireAuth()` 解析 Bearer token，将 `user_id` 注入请求上下文。Token 有效期 7 天，签名算法 HS256。
2. **参数校验**: Hertz 服务配置了 `go-playground/validator/v10`，通过 struct tag 声明校验规则（如 `validate:"required,min=3,max=50"`）。
3. **日记 CRUD**: 所有日记操作按 `user_id` 隔离，用户只能操作自己的日记。PUT 为全量替换，PATCH 为部分更新。日记支持 `happened_at` 字段记录事情发生时间。
4. **草稿系统**: Redis 作为热存储，支持乐观锁（version 字段防并发冲突）。后台 flusher 每 200ms 将脏草稿异步落盘到 MySQL。活跃草稿 TTL 30 天，删除草稿 TTL 1 小时。所有 Redis 写操作使用 Lua 脚本保证原子性。
5. **图片上传**: 支持图片上传到本地 `data/uploads/` 目录，通过 Hertz StaticFS 提供静态文件服务。
6. **自动迁移**: 启动时 GORM AutoMigrate 自动同步数据库 schema，修改 Model 后无需手动迁移。
7. **无 session 机制**: 后端完全无状态，不存储 session。登出仅为前端清除 token 的占位操作。

### 环境变量

| 变量 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `CATDIARY_MYSQL_DSN` | 是 | 无 | MySQL 连接串（写在 `.env` 中） |
| `CATDIARY_ENV_FILE` | 否 | 无 | 自定义 .env 文件路径 |
| `CATDIARY_REDIS_ADDR` | 否 | `127.0.0.1:6379` | Redis 地址 |
| `CATDIARY_REDIS_PASSWORD` | 否 | 空 | Redis 密码 |
| `CATDIARY_REDIS_DB` | 否 | `0` | Redis 数据库编号 |
| `JWT_SECRET` | 否 | `catdiary_default_secret_key_please_change` | JWT 签名密钥 |

### 前端路由

| 路径 | 页面 | 认证 | 说明 |
|------|------|------|------|
| `/` | LandingPage | 否 | 公开落地页（猫咪主题） |
| `/login` | LoginPage | 否 | 登录 |
| `/register` | RegisterPage | 否 | 注册 |
| `/app` | AppLayout | 是 | 认证壳，重定向到 `/app/diaries` |
| `/app/diaries` | DiaryListPage | 是 | 日记列表 |
| `/app/diaries/new` | DiaryNewPage | 是 | 新建日记 |
| `/app/diaries/:id` | DiaryDetailPage | 是 | 日记详情（查看/编辑/删除） |
| `/app/drafts` | DraftListPage | 是 | 草稿列表 |
| `/app/drafts/:id` | DraftEditPage | 是 | 草稿编辑 |
| `/app/profile` | ProfilePage | 是 | 个人资料 |
| `/app/settings/security` | SecurityPage | 是 | 安全设置 |

### 已知问题

1. **`.env.example` 不存在**: README 中提到的示例文件尚未创建。

### 提交规范

遵循 `.trae/rules/git-commit-message.md`：中文提交信息，每行以"喵~"结尾，不含 emoji，标题不超过 72 字，祈使语气。常用前缀：新增/修复/调整/重构/清理/优化。
