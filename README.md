# 记忆连接 · ReunionAtAnOldPlace（分支B）

> 你帮我再看一眼，我把记忆还给你。

帮助离城青年通过他人的镜头，代看回忆角落的 Web 应用。

## 功能

- **首页 Feed 流**：展示所有"求看"需求，可展开**双面明信片**——
  正面是「过去回忆 + 老照片」，翻面是「当下现场 + 新照片」。
- **个人中心**：记忆硬币余额（可用 / 冻结）、我的发布、我的回应、硬币流水。
- **记忆硬币循环生态**：
  - 注册即赠送 100 枚初始硬币；
  - 发布求看需求时冻结悬赏硬币（可随时追加）；
  - 他人点击「替他去拍」上传现场新照与寄语；
  - 发起人确认采纳后，冻结的悬赏硬币结算给拍摄者；
  - 取消需求则冻结硬币全额退回。

## 技术栈

| 端 | 技术 |
|----|------|
| 前端 | Vue 3 + Vite + vue-router + axios |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL（192.168.5.100 docker） |

## 目录结构

```
backend/            Go 后端
  ├── config/       环境变量配置
  ├── models/       用户表 / 故事表 / 回应表 / 硬币流水表
  ├── handlers/     RESTful API（发帖冻结、采纳结算均走数据库事务 + 行锁）
  └── main.go
frontend/           Vue 前端
  └── src/
      ├── views/    Feed.vue（首页）、Profile.vue（个人中心）
      └── components/  Postcard.vue（双面明信片）、StoryCard.vue、AuthModal.vue、PhotoUpload.vue
```

## 启动方式

### 后端（端口 8081）

```bash
cd backend
go mod tidy
go run .
```

可用环境变量（默认值）：

| 变量 | 默认值 |
|------|--------|
| `DB_HOST` | `192.168.5.100` |
| `DB_PORT` | `3306` |
| `DB_USER` | `root` |
| `DB_PASSWORD` | `123456` |
| `DB_NAME` | `memory_link_b` |
| `PORT` | `8081` |

> 注：本机 8080 端口与 `memory_link` 库被分支A实例占用，故分支B默认使用
> 8081 端口与独立的 `memory_link_b` 数据库，两套数据互不影响。

### 前端（端口 5273）

```bash
cd frontend
npm install
npm run dev
```

浏览器打开 <http://localhost:5273>（Vite 已配置 `/api`、`/uploads` 代理到 8081）。

## API 一览

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/register` | 注册（赠送 100 硬币） |
| POST | `/api/login` | 登录，返回 token |
| GET | `/api/stories` | Feed 流（公开） |
| POST | `/api/stories` | 发布需求（事务冻结悬赏） |
| POST | `/api/stories/:id/append` | 追加悬赏（事务追加冻结） |
| POST | `/api/stories/:id/cancel` | 取消需求（事务退回冻结） |
| POST | `/api/stories/:id/respond` | 替他去拍（新照片 + 寄语） |
| POST | `/api/responses/:id/accept` | 确认采纳（事务结算悬赏） |
| GET | `/api/me` | 当前用户信息 |
| GET | `/api/my/stories` | 我的发布 |
| GET | `/api/my/responses` | 我的回应 |
| GET | `/api/my/transactions` | 硬币流水 |
| POST | `/api/upload` | 照片上传 |

账目安全：发帖冻结、追加冻结、取消退款、采纳结算均在 **MySQL 事务**内完成，
并对用户 / 故事行加 `SELECT ... FOR UPDATE` 行锁，防止并发超扣或重复结算。
