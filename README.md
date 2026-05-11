# 🐱 CatDiary（猫咪日记）

> 喵～欢迎来到猫咪日记！这是本喵用 Go 和 Vue 精心打造的全栈日记小窝，让你可以安安静静地写日记、存草稿、贴图片，还能用 Markdown 写得漂漂亮亮的喵～

## ✨ 这里有什么好玩的喵？

- **沉浸式写作喵** 🖊️ — Markdown 编辑器加持，写日记就像在云朵上踩奶一样舒服～
- **双重保护喵** 🛡️ — 热数据在 Redis 里打滚，冷数据在 MySQL 里睡大觉，后台自动把草稿从 Redis 搬到 MySQL，不怕丢喵！
- **乐观锁防冲突喵** 🔒 — 草稿编辑带了 version 字段，两只猫同时挠同一个草稿也不会打架～
- **图片上传喵** 📸 — 可以把猫猫照片上传到日记里，存在本地 `data/uploads/` 目录下～
- **JWT 认证喵** 🔑 — 无状态 token 认证，7 天有效，不用每次登录都要按爪印～
- **AI 预留喵** 🤖 — 架构上给大模型留了位置，以后可以让 AI 帮你给日记打标签、分析心情、写年度总结～

## 🛠️ 技术栈小鱼干

| 层级 | 吃的东西 |
|------|----------|
| 后端 | Go 1.26 + CloudWeGo Hertz + GORM + go-redis |
| 前端 | Vue 3 + Vite + Vue Router 4 + Tailwind CSS 4 + daisyUI 5 |
| 数据库 | MySQL 8.0（主存储）+ Redis 7.4（草稿热存储） |
| 基础设施 | Docker Compose 一键起喵～ |

## 🚀 快速开始

### 环境要求

| 工具 | 说明 |
|------|------|
| Go 1.26+ | 跑后端喵 |
| Node.js 20+ | 跑前端喵 |
| MySQL 8.0 + Redis 7.4 | 有 Docker 就自动搞定；没有就自己装喵～ |

---

### 路径一：有 Docker（推荐，两爪子搞定 🐾）

**第一步** — 复制环境变量文件，改一下 MySQL 密码：

```bash
cp .env.example .env
# 编辑 .env，把 CATDIARY_MYSQL_DSN 里的 password 改成你 Docker MySQL 的密码
```

**第二步** — 一键起飞：

```bash
make dev           # Docker 拉起来 → 后端启起来

# 开另一个终端跑前端：
cd frontend && npm install && npm run dev
```

后端 `http://localhost:8080`，前端 `http://localhost:5173`，搞定了喵～

---

### 路径二：没有 Docker（需要本地装好 MySQL + Redis）

**第一步** — 同上，复制 `.env.example` 为 `.env`，把 `CATDIARY_MYSQL_DSN` 指向你本地的 MySQL。

**第二步** — 确保 MySQL 和 Redis 在本地跑着，然后：

```bash
# 终端 1：启动后端
make run           # 或者 go run ./backend/cmd/api/main.go

# 终端 2：启动前端
cd frontend && npm install && npm run dev
```

Docker 只是方便起依赖而已，不用 Docker 也能跑，把 MySQL 和 Redis 地址配对就行喵～

## 📁 猫窝结构一览

```text
CatDiary/
├── backend/                      # Go 后端（本喵的窝）
│   ├── cmd/api/main.go           # 入口：从配置文件开始，一路初始化到路由
│   └── internal/
│       ├── config/               # 环境变量加载（godotenv）
│       ├── model/                # GORM 数据模型（User, Diary, DraftDiary）
│       ├── repository/           # 数据访问层（MySQL + Redis + Lua 脚本）
│       ├── service/              # 业务逻辑层（所有聪明的小脑瓜都在这）
│       ├── handler/              # HTTP 处理层（接收请求、校验、调用 service）
│       ├── middleware/           # JWT 认证中间件
│       ├── router/               # 路由注册 + 静态文件服务
│       └── worker/               # 后台草稿刷新器（200ms 轮询 Redis → MySQL）
├── frontend/                     # Vue 3 前端（猫爬架）
│   └── src/
│       ├── router/index.js       # 10 条路由，带认证守卫
│       ├── services/             # API 封装 + Token 管理
│       └── views/                # 页面组件（着陆页、登录注册、日记、草稿、个人资料）
├── docker-compose.yml            # MySQL 8.0 + Redis 7.4 编排
├── Makefile                      # 开发命令小助手
└── docs/api-test/                # 手动 API 测试用例文档（没有自动化测试喵…）
```

## 🐾 API 速览

所有接口都在 `/api/v1` 下面喵：

| 分类 | 路由 | 说喵 |
|------|------|------|
| 健康检查 | `GET /healthz`, `GET /readyz` | K8s 探针戳一戳～ |
| 认证 | `POST /auth/register`, `POST /auth/login` | 注册登录，拿到 JWT 小鱼干 |
| 用户 | `GET /users/me`, `PATCH /users/me` | 看资料、改资料、改密码 |
| 日记 | `CRUD /diaries` | 增删改查，支持 `happened_at` 字段记事情发生的日子 |
| 草稿 | `CRUD /drafts` | Redis 热存储 + 乐观锁 + 后台落盘 MySQL |
| 上传 | `POST /uploads`, `GET /uploads/:filename` | 传图片喵～ |

## 🧶 数据流

```
浏览器 → Vue Router → API Service (fetch) → Hertz Handler → Service → Repository → GORM/Redis → MySQL/Redis
                                                           ↑ 校验 (validator/v10)
                                                           ↑ JWT 中间件验 token
```

## 📝 开发约定

- **提交信息**: 中文，每行以"喵~"结尾，祈使语气，标题 ≤72 字（详见 `.trae/rules/git-commit-message.md`）
- **表名**: GORM Model 用单数表名（`TableName()` 显式指定）
- **用户隔离**: 所有日记/草稿操作按 `user_id` 隔离，每只猫只能碰自己的东西～

## ⚠️ 还没做完的事情

- [ ] `.env.example` 示例文件还没创建喵…
- [ ] 没有自动化测试（`*_test.go` 不存在），现在靠 `docs/api-test/` 里的手动测试文档撑着
- [ ] `pkg/eino/` LLM 集成还只是预留的空壳
- [ ] `pkg/logger/` 结构化日志还没接入

---

> 喵～如果你喜欢这个项目，就给颗小鱼干（Star）吧！有问题欢迎提 Issue，本喵会抽空挠键盘回复的～ 🐾
