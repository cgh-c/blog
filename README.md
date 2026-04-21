# Personal Blog

一个前后端分离的个人博客系统，使用 Go + Vue 3 构建，适合学习全栈 Web 开发。

## 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 后端 | Go + Gin | RESTful API，分层架构 |
| 前端 | Vue 3 + TypeScript | Composition API，Vite 构建 |
| 管理��台 | Vue 3 + TypeScript | 独立 SPA，文章/分类/标签管理 |
| 数据库 | SQLite + GORM | 零配置，WAL 模式 |
| 部署 | Docker + Nginx | 多阶���构建，反向代理 |

## 项目结构

```
blog/
├── backend/                    # Go 后端
│   ├── main.go                 # 入口，加载配置 → 初始化DB → 启动路由
│   ├── internal/
│   │   ├── config/             # 环境变量配置
│   │   ├── database/           # 数据库初始化 + 种子数据
│   │   ├── model/              # 数据模型（Article, Category, Tag, User, Setting）
│   │   ├── repository/         # 数据访问层（SQL 查询）
│   │   ├── service/            # 业务逻辑层
│   │   ├── handler/            # HTTP 处理器（请求 → 响应）
│   │   └── router/             # 路由注册，中间件组装
│   └── pkg/                    # 可复用工具包
│       ├── auth/               # JWT 签发/验证 + bcrypt 密码哈希
│       ├── ratelimit/          # IP 限流
│       ├── response/           # 统一 JSON 响应格式
│       ├── security/           # 安全头 + CORS
│       ├── upload/             # 文件上传
│       └── validator/          # 请求参数校验
│
├── frontend/                   # 博客前端
│   └── src/
│       ├── api/index.ts        # Axios 封装，统一拦截
│       ├── router/index.ts     # 路由定义
│       ├── stores/site.ts      # Pinia 全局状态（设置、暗色模式）
│       ├── styles/main.css     # 全局样式 + CSS 变量
│       ├── components/         # 复用组件
│       │   ├── Header.vue      # 顶部导航
│       │   ├── Footer.vue      # 页脚
│       │   ├── ArticleCard.vue # 文章卡片
│       │   ├── Pagination.vue  # 分页
│       │   ├── TagCloud.vue    # 标签云
│       │   ├── BackToTop.vue   # 回到顶部
│       │   └── TableOfContents.vue  # 文章目录
│       └── views/              # 页面
│           ├── Welcome.vue     # 欢迎页
│           ├── Home.vue        # 文章列表（主页）
│           ├── Article.vue     # 文章详情
│           ├── Category.vue    # 分类筛选
│           ├── Tag.vue         # 标签筛选
│           ├── Archive.vue     # 归档
│           └── About.vue       # 关于
│
├── admin/                      # 管理后台（独立 SPA）
│
├── Dockerfile.backend          # 后端多阶段构建
├── Dockerfile.frontend         # 前端 + Admin 多阶段构建 → Nginx
├── docker-compose.yml          # 容器编排
└── nginx.conf                  # Nginx 反向代理配置
```

## 阅读代码建议

如果你想通过这个项目学���，建议按以下顺序阅读：

### 第一步：理解后端分层架构

从一个请求的完整生命周期入手，以「获取文章列表」为例：

```
请求 GET /api/articles
  → router/router.go          定义路由，绑定 handler
  → handler/public_article.go  解析参数，调用 service
  → service/article.go         业务逻辑，调用 repository
  → repository/article.go      GORM 查询数据库
  → model/article.go           ���据结构定义
  → 响应 JSON
```

**核心文件阅读顺序：**
1. `main.go` — 程序入口，了解启动流程
2. `internal/model/` — 先看数据长什么样
3. `internal/router/router.go` ��� 看有哪些 API
4. `internal/handler/public_article.go` — 跟一个完整请求
5. `internal/service/article.go` — 看业务逻辑
6. `internal/repository/article.go` — 看 GORM 怎么查数据库

### 第二步：理解认证机制

1. `pkg/auth/jwt.go` — JWT Token 的签发和验证
2. `pkg/auth/middleware.go` — 中间件如何拦截未登录请求
3. `handler/auth.go` — 登录/登出逻辑

### 第三步：理解前端数据流

```
用户访问页面
  → router/index.ts        路由匹配组件
  → views/Home.vue         页面组件，onMounted 发起请求
  → api/index.ts           Axios 调后端 API
  → 拿到数据，渲染模板
```

**核心文件阅读顺序：**
1. `main.ts` — Vue 应用初始化
2. `router/index.ts` — 路由表
3. `api/index.ts` — 看拦截器怎么统一处理响应
4. `views/Home.vue` — 看一个完整页面怎么加载数据
5. `stores/site.ts` — Pinia 状态管理
6. `styles/main.css` — CSS 变量和暗色模式实现

### 第四步：理解部署

1. `Dockerfile.backend` — Go 多阶段构建（编译 → 运行）
2. `Dockerfile.frontend` — Node 构建前端 → Nginx 托管
3. `nginx.conf` — 反向代理，SPA fallback
4. `docker-compose.yml` — 服务编排

## 本地开发

```bash
# 1. 启动后端
cd backend
cp ../.env.example .env  # 按需修改
go run .

# 2. 启动前端（另一个终端）
cd frontend
npm install
npm run dev              # http://localhost:5173

# 3. 启动管理后台（另一个终端）
cd admin
npm install
npm run dev              # http://localhost:5174
```

## Docker 部署

```bash
# 创建环境变量
cp .env.example .env
# 编辑 .env 设置 JWT_SECRET、ADMIN_USERNAME、ADMIN_PASSWORD 等

# ��建并启动
docker compose build
docker compose up -d

# 查看日志
docker compose logs -f
```

部署后访问：
- 博客前端：`http://your-server`
- 管理后台：`http://your-server/admin/`

## API 概览

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /api/articles | 文章列表（支持 category_id、tag_id 筛选） | - |
| GET | /api/articles/:slug | 文章详情 | - |
| GET | /api/archives | 归档统计 | - |
| GET | /api/categories | 分类列表 | - |
| GET | /api/tags | 标签列表 | - |
| GET | /api/settings | 站点设置（公开） | - |
| POST | /api/auth/login | 登录 | - |
| POST | /api/auth/logout | 登出 | - |
| GET | /api/admin/articles | 文章管理列表 | JWT |
| POST | /api/admin/articles | 创建文章 | JWT |
| PUT | /api/admin/articles/:id | 更新文章 | JWT |
| DELETE | /api/admin/articles/:id | 删除文章 | JWT |
| POST | /api/admin/upload | 上传文件 | JWT |
| GET | /api/admin/dashboard | 仪表盘统计 | JWT |

## License

MIT
