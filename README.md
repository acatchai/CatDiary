# 🐱 CatDiary (猫咪日记)

> 一个基于 Go 构建的高可用、全栈日记管理系统，支持沉浸式 Markdown 写作、多级数据落盘与大模型（LLM）智能化扩展。

## 🌟 项目特性

* **沉浸式写作**：集成现代 Markdown 编辑器，支持流畅的文本与多媒体输入。
* **高可用存储架构**：
    * **持久层**：MySQL 8.0 (原生支持 `utf8mb4` 与 Emoji 表情)。
    * **缓存层**：Redis 7.4 (AOF 持久化开启)，用于高性能鉴权与草稿防丢失。
    * **容灾层**：后端支持 JSONL 本地文件追加写入，应对数据库瞬时故障。
* **安全鉴权**：基于 Cache-Aside 模式的自定义 Session 机制，保障 API 安全。
* **AI 赋能预留**：架构级预留大模型接入层，未来支持日记自动标签化、情感分析与年度记忆总结。

## 🛠️ 技术栈选型

* **后端 (Backend)**: Go 1.21+, Hertz / Gin, GORM, go-redis, Viper (配置管理).
* **前端 (Frontend)**: React 18 / Vue 3, Vite, Tailwind CSS, shadcn/ui, Vditor / Milkdown.
* **基础设施 (Infra)**: Docker & Docker Compose, WSL2.

---

## 📁 核心目录结构与模块说明

本项目采用 **Monorepo（单体仓库）** 结构，后端遵循 Go Standard Project Layout 规范。

```text
diary-project/
├── backend/                  # Go 后端微单体服务
│   ├── cmd/api/main.go       # [入口点] 服务的启动引导程序，负责依赖注入和启动 Server
│   ├── internal/             # [核心业务] 外部应用无法导入的私有核心代码
│   │   ├── config/           # 配置解析：读取 .env 并反序列化为 Go Struct
│   │   ├── handler/          # 表现层：解析 HTTP Request，调用 Service，返回 JSON Response
│   │   ├── service/          # 业务逻辑层：核心算法、数据流转（校验、核心计算）
│   │   ├── repository/       # 数据访问层：封装 MySQL、Redis、JSONL 的具体读写操作
│   │   ├── model/            # 领域模型：GORM 结构体 (Entities) 与接口 DTOs
│   │   ├── middleware/       # 中间件：Session 拦截校验、全局错误捕获、CORS
│   │   └── router/           # 路由层：API 路径定义与 Handler 绑定
│   └── pkg/                  # [公共库] 可供其他项目复用的通用组件
│       ├── errors/           # 错误处理上下文封装
│       └── logger/           # 结构化日志封装
├── frontend/                 # 前端视图服务 (Vite)
├── deploy/                   # 基础设施部署脚本与初始化 SQL
├── docker-compose.yml        # 核心依赖 (MySQL + Redis) 容器编排
├── Makefile                  # 本地开发环境快捷命令集合
└── .env.example              # 环境变量配置模板
