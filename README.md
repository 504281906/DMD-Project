# DMD-Project (Device MDM Platform)

开源移动设备管理平台，专注于中小企业市场，提供设备管理、应用分发、安全策略等核心功能。

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 📌 项目简介

DMD-Project 是一款面向中小企业的开源 MDM（Mobile Device Management）解决方案，采用 Go + React 技术栈，支持本地化部署，数据自主可控。

### 核心特性

- **设备管理** - 设备注册、状态监控、分组管理
- **应用分发** - 企业应用商店、批量推送、安装统计
- **安全策略** - 密码策略、加密要求、合规检测、远程锁定/擦除

---

## ✅ 已实现功能

### 后端 API (Go + Gin + GORM)

| 模块 | 功能 | 状态 |
|------|------|------|
| **设备管理** | 设备 CRUD、远程锁定、远程擦除 | ✅ |
| **应用分发** | 应用 CRUD、设备安装应用、安装统计 | ✅ |
| **安全策略** | 策略 CRUD、分配到设备、规则验证 | ✅ |
| **数据库** | SQLite 本地存储（生产环境可切换 PostgreSQL） | ✅ |

### 前端管理后台 (React + Ant Design)

| 页面 | 功能 | 状态 |
|------|------|------|
| **仪表盘** | 统计卡片、设备分布图表、最近设备列表 | ✅ |
| **设备列表** | 搜索/分页、锁定/擦除/详情 | ✅ |
| **应用列表** | 筛选、上传、编辑、删除 | ✅ |
| **策略列表** | 创建、编辑、删除、分配到设备 | ✅ |

### API 路由 (21个)

```
设备管理:
  POST   /api/v1/devices              # 创建设备
  GET    /api/v1/devices              # 列出设备
  GET    /api/v1/devices/:id          # 获取设备详情
  PUT    /api/v1/devices/:id          # 更新设备
  DELETE /api/v1/devices/:id          # 删除设备
  POST   /api/v1/devices/:id/lock     # 远程锁定
  POST   /api/v1/devices/:id/wipe     # 远程擦除

应用管理:
  POST   /api/v1/apps                 # 创建应用
  GET    /api/v1/apps                 # 列出应用
  GET    /api/v1/apps/:id             # 获取应用详情
  PUT    /api/v1/apps/:id             # 更新应用
  DELETE /api/v1/apps/:id             # 删除应用

设备应用:
  POST   /api/v1/devices/:id/apps     # 设备安装应用
  GET    /api/v1/devices/:id/apps     # 列出设备已安装应用

策略管理:
  POST   /api/v1/policies             # 创建策略
  GET    /api/v1/policies             # 列出策略
  GET    /api/v1/policies/stats       # 策略统计
  GET    /api/v1/policies/:id         # 获取策略详情
  PUT    /api/v1/policies/:id         # 更新策略
  DELETE /api/v1/policies/:id         # 删除策略

设备策略:
  POST   /api/v1/devices/:id/policies              # 分配策略到设备
  DELETE /api/v1/devices/:id/policies/:policyId    # 移除设备策略
  GET    /api/v1/devices/:id/policies             # 列出设备关联策略
  POST   /api/v1/devices/:id/policies/:policyId/apply  # 应用策略
```

---

## 🚧 待完成需求

### P0 (必须)

- [ ] **iOS MDM 协议对接** - 实现 Apple MDM 协议，支持 iOS 设备管理
- **Android EMM 对接** - 实现 Android Enterprise Device Admin
- **客户端 SDK** - iOS/Android Agent 应用开发

### P1 (重要)

- [ ] **用户认证** - JWT 认证、角色权限管理
- [ ] **多租户** - 支持多个企业/组织隔离
- [ ] **消息推送** - APNs (iOS) / FCM (Android) 集成
- [ ] **设备分组** - 动态分组、标签管理

### P2 (可选)

- [ ] **MAC/PC 支持** - 扩展到桌面端设备管理
- [ ] **SaaS 部署** - Docker 一键部署方案
- [ ] **API 开放** - 第三方系统集成
- [ ] **审计日志** - 操作记录、合规报表

---

## 🚀 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+ (前端开发)
- SQLite / PostgreSQL

### 1. 克隆代码

```bash
git clone https://github.com/504281906/DMD-Project.git
cd DMD-Project
```

### 2. 启动后端

```bash
# 安装依赖
go mod tidy

# 编译运行
go build -o openmdm-api ./cmd/api/...
./openmdm-api

# 服务启动后访问 http://localhost:8080
```

### 3. 启动前端

```bash
cd web
npm install
npm run dev

# 访问 http://localhost:3000
```

### 4. 生产构建

```bash
# 前端构建
cd web
npm run build

# 后端构建
cd ..
go build -o openmdm-api ./cmd/api/...
```

### 5. Docker 部署 (TODO)

```bash
docker-compose up -d
```

---

## 📁 项目结构

```
DMD-Project/
├── cmd/
│   └── api/
│       └── main.go           # 主入口
├── internal/
│   ├── model/                # 数据模型
│   │   ├── device.go        # 设备模型
│   │   ├── application.go    # 应用模型
│   │   └── policy.go        # 策略模型
│   ├── repository/           # 数据访问层
│   ├── service/              # 业务逻辑层
│   ├── handler/              # API 处理器
│   └── middleware/          # 中间件
├── web/                     # React 前端
│   ├── src/
│   │   ├── pages/           # 页面组件
│   │   ├── components/      # 通用组件
│   │   └── services/       # API 调用
│   └── dist/               # 构建产物
├── config.yaml             # 配置文件
├── openmdm.db              # SQLite 数据库
├── go.mod / go.sum         # Go 依赖
└── README.md
```

---

## ⚙️ 配置说明

配置文件 `config.yaml`:

```yaml
server:
  host: 0.0.0.0
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: openmdm
  sslmode: disable
  # 切换为 sqlite:
  # driver: sqlite
  # dsn: openmdm.db

rabbitmq:
  host: localhost
  port: 5672
  user: guest
  password: guest

app:
  name: OpenMDM
  version: 0.1.0
  debug: true
```

---

## 🔧 技术栈

### 后端

- **语言**: Go 1.22+
- **框架**: Gin Web Framework
- **ORM**: GORM
- **数据库**: SQLite (开发) / PostgreSQL (生产)

### 前端

- **框架**: React 18
- **构建**: Vite
- **UI 库**: Ant Design 5.x
- **语言**: TypeScript

---

## 📄 开源协议

本项目基于 **MIT 协议** 开源。

```
MIT License

Copyright (c) 2026 DMD-Project

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

---

## 📬 联系方式

- **GitHub Issues**: https://github.com/504281906/DMD-Project/issues

---

<p align="center">
  <strong>DMD-Project</strong> - 开源 MDM解决方案
</p>
