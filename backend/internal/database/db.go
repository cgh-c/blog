package database

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"blog-backend/internal/model"
	"blog-backend/pkg/auth"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init opens the SQLite database, runs auto-migration, and returns the DB instance.
func Init(dbPath string, isProduction bool) *gorm.DB {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	logLevel := logger.Info
	if isProduction {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Enable WAL mode for better concurrent read performance
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		&model.Setting{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// SeedAdmin creates the initial admin user if no users exist.
func SeedAdmin(db *gorm.DB, username, password string) {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	admin := model.User{
		Username:     username,
		PasswordHash: hash,
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Printf("Admin user '%s' created successfully", username)
}

// SeedSettings creates default site settings if none exist.
func SeedSettings(db *gorm.DB) {
	var count int64
	db.Model(&model.Setting{}).Count(&count)
	if count > 0 {
		return
	}

	defaults := []model.Setting{
		{Key: "site_name", Value: "My Blog", Type: "string", IsPublic: true, Description: "Blog name displayed in navbar"},
		{Key: "site_description", Value: "A personal blog", Type: "string", IsPublic: true, Description: "Blog description for SEO and footer"},
		{Key: "about_content", Value: "# About Me\n\nWelcome to my blog!", Type: "text", IsPublic: true, Description: "About page content in Markdown"},
		{Key: "avatar_url", Value: "", Type: "string", IsPublic: true, Description: "Avatar URL for about page"},
		{Key: "social_links", Value: "[]", Type: "json", IsPublic: true, Description: "Social media links as JSON array"},
		{Key: "icp_number", Value: "", Type: "string", IsPublic: true, Description: "ICP filing number for footer"},
		{Key: "admin_email", Value: "", Type: "string", IsPublic: false, Description: "Admin email (internal use only)"},
	}

	db.Create(&defaults)
	log.Println("Default settings seeded")
}

// SeedSampleData creates sample categories, tags, and articles for demonstration.
// It skips any data that already exists (by slug), so it's safe to run on a DB with existing content.
func SeedSampleData(db *gorm.DB) {
	// Categories — insert only if slug doesn't exist
	categoryDefs := []model.Category{
		{Name: "技术笔记", Slug: "tech-notes", Description: "编程学习与技术分享", SortOrder: 1},
		{Name: "项目实战", Slug: "projects", Description: "实际项目开发经验总结", SortOrder: 2},
		{Name: "随笔思考", Slug: "thoughts", Description: "生活感悟与思考记录", SortOrder: 3},
	}
	var categories []model.Category
	for _, c := range categoryDefs {
		var existing model.Category
		if db.Where("slug = ?", c.Slug).First(&existing).Error == nil {
			categories = append(categories, existing)
		} else {
			db.Create(&c)
			categories = append(categories, c)
		}
	}

	// Tags — insert only if slug doesn't exist
	tagDefs := []model.Tag{
		{Name: "Go", Slug: "go"},
		{Name: "Vue", Slug: "vue"},
		{Name: "Docker", Slug: "docker"},
		{Name: "数据库", Slug: "database"},
		{Name: "前端", Slug: "frontend"},
		{Name: "后端", Slug: "backend"},
		{Name: "入门教程", Slug: "tutorial"},
	}
	var tags []model.Tag
	for _, t := range tagDefs {
		var existing model.Tag
		if db.Where("slug = ?", t.Slug).First(&existing).Error == nil {
			tags = append(tags, existing)
		} else {
			db.Create(&t)
			tags = append(tags, t)
		}
	}

	now := time.Now()

	articles := []struct {
		article model.Article
		tagIDs  []uint
	}{
		{
			article: model.Article{
				Title:      "Go 语言入门：从零开始构建 Web 服务",
				Slug:       "go-web-getting-started",
				CategoryID: &categories[0].ID,
				Visibility: model.VisibilityPublic,
				ViewCount:  128,
				PublishedAt: &now,
				Summary:    "本文介绍如何使用 Go 语言和 Gin 框架从零搭建一个 Web 服务，涵盖路由、中间件、数据库连接等核心概念。",
				Content: `# Go 语言入门：从零开始构建 Web 服务

## 为什么选择 Go？

Go 语言（Golang）由 Google 开发，具有以下优势：

- **编译速度快**：几秒内完成编译
- **并发支持**：goroutine 和 channel 让并发编程变得简单
- **部署简单**：编译为单个二进制文件，无需运行时依赖
- **性能优秀**：接近 C 语言的性能

## 第一个 Web 服务

使用 Gin 框架创建一个简单的 HTTP 服务：

` + "```go" + `
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()

    r.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, World!",
        })
    })

    r.Run(":8080")
}
` + "```" + `

## 项目结构

一个良好的 Go 项目应该遵循清晰的目录结构：

` + "```" + `
├── cmd/            # 应用入口
├── internal/       # 项目内部代码（不可被外部引用）
│   ├── handler/    # HTTP 处理器
│   ├── service/    # 业务逻辑
│   ├── model/      # 数据模型
│   └── repository/ # 数据访问层
├── pkg/            # 可复用的公共包
├── go.mod          # Go 模块定义
└── main.go         # 入口文件
` + "```" + `

## 中间件

Gin 的中间件机制非常强大，可以用于日志记录、认证、跨域处理等：

` + "```go" + `
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "未授权"})
            return
        }
        c.Next()
    }
}
` + "```" + `

## 总结

Go 语言的简洁语法和强大的标准库使其成为构建 Web 服务的理想选择。配合 Gin 框架，可以快速开发高性能的 RESTful API。

> 下一篇我们将探讨如何连接数据库和实现 CRUD 操作。`,
			},
			tagIDs: []uint{tags[0].ID, tags[5].ID, tags[6].ID},
		},
		{
			article: model.Article{
				Title:      "Vue 3 组合式 API 实战指南",
				Slug:       "vue3-composition-api-guide",
				CategoryID: &categories[0].ID,
				Visibility: model.VisibilityPublic,
				ViewCount:  96,
				PublishedAt: &now,
				Summary:    "深入理解 Vue 3 的 Composition API，通过实际案例学习 ref、reactive、computed、watch 等核心概念。",
				Content: `# Vue 3 组合式 API 实战指南

## 从 Options API 到 Composition API

Vue 3 引入了组合式 API（Composition API），它解决了 Options API 在大型组件中代码组织混乱的问题。

## 核心概念

### ref 和 reactive

` + "```typescript" + `
import { ref, reactive } from 'vue'

// ref 用于基本类型
const count = ref(0)
console.log(count.value) // 0

// reactive 用于对象
const state = reactive({
  name: '张三',
  age: 25
})
console.log(state.name) // 张三
` + "```" + `

### computed 计算属性

` + "```typescript" + `
import { ref, computed } from 'vue'

const firstName = ref('张')
const lastName = ref('三')

const fullName = computed(() => {
  return firstName.value + lastName.value
})
` + "```" + `

### watch 侦听器

` + "```typescript" + `
import { ref, watch } from 'vue'

const keyword = ref('')

watch(keyword, (newVal, oldVal) => {
  console.log('搜索词变化：', oldVal, '->', newVal)
  // 执行搜索请求...
})
` + "```" + `

## 组件通信

### Props 和 Emits

` + "```vue" + `
<script setup lang="ts">
// 定义 props
const props = defineProps<{
  title: string
  count?: number
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update', value: number): void
}>()

function handleClick() {
  emit('update', props.count ? props.count + 1 : 1)
}
</script>
` + "```" + `

### Pinia 状态管理

对于跨组件共享状态，推荐使用 Pinia：

` + "```typescript" + `
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    name: '',
    isLoggedIn: false
  }),
  actions: {
    login(name: string) {
      this.name = name
      this.isLoggedIn = true
    }
  }
})
` + "```" + `

## 总结

Composition API 让我们能够更灵活地组织代码逻辑，特别适合复杂组件的开发。结合 TypeScript，可以获得更好的类型安全和开发体验。`,
			},
			tagIDs: []uint{tags[1].ID, tags[4].ID, tags[6].ID},
		},
		{
			article: model.Article{
				Title:      "Docker 部署实战：从开发到上线",
				Slug:       "docker-deployment-practice",
				CategoryID: &categories[1].ID,
				Visibility: model.VisibilityPublic,
				ViewCount:  75,
				PublishedAt: &now,
				Summary:    "记录使用 Docker 和 Docker Compose 将全栈项目部署到云服务器的完整过程，包含踩坑经验。",
				Content: `# Docker 部署实战：从开发到上线

## 为什么用 Docker？

传统部署方式需要在服务器上安装各种运行环境，不同项目的依赖可能冲突。Docker 通过容器化解决了这个问题：

- **环境一致性**：开发环境 = 生产环境
- **隔离性**：每个服务独立运行，互不影响
- **可移植**：一次构建，到处运行

## Dockerfile 编写

### 后端（Go）多阶段构建

` + "```dockerfile" + `
# 第一阶段：编译
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o server .

# 第二阶段：运行（镜像仅 ~20MB）
FROM alpine:3.19
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]
` + "```" + `

多阶段构建的好处是最终镜像非常小，不包含编译工具链。

### 前端（Nginx）

` + "```dockerfile" + `
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
` + "```" + `

## Docker Compose 编排

` + "```yaml" + `
version: '3.8'
services:
  backend:
    build:
      context: .
      dockerfile: Dockerfile.backend
    ports:
      - "8080:8080"
    volumes:
      - blog-data:/app/data

  nginx:
    build:
      context: .
      dockerfile: Dockerfile.frontend
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  blog-data:
` + "```" + `

## 常用命令

| 命令 | 说明 |
|------|------|
| ` + "`docker compose up -d`" + ` | 后台启动所有服务 |
| ` + "`docker compose logs -f`" + ` | 查看实时日志 |
| ` + "`docker compose down`" + ` | 停止并移除容器 |
| ` + "`docker compose build --no-cache`" + ` | 强制重新构建 |

## 踩坑记录

### 1. 国内镜像拉取超时

需要配置 Docker 镜像加速源，在 ` + "`/etc/docker/daemon.json`" + ` 中添加国内镜像地址。

### 2. SQLite 需要 CGO

Go 编译 SQLite 驱动需要开启 CGO，Alpine 镜像需要安装 ` + "`gcc`" + ` 和 ` + "`musl-dev`" + `。

### 3. 数据持久化

一定要使用 Docker Volume 挂载数据目录，否则容器重建后数据会丢失！

## 总结

Docker 大幅简化了部署流程，虽然初次配置需要一些学习成本，但一旦配好，后续的更新部署只需要 ` + "`git pull && docker compose up -d --build`" + ` 一条命令。`,
			},
			tagIDs: []uint{tags[2].ID, tags[0].ID, tags[5].ID},
		},
		{
			article: model.Article{
				Title:      "SQLite + GORM：轻量级数据库方案",
				Slug:       "sqlite-gorm-guide",
				CategoryID: &categories[0].ID,
				Visibility: model.VisibilityPublic,
				ViewCount:  52,
				PublishedAt: &now,
				Summary:    "介绍如何在 Go 项目中使用 GORM 操作 SQLite 数据库，包括模型定义、CRUD 操作和查询优化。",
				Content: `# SQLite + GORM：轻量级数据库方案

## 为什么选 SQLite？

对于个人项目和小型应用，SQLite 是一个绝佳选择：

- **零配置**：无需安装数据库服务，数据存储在单个文件中
- **性能优秀**：对于读多写少的场景，性能不输 MySQL
- **易于备份**：复制一个文件即可完成备份
- **资源占用低**：非常适合小型云服务器

## GORM 基础

### 模型定义

` + "```go" + `
type Article struct {
    ID         uint       ` + "`gorm:\"primaryKey\" json:\"id\"`" + `
    Title      string     ` + "`gorm:\"size:200;not null\" json:\"title\"`" + `
    Content    string     ` + "`gorm:\"type:text\" json:\"content\"`" + `
    CategoryID *uint      ` + "`gorm:\"index\" json:\"category_id\"`" + `
    Category   *Category  ` + "`gorm:\"foreignKey:CategoryID\" json:\"category\"`" + `
    Tags       []Tag      ` + "`gorm:\"many2many:article_tags\" json:\"tags\"`" + `
    CreatedAt  time.Time  ` + "`json:\"created_at\"`" + `
}
` + "```" + `

### 自动迁移

` + "```go" + `
db.AutoMigrate(&Article{}, &Category{}, &Tag{})
` + "```" + `

GORM 会自动创建表和索引，开发阶段非常方便。

### CRUD 操作

` + "```go" + `
// 创建
db.Create(&Article{Title: "Hello", Content: "World"})

// 查询
var article Article
db.Preload("Category").Preload("Tags").First(&article, 1)

// 更新
db.Model(&article).Updates(map[string]interface{}{
    "title": "New Title",
})

// 删除
db.Delete(&article)
` + "```" + `

### 分页查询

` + "```go" + `
func ListArticles(db *gorm.DB, page, size int) ([]Article, int64) {
    var articles []Article
    var total int64

    db.Model(&Article{}).Count(&total)
    db.Offset((page - 1) * size).
       Limit(size).
       Order("created_at DESC").
       Find(&articles)

    return articles, total
}
` + "```" + `

## SQLite 优化建议

### 开启 WAL 模式

` + "```go" + `
db.Exec("PRAGMA journal_mode=WAL")
` + "```" + `

WAL（Write-Ahead Logging）模式允许读写并发，显著提升性能。

### 开启外键约束

` + "```go" + `
db.Exec("PRAGMA foreign_keys=ON")
` + "```" + `

SQLite 默认不启用外键约束，需要手动开启。

## 总结

SQLite + GORM 的组合非常适合个人项目，简单高效。当项目规模增长后，GORM 的数据库驱动可以轻松切换到 MySQL 或 PostgreSQL，业务代码几乎不需要修改。`,
			},
			tagIDs: []uint{tags[0].ID, tags[3].ID, tags[6].ID},
		},
		{
			article: model.Article{
				Title:      "一个程序员的独立博客之路",
				Slug:       "why-build-personal-blog",
				CategoryID: &categories[2].ID,
				Visibility: model.VisibilityPublic,
				ViewCount:  200,
				PublishedAt: &now,
				Summary:    "为什么要搭建个人博客？记录自己从零开始构建博客系统的心路历程和收获。",
				Content: `# 一个程序员的独立博客之路

## 为什么要写博客？

在这个信息爆炸的时代，为什么还要费力搭建一个个人博客？

### 1. 学习的最好方式是输出

费曼学习法告诉我们：**如果你不能简单地解释一件事，说明你还没有真正理解它。** 写技术博客就是一种强迫自己深入理解的过程。

### 2. 建立个人品牌

在技术社区中，一个持续输出高质量内容的博客，比任何简历都更有说服力。

### 3. 记录成长轨迹

回顾几年前写的文章，你会清晰地看到自己的成长。这种感觉无可替代。

## 为什么要自己造轮子？

市面上有很多成熟的博客方案：WordPress、Hexo、Hugo……为什么还要自己写？

**因为这本身就是最好的学习项目。**

一个博客系统虽然看起来简单，但涵盖了 Web 开发的方方面面：

- 前后端分离架构设计
- 用户认证和权限控制
- 数据库设计和 ORM 使用
- RESTful API 设计
- Docker 容器化部署
- Nginx 反向代理配置

每一个环节都是实打实的技能点。

## 技术选型的思考

| 技术 | 选择 | 原因 |
|------|------|------|
| 后端语言 | Go | 想学习新语言，编译部署简单 |
| Web 框架 | Gin | Go 生态最流行，文档丰富 |
| 数据库 | SQLite | 个人项目够用，部署零配置 |
| 前端框架 | Vue 3 | 学习曲线友好，生态成熟 |
| 部署方式 | Docker | 环境一致，一键部署 |

## 写在最后

> "种一棵树最好的时间是十年前，其次是现在。"

不管你用什么技术栈，不管你写的文章有没有人看，**开始写**才是最重要的。

这个博客就是我的第一步。希望多年后回看，能为当初的决定感到骄傲。`,
			},
			tagIDs: []uint{},
		},
	}

	created := 0
	for _, item := range articles {
		var existing model.Article
		if db.Where("slug = ?", item.article.Slug).First(&existing).Error == nil {
			continue // already exists, skip
		}
		db.Create(&item.article)
		if len(item.tagIDs) > 0 {
			var articleTags []model.Tag
			db.Where("id IN ?", item.tagIDs).Find(&articleTags)
			db.Model(&item.article).Association("Tags").Replace(articleTags)
		}
		created++
	}

	if created > 0 {
		log.Printf("Sample data seeded: %d new articles added", created)
	} else {
		log.Println("Sample data: all articles already exist, nothing to add")
	}
}
