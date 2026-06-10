# ecommerce-go

Skeleton e-commerce backend lấy cảm hứng từ kiến trúc Medusa 2.0.
Single repo, modular monolith, sẵn sàng tách microservice.

**Stack:** Go 1.22 · chi router · sqlx + lib/pq · zap · goose (SQL migrations) · PostgreSQL 16

---

## Kiến trúc

- **Modular monolith:** 13 domain module trong `modules/`, mỗi module map 1:1 với 1 PostgreSQL schema.
- **Module isolation:** module chỉ giao tiếp qua `Service` interface, không import repository của nhau → tách microservice sau chỉ cần đổi implementation interface.
- **Cross-schema:** dùng SNAPSHOT (copy data lúc create), SOFT REF (lưu ID không FK), application-level join. Không cross-schema FK.
- **Master/slave ready:** `db.DB` tách Write/Read connection; ban đầu trỏ chung DSN, scale chỉ cần set `DATABASE_REPLICA_URL`.

## Cấu trúc

```
apps/
  api/              # HTTP API server (chi)
    main.go
    cmd/migrate/    # migration runner
    internal/
      server/       # wire modules + mount routes
      middleware/   # logger, recovery (zap)
      handler/      # HTTP handlers
  worker/           # background jobs (skeleton)
modules/            # 13 domain modules
  <module>/
    model.go        # structs với db: tags
    service.go      # business logic (interface + impl)
    repository.go   # SQL viết tay (sqlx)
    migrations/     # goose .sql files
packages/
  db/               # sqlx connection, tx, master/slave
  events/           # event bus (LogBus; thay RedisBus khi prod)
  logger/           # zap
  httpx/            # JSON response + error mapping
  migrate/          # goose-compatible runner
  types/            # ID gen, errors, pagination
```

## 13 schemas

`auth` `identity` `customer` `product` `pricing` `inventory` `cart` `ordering` `payment` `fulfillment` `promotion` `region` `notification` — tổng 55 tables.

`product` module được implement đầy đủ (CRUD + variants) làm pattern mẫu. Các module khác có model đầy đủ + service/repository stub.

---

## Chạy

```bash
# 1. Khởi động infra
docker compose up -d

# 2. Chạy migrations (quét modules/*/migrations)
make migrate

# 3. Chạy API
make run
# hoặc gộp cả 3: make dev
```

Test:

```bash
curl localhost:8080/health

curl -X POST localhost:8080/v1/products \
  -H 'Content-Type: application/json' \
  -d '{"title":"Áo thun basic","handle":"ao-thun-basic","status":"published"}'

curl "localhost:8080/v1/products?per_page=5"
```

---

## Mở rộng

**Thêm endpoint cho module khác:** làm theo `modules/product` —
1. Thêm method vào `Service` interface + impl trong `service.go`
2. Viết SQL trong `repository.go`
3. Tạo `apps/api/internal/handler/<module>.go`
4. Mount route trong `server.New()`

**Tách microservice:** đổi wire trong `server.WireModules` — thay `product.NewService(...)` bằng HTTP/gRPC client implement cùng `product.Service` interface. Code handler không đổi.

**Event bus production:** thay `events.NewLogBus` bằng implementation Redis Streams / NATS, giữ nguyên `EventBus` interface.

---

## Lưu ý

- Migration runner trong `packages/migrate` là bản nhẹ tương thích format goose. File `.sql` vẫn chạy được bằng goose CLI thật: `goose -dir modules/<m>/migrations postgres "$DATABASE_URL" up`.
- Tiền tệ lưu bằng số nguyên (cents / smallest unit), tránh float.
