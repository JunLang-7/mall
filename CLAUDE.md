# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
# Build
go build -o mall.backend main.go

# Run locally (uses mall_local.yml config, with ${ENV_VAR} expansion)
./mall.backend -c mall_local.yml

# Run with remote etcd config
ETCD_ADDR=http://etcd:2379 ./mall.backend -r http://etcd:2379

# Docker Compose (starts etcd, mysql, redis, backend)
docker compose up -d
```

```bash
# Generate DB models with gorm/gen (reads adaptor/repo/gen.yaml)
make gendb
```

There are no tests yet.

## Architecture

This is a Go monolith web server for a course mall platform, using **Gin** for HTTP, **GORM** (with gorm/gen for model generation) for MySQL, and **go-redis** for caching and token storage. DingTalk/Lark OAuth is used for admin login; WeChat OAuth for customer login.

### Dependency Injection via Adaptor

`adaptor.IAdaptor` is the central DI container. It wraps `*config.Config`, `*gorm.DB`, and `*redis.Client`. Constructors for services and repos accept `IAdaptor` and pull what they need from it. See `adaptor/adaptor.go:9-13`.

### Layered flow

```
api/          — Gin handlers; parse request, call service, write JSON response via common.Errno
service/      — Business logic; stateless singletons. service/user handles both admin & customer auth
  do/         — Domain objects (structs mapping to DB tables)
  dto/        — Request/response DTOs
adaptor/
  repo/       — Data access layer; query/ folder holds gorm/gen generated query code; model/ holds generated models; separate packages per domain (admin/, goods/, customer/)
  redis/      — Redis helpers (verify codes, tokens, locks, QR code state, order fee locks)
  rpc/        — External API clients (Lark/DingTalk, Tencent COS)
router/
  router.go   — Route definitions, admin and customer groups with auth middleware
  auth.go     — Token-based auth middleware (admin and customer variants)
  access.go   — Access log middleware with request/response body capture
  white_list.go — Unauthenticated routes (login, captcha, WeChat callbacks)
common/       — Shared types: Errno (error code + message pattern), User/AdminUser, Pager
consts/       — Enums: order statuses, SMS code scene constants, token TTLs
config/       — YAML config loading (local file with env var expansion, or remote etcd with hot reload)
utils/        — Logger, captcha, snowflake ID generator, goroutine pool, encryption helpers
web/          — Static HTML pages (login, upload, etc.)
```

### Auth model

Two separate auth domains:
- **Admin** (`/api/mall/admin/...`) — token in `token` header, stored in `adaptor/repo/model/admin_user.gen.go`. Lark OAuth, mobile + SMS code/password.
- **Customer** (`/api/mall/customer/...`) — token in `token` header, stored via `adaptor/repo/model/app_user.gen.go` / `wechat_user.gen.go` / `mobile_user.gen.go`. WeChat OAuth, mobile + SMS code/password.

Auth middleware (`router/auth.go`) checks the white list first, then validates the token header. On success, the user object is stored in `gin.Context` via `ctx.Set()` and retrieved by handlers through `api.GetUserFromCtx()` / `api.GetAdminUserFromCtx()`.

### Error handling pattern

All errors use `common.Errno` (code + message + optional detail). Service methods return `common.Errno` by value. Handlers call `api.WriteResp(ctx, data, errno)` to render JSON. Never use Go's native `error` as the public API response — wrap underlying errors with `.WithErr(err)` on the Errno.

### Config

- Local: `mall_local.yml` with `${ENV_VAR}` expansion (see `.env.example` for required vars). The `-c` flag specifies the path.
- Remote: etcd key `/configs/mall/system`. The `-r` flag specifies the etcd address (or set `ETCD_ADDR` env var). Config hot-reloads every minute.
- Sensitive files (`mall_local.yml`, `.env`) are gitignored; `.example` files are committed.

### Background tasks

`service/user/tasks.go` starts a goroutine on boot that runs every minute: cancel timed-out unpaid orders, auto-confirm received orders after `AutoReceiveDays`.
