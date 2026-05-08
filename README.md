# Mall — 在线课程商城后端系统

基于 **Go + Gin** 的单体 Web 服务，为课程商城提供完整的后台管理与 C 端接口能力。

## 技术栈

| 组件 | 技术选型 |
|------|---------|
| HTTP 框架 | Gin |
| ORM | GORM + gorm/gen (模型代码生成) |
| 数据库 | MySQL 8.0 |
| 缓存 / Token 存储 | Redis 7 |
| 配置中心 | 本地 YAML / etcd 远程配置（支持热更新） |
| 对象存储 | 腾讯云 COS |
| OAuth 登录 | 飞书 (Lark)、微信公众号、微信小程序 |
| 支付 | 微信支付 (API V3) |
| 日志 | go.uber.org/zap |
| 协程池 | ants |
| ID 生成 | Snowflake |
| 容器化 | Docker + Docker Compose |

## 功能特性

### 管理后台 (`/api/mall/admin`)

- **管理员认证** — 手机号 + 短信验证码 / 密码登录；飞书 OAuth 扫码登录
- **管理员管理** — 增删改查、角色权限分配 (RBAC)
- **权限管理** — 菜单权限树、角色绑定
- **课程管理** — 课程 CRUD、课程目录编排、录播课时管理、上下架控制
- **订单管理** — 订单列表、详情、统计、退款处理
- **C 端用户管理** — 用户列表查询、状态变更（启用/禁用）

### C 端 (`/api/mall/customer`)

- **用户认证** — 手机号 + 短信验证码 / 密码登录；微信公众号扫码登录；微信小程序登录
- **个人中心** — 用户信息、密码修改、微信绑定/解绑
- **课程浏览** — 课程列表、详情、课时信息
- **学习记录** — 课时学习进度上报与查询、继续学习、已购课程
- **购物车** — 添加/移除/查看购物车商品
- **订单** — 费用计算、立即支付/稍后支付、取消订单、订单列表与详情
- **微信回调** — 支付成功回调、退款回调

## 项目结构

```
mall/
├── main.go                  # 入口：初始化 MySQL / Redis / 配置，启动服务
├── Makefile                 # gendb 生成数据库模型
├── Dockerfile               # 多阶段构建 (golang:1.26-alpine → alpine:3.21)
├── docker-compose.yml       # 本地开发环境 (MySQL + Redis + etcd + 后端)
├── mall_local.yml.example   # 本地配置文件模板
│
├── api/                     # HTTP 处理器层 (Gin handlers)
│   ├── admin/               #   管理后台接口
│   ├── customer/            #   C 端接口
│   └── resp.go              #   统一 JSON 响应封装
│
├── service/                 # 业务逻辑层 (无状态单例)
│   ├── admin/               #   管理员业务
│   ├── user/                #   用户认证（管理端 & C 端）
│   ├── goods/               #   课程商品业务
│   ├── perm/                #   权限业务
│   ├── role/                #   角色业务
│   ├── storage/             #   对象存储业务
│   ├── token/               #   Token 管理
│   ├── do/                  #   领域对象 (映射数据库表)
│   └── dto/                 #   请求/响应 DTO
│
├── adaptor/                 # 依赖注入 & 数据访问层
│   ├── adaptor.go           #   IAdaptor 接口 (DI 容器)
│   ├── repo/                #   数据库仓库 (gorm/gen 生成)
│   │   ├── query/           #     查询代码
│   │   ├── model/           #     数据模型
│   │   └── gen.yaml         #     gorm/gen 配置文件
│   ├── redis/               #   Redis 操作封装
│   └── rpc/                 #   外部 API 客户端 (飞书/Lark、腾讯云 COS)
│
├── router/                  # 路由 & 中间件
│   ├── router.go            #   路由注册入口
│   ├── auth.go              #   Token 鉴权中间件（管理端 & C 端）
│   ├── white_list.go        #   白名单路由（无需鉴权）
│   ├── access.go            #   访问日志中间件
│   └── pprof.go             #   pprof 性能分析
│
├── common/                  # 通用类型 (Errno、User、Pager)
├── consts/                  # 枚举常量 (订单状态、验证码场景、Token TTL)
├── config/                  # 配置加载 (本地 YAML / etcd 远程)
├── utils/                   # 工具库 (日志、验证码、Snowflake ID、协程池、加密)
├── web/                     # 静态前端页面 (登录页等)
└── data/                    # Docker 挂载数据目录 (gitignored)
```

## 架构分层

```
api/  (Gin handlers — 解析请求 → 调用 service → 返回 JSON)
  │
service/  (业务逻辑 — 编排 repo/redis/rpc 调用)
  │
adaptor/
  ├── repo/    (数据访问 — gorm/gen 生成的查询代码)
  ├── redis/   (缓存 & Token — 验证码、令牌、分布式锁、二维码状态)
  └── rpc/     (外部 API — 飞书、腾讯云 COS)
```

所有依赖通过 `adaptor.IAdaptor` 接口注入。Service 和 Repo 构造函数接收 `IAdaptor`，从中获取 `*gorm.DB`、`*redis.Client`、`*config.Config`。

## 快速开始

### 1. 环境准备

- Go 1.26+
- Docker & Docker Compose
- 配置 `.env` 文件和 `mall_local.yml`（参考 `.example` 模板）

### 2. 本地开发

```bash
# 启动依赖服务 (MySQL + Redis + etcd)
docker-compose up -d mysql redis etcd

# 编译 & 运行
go build -o mall.backend main.go
./mall.backend -c mall_local.yml

# 或直接开发运行
go run main.go -c mall_local.yml
```

### 3. Docker 完整部署

```bash
# 启动全部服务 (包括后端应用)
docker-compose up -d

# 查看服务状态
docker-compose ps
```

### 4. 生成数据库模型

```bash
make gendb
```

## 配置说明

| 配置项 | 说明 |
|--------|------|
| `server.http_port` | HTTP 服务端口 |
| `server.log_level` | 日志级别 (debug / info / warn / error) |
| `server.enable_pprof` | 是否开启 pprof 性能分析 |
| `mysql.*` | MySQL 连接配置 (支持 `${ENV_VAR}` 环境变量展开) |
| `redis.*` | Redis 连接配置 |
| `security.*` | 手机号 AES 加密密钥、密码哈希盐值、Token TTL |
| `order.*` | 订单支付超时、自动确认收货天数、Snowflake 节点 ID |
| `wechat_pay.*` | 微信支付 API V3 商户配置 |
| `storage.*` | 腾讯云 COS 对象存储配置 |
| `app_conf.*` | 第三方应用配置 (微信公众号、飞书) |

支持两种配置加载方式：

- **本地模式** (`-c mall_local.yml`) — 从 YAML 文件加载，支持 `${VAR}` 环境变量展开
- **远程模式** (`-r http://etcd:2379`) — 从 etcd 读取 `/configs/mall/system`，每 1 分钟自动热更新

## 统一错误响应

所有接口返回统一格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

- `code` 为 `0` 表示成功，非 `0` 为业务错误码
- `msg` 为可读的错误描述
- `data` 为业务数据（错误时可能为 `null`）

## 后台任务

服务启动后自动开启定时任务，每分钟执行：

- 取消超时未支付的订单
- 自动确认已收货超期的订单
