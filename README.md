# Medusa-Inspired Go Modular Monolith Skeleton

A production-ready, high-performance modular monolith e-commerce API built in **Go**, inspired by the architecture of **Medusa 2.0**. Designed for complete domain isolation, clean boundaries, and seamless microservice extraction.

**Tech Stack:**
* **Language/Runtime:** Go 1.22+
* **Routing:** `chi` Router
* **Database & Querying:** PostgreSQL 16+, `sqlx` + `lib/pq`
* **Logging & Telemetry:** `zap` Logger
* **Migrations:** Goose-compatible raw SQL migrations
* **In-Memory & Pub/Sub:** Redis 7+

---

## Architectural Principles

1. **Modular Monolith Layout**: The project contains 13 isolated domain modules under `modules/`. Each module maps 1:1 to its own PostgreSQL database schema.
2. **Domain Isolation**: Modules are strictly decoupled. They can only communicate with other domains through their declared `Service` interfaces. No module is allowed to import another module's repository or query its database schema directly.
3. **Cross-Schema Decoupling**: Rather than physical foreign keys between schemas, this design uses **snapshots** (copying data at creation time), **soft references** (storing UUIDs as plain text fields), and **application-level joins** to keep databases independent. This ensures zero friction if a module needs to be extracted into its own database/microservice.
4. **Read/Write DB Separation (Master/Slave)**: The `db.DB` wrapper splits read and write queries. Scaling reads only requires setting the `DATABASE_REPLICA_URL` environment variable.
5. **Plug-and-Play Factories**: Key infrastructure packages (cache, events, queue, search, monitoring, payment, mail, upload) use the Factory Pattern. You can switch between cloud providers (e.g., MinIO/AWS S3/Cloudflare R2, or Redis/RabbitMQ/Kafka) purely via `.env` configuration without changing any business logic.

---

## Repository Structure

```text
apps/
  api/                    # HTTP API Server
    main.go               # Server entrypoint & dependency wiring
    cmd/migrate/          # Database migration runner entrypoint
    internal/
      server/             # Dependency injection (WireModules) & routing
      middleware/         # Global middlewares (Logger, Recovery, RequestID)
      handler/            # HTTP handlers/controllers (ProductHandler, etc.)
  worker/                 # Background job worker (Skeleton)
modules/                  # 13 Independent Domain Modules
  <module_name>/          # E.g., modules/product
    model.go              # Database model structs with sqlx/json tags
    service.go            # Business logic contract (interface & implementation)
    repository.go         # Handwritten SQL query logic using sqlx
    migrations/           # Module-specific raw SQL migrations
packages/                 # Shared internal utilities & infrastructure
  cache/                  # Cache factory (supports Redis, InMemory with TTL)
  config/                 # Environment configuration loader (.env parser)
  db/                     # DB client with master/slave routing and transaction wrapper
  events/                 # Pub/Sub event bus (supports console log, InMemory, Redis, Kafka, RabbitMQ)
  httpx/                  # Unified HTTP JSON responses and metadata injectors
  logger/                 # zap logger initialization wrapper
  mail/                   # Mailer factory (supports SMTP, Console logger, SendGrid, Mailgun)
  migrate/                # Lightweight goose-compatible migration runner
  monitoring/             # Telemetry factory (supports Console, Prometheus, Datadog, Sentry)
  payment/                # Payment gateway factory (supports Stripe, PayPal, MoMo, VNPay)
  queue/                  # Message Queue factory (supports InMemory, RabbitMQ, Kafka, Redis)
  search/                 # Search engine factory (supports InMemory mock, Meilisearch, Elasticsearch, Algolia)
  types/                  # UUID generation (v7), standard errors, query parameters
  upload/                 # Object storage uploader (S3-compatible: MinIO, Cloudflare R2, AWS S3)
```

---

## Supported Domain Modules

The skeleton contains schemas and stubs for all **13 core Medusa domains** (totaling 55 database tables):
* `auth`
* `identity`
* `customer`
* `product` *(Fully implemented CRUD + variants as a reference pattern)*
* `pricing`
* `inventory`
* `cart`
* `ordering`
* `payment`
* `fulfillment`
* `promotion`
* `region`
* `notification`

---

## Getting Started

### 1. Prerequisites
Make sure you have Go 1.22+, Docker, and Docker Compose installed.

### 2. Run Infrastructure
Start PostgreSQL, Redis, MinIO, and Mailhog using Docker Compose:
```bash
docker compose up -d
```

### 3. Setup Configuration
Copy `.env.example` to `.env` and adjust variables if needed:
```bash
cp .env.example .env
```

### 4. Run Migrations & Start API Server
Run database schema migrations and boot up the server:
```bash
# Apply migrations for all 13 modules
make migrate

# Launch the API server
make run

# Alternatively, run everything with one command (infra up + migrate + run):
make dev
```

The server will be listening on http://localhost:8080.

---

## Testing API Endpoints

You can verify the setup by running these queries:

* **Health Check**:
  ```bash
  curl http://localhost:8080/health
  ```
  *(Responses will contain a standardized metadata block at the end indicating request latency, timestamp, and UUIDv4 Request ID).*

* **Create Product**:
  ```bash
  curl -X POST http://localhost:8080/api/v1/products \
    -H 'Content-Type: application/json' \
    -d '{"title":"Basic T-Shirt","handle":"basic-t-shirt","status":"published"}'
  ```

* **List Products**:
  ```bash
  curl "http://localhost:8080/api/v1/products?per_page=5"
  ```

---

## Development Guidelines

### Adding Endpoints for a New Module
Follow the pattern established in the `product` module:
1. Define the SQL schemas and migrations under `modules/<module>/migrations/`.
2. Map fields to Go structs in `modules/<module>/model.go`.
3. Add business contract methods to the `Service` interface inside `modules/<module>/service.go`.
4. Implement SQL queries inside `modules/<module>/repository.go`.
5. Create a new controller handler under `apps/api/internal/handler/<module>.go`.
6. Mount your new routes in `apps/api/internal/server/server.go`.

### Extracting a Domain into a Microservice
If a module (e.g., `payment`) grows and requires extraction:
1. Move the module folder `modules/payment` to a separate repository.
2. In the monolith, replace `payment.NewService(...)` in `apps/api/internal/server/server.go` with an HTTP or gRPC client implementation that satisfies the exact same `payment.Service` interface.
3. Because the domain modules do not import other repository layers or use direct cross-schema joins, **no handler, middleware, or business logic code in other modules needs to be changed**.

### Telemetry & Metadata Standard
All JSON responses generated through `httpx.JSON` or `httpx.Error` are formatted as:
```json
{
  "your_data_key": "data_value",
  "metadata": {
    "requestId": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2026-06-10T16:45:00Z",
    "latencyUs": 1024
  }
}
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
