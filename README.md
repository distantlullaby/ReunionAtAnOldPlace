# 📮 记忆连接 ReunionAtAnOldPlace

> **你帮我再看一眼，我把记忆还给你**

帮助离城青年，通过他人的镜头"代看"故乡的回忆角落。贴上老照片与回忆发布"求看"并悬赏**记忆硬币**，此刻路过那里的人拍下现场新照、附上寄语，发起人确认后悬赏硬币结算给对方——一张由「过去回忆+老照片」与「当下现场+新照片」组成的**双面明信片**就此重逢。

## ✨ 核心功能

- **首页 Feed 流**：求看需求列表，卡片可展开为 3D 翻转的**双面明信片**
  - 过去面：回忆文字 + 老照片
  - 当下面：他人上传的现场新照 + 寄语
- **记忆硬币循环生态**
  - 注册即赠送初始硬币（默认 100）
  - 发布求看 → 从余额**冻结扣除**悬赏（支持**追加**）
  - 他人「替他去拍」上传新照与寄语
  - 发起人「确认采纳」→ 托管硬币**结算**给拍摄者
- **个人中心**：硬币余额、我发布的、我去拍的、完整硬币流水
- 账目安全：发帖冻结 / 追加 / 采纳结算均在 **Go 服务端数据库事务 + 行级锁(`SELECT … FOR UPDATE`)** 中完成，杜绝透支、重复结算与并发错乱

## 🧱 技术栈

| 端 | 技术 |
|---|---|
| 前端 | Vue 3 + Vite + Vue Router + Axios |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8（utf8mb4），JWT 鉴权，bcrypt 密码 |
| 图片 | 后端本地存储 + 静态服务，前端上传组件 |

## 📁 目录结构

```
ReunionAtAnOldPlace/
├─ backend/
│  ├─ cmd/server/main.go      # 服务入口
│  ├─ cmd/seed/main.go        # 可选：生成自洽的演示数据
│  └─ internal/
│     ├─ config/              # 环境变量配置
│     ├─ db/                  # 建库连接 + AutoMigrate
│     ├─ model/               # user / story / response / coin_transaction
│     ├─ service/             # 业务逻辑（硬币事务核心）
│     ├─ handler/             # HTTP 处理器
│     ├─ middleware/          # JWT 认证
│     ├─ router/              # 路由与依赖装配
│     └─ utils/               # JWT / 响应 / 密码
└─ frontend/
   └─ src/
      ├─ api/                 # axios 封装与接口
      ├─ store/auth.js        # 登录态
      ├─ router/
      ├─ components/          # PostcardCard 双面明信片、ImageUploader、Modal
      └─ views/               # FeedView 首页、LoginView、ProfileView
```

## 🗄️ 数据库表

- `users`：用户，含 `coin_balance` 可用余额
- `stories`：求看（明信片过去面），含 `reward` 托管悬赏、`status`(open/completed)、`accepted_response_id`
- `responses`：替拍回应（明信片当下面），含 `new_photo`、`message`
- `coin_transactions`：硬币流水（register_gift / freeze / append_freeze / reward_income / refund），记录金额与变动后余额

## 🚀 快速开始

### 0. 前置

- Go 1.22（如在 `C:\mine\code\tool\Google\go\go1.22.1\sdk\bin`）
- Node.js（nvm 管理）
- MySQL 已启动（默认 `192.168.5.100:3306`，root/123456）。库 `memory_link` 由后端**自动创建**并迁移表结构

### 1. 启动后端（:8080）

```bash
cd backend
go run ./cmd/server
# 或：go build -o server ./cmd/server && ./server
```

### 2. 启动前端（:5173）

```bash
cd frontend
npm install
npm run dev
```

浏览器打开 http://localhost:5173 ，`/api`、`/uploads` 已由 Vite 代理到 8080。

### 3.（可选）生成演示数据

会**清空并重建**一套硬币账目自洽的演示数据（含自动生成的示意图）：

```bash
cd backend
go run ./cmd/seed      # 或 go build -o seed ./cmd/seed && ./seed
```

演示账号（密码均为 `secret123`）：`alice`(阿篱·发起人)、`bob`(小城拍摄员·拍摄者)、`chuan`(阿川)。

## ⚙️ 配置（环境变量，均有默认值）

| 变量 | 默认 | 说明 |
|---|---|---|
| `SERVER_PORT` | 8080 | 服务端口 |
| `DB_HOST` / `DB_PORT` | 192.168.5.100 / 3306 | MySQL |
| `DB_USER` / `DB_PASSWORD` | root / 123456 | 账号密码 |
| `DB_NAME` | memory_link | 库名 |
| `INITIAL_COINS` | 100 | 注册赠送硬币 |
| `JWT_SECRET` | （内置） | 生产环境请覆盖 |
| `UPLOAD_DIR` | uploads | 图片存储目录 |

## 🔌 主要 API

| 方法 | 路径 | 鉴权 | 说明 |
|---|---|---|---|
| POST | /api/auth/register | – | 注册并赠送初始硬币 |
| POST | /api/auth/login | – | 登录 |
| GET | /api/stories | – | Feed 分页（?status=&page=&size=） |
| GET | /api/stories/:id | – | 故事详情含全部回应 |
| POST | /api/stories | ✅ | 发布求看（冻结悬赏） |
| POST | /api/stories/:id/append | ✅ | 追加悬赏冻结 |
| POST | /api/stories/:id/respond | ✅ | 替他去拍 |
| POST | /api/stories/:id/accept/:rid | ✅ | 确认采纳并结算 |
| POST | /api/upload | ✅ | 上传图片 |
| GET/PUT | /api/user/profile | ✅ | 个人中心 / 更新资料 |

统一响应：`{ "code": 0, "msg": "ok", "data": ... }`，鉴权头 `Authorization: Bearer <token>`。
