.
├── cmd/
│   └── api/
│       └── main.go                                  # Application entrypoint: load config, build DI graph, start servers, handle graceful shutdown
├── configs/
│   ├── config.development.yaml                      # Development environment configuration (local development)
│   ├── config.staging.yaml                          # Staging environment configuration (pre-production testing)
│   ├── config.production.yaml                       # Production environment configuration (production deployment)
│   └── config.go                                    # Configuration loading/parsing/validation (env-aware: reads config.{env}.yaml based on APP_ENV)
├── docs/
│   ├── openapi/
│   │   ├── openapi.yaml                             # Root OpenAPI spec (info/servers + references to paths/components)
│   │   ├── components/
│   │   │   ├── schemas.yaml                         # Reusable schema definitions (request/response DTOs)
│   │   │   ├── responses.yaml                       # Reusable response templates (standard errors, common responses)
│   │   │   └── parameters.yaml                      # Reusable parameter definitions (pagination, headers, etc.)
│   │   └── paths/
│   │       ├── users.yaml                           # User endpoints (paths/operations)
│   │       └── auth.yaml                            # Auth endpoints (paths/operations)
│   └── postman/
│       ├── collections/
│       │   └── myapp.postman_collection.json        # Postman collection for manual QA / development
│       └── environments/
│           ├── local.postman_environment.json       # Postman environment for local testing
│           └── staging.postman_environment.json     # Postman environment for staging testing
├── db/
│   ├── migrations/
│   │   ├── 0001_init.sql                            # Migration: initial schema setup
│   │   └── 0002_add_indexes.sql                     # Migration: add indexes/constraints for performance/integrity
│   ├── queries/
│   │   ├── user/
│   │   │   ├── user.sql                             # sqlc queries for user domain (CRUD, list, search)
│   │   │   └── user_admin.sql                       # sqlc queries for admin-grade user operations
│   │   ├── auth/
│   │   │   └── auth.sql                             # sqlc queries for authentication/session/token flows
│   │   └── shared/
│   │       └── outbox.sql                           # sqlc queries for outbox pattern (optional event publishing)
│   └── schema/
│       └── schema.sql                               # Schema snapshot/reference (not a replacement for migrations)
├── sqlc.yaml                                        # sqlc configuration (engine, schema/queries inputs, code generation output)
├── scripts/
│   ├── migrate.sh                                   # Run database migrations (up/down/status) for local/CI usage
│   ├── sqlc_generate.sh                             # Generate sqlc code from db/queries and db/schema
│   └── lint.sh                                      # Lint/format/static analysis runner for CI/local
├── internal/
│   ├── app/
│   │   ├── bootstrap/
│   │   │   ├── wire.go                              # Dependency wiring/composition root (wire or manual constructors)
│   │   │   ├── http_server.go                       # HTTP server construction/startup (router + middleware chain)
│   │   │   ├── grpc_server.go                       # gRPC server construction/startup (optional)
│   │   │   └── cron.go                              # Scheduler/cron initialization (optional background jobs)
│   │   ├── di/
│   │   │   └── container.go                         # DI container/provider definitions (interfaces -> implementations)
│   │   └── lifecycle/
│   │       └── shutdown.go                          # Graceful shutdown orchestration (signals, context cancellation, resource cleanup)
│   ├── shared/
│   │   ├── logger/
│   │   │   ├── logger.go                            # Logger interface + contextual field helpers
│   │   │   └── zap_logger.go                        # zap-based implementation of the logger interface
│   │   ├── errors/
│   │   │   ├── app_error.go                         # Canonical application error type (code/message/cause/metadata)
│   │   │   ├── codes.go                             # Error code catalog (domain/usecase/http mapping anchors)
│   │   │   └── http_mapper.go                       # Maps AppError -> HTTP status + standardized error response body
│   │   ├── validation/
│   │   │   └── validator.go                         # Validation wrapper (struct tags + custom rules)
│   │   ├── auth/
│   │   │   ├── jwt.go                               # JWT signing/verification + claim parsing utilities
│   │   │   ├── password.go                          # Password hashing/verification utilities (bcrypt/argon2 wrapper)
│   │   │   └── permissions.go                       # RBAC primitives (roles/permissions helpers)
│   │   ├── security/
│   │   │   ├── headers.go                           # Security header configuration utilities (HSTS, nosniff, etc.)
│   │   │   ├── cors.go                              # CORS policy configuration utilities (allowlist-based)
│   │   │   └── csrf.go                              # CSRF utilities (only for cookie-based auth flows)
│   │   ├── observability/
│   │   │   ├── request_context.go                   # Context keys (request_id/user_id/trace_id) + helpers
│   │   │   ├── metrics.go                           # Metrics primitives/exporter wiring (Prometheus/OTel)
│   │   │   └── tracing.go                           # Tracing primitives (OTel tracer/span helpers)
│   │   └── pagination/
│   │       └── pagination.go                        # Pagination utilities (parse page/pageSize or limit/offset, calculate metadata)
│   ├── platform/
│   │   ├── db/
│   │   │   ├── postgres.go                          # PostgreSQL connection pool initialization (pgx) + health checks
│   │   │   ├── tx.go                                # Transaction helpers (begin/commit/rollback) with context propagation
│   │   │   └── sqlc/
│   │   │       ├── store.go                         # Store wrapper to expose sqlc Queries and transaction-scoped execution
│   │   │       └── gen/                             # sqlc-generated code (DO NOT EDIT)
│   │   │           ├── db.go                        # sqlc core types: Queries struct + constructors
│   │   │           ├── models.go                    # sqlc model structs/types generated from schema
│   │   │           ├── user.sql.go                  # Generated query methods for db/queries/user/*.sql
│   │   │           ├── auth.sql.go                  # Generated query methods for db/queries/auth/*.sql
│   │   │           └── outbox.sql.go                # Generated query methods for db/queries/shared/outbox.sql
│   │   ├── cache/
│   │   │   └── redis.go                             # Redis client initialization + cache helpers (namespacing, TTL)
│   │   ├── mq/
│   │   │   ├── kafka.go                             # Kafka producer/consumer initialization (optional)
│   │   │   └── rabbitmq.go                          # RabbitMQ connection/channel initialization (optional)
│   │   └── httpclient/
│   │       └── client.go                            # Outbound HTTP client (timeouts/retries/circuit-breaking hooks)
│   ├── transport/
│   │   ├── http/
│   │   │   ├── server.go                            # net/http server configuration (timeouts, base settings)
│   │   │   ├── router.go                            # Route registration (module routers) + middleware assembly
│   │   │   ├── response.go                          # Standard response envelope (data/error/meta) helpers
│   │   │   └── middleware/
│   │   │       ├── request_id.go                    # Assign/propagate request ID (X-Request-ID) into context/logs
│   │   │       ├── real_ip.go                       # Resolve real client IP using trusted proxy headers
│   │   │       ├── recover.go                       # Catch panics and convert to standardized 500 errors
│   │   │       ├── error_handler.go                 # Centralized error rendering (AppError -> HTTP JSON)
│   │   │       ├── access_log.go                    # Access logging (method/path/status/latency/bytes/request_id)
│   │   │       ├── audit_log.go                     # Audit logging for sensitive operations (who did what/when)
│   │   │       ├── authn_jwt.go                     # JWT authentication middleware (principal injection into context)
│   │   │       ├── authz_rbac.go                    # Authorization middleware (RBAC/permission checks)
│   │   │       ├── cors.go                          # CORS enforcement middleware
│   │   │       ├── security_headers.go              # Apply security headers (HSTS, nosniff, frame options, etc.)
│   │   │       ├── csrf.go                          # CSRF protection middleware (cookie-based auth only)
│   │   │       ├── rate_limit.go                    # Rate limiting middleware (by IP/user/token with burst/window)
│   │   │       ├── timeout.go                       # Request timeout middleware (context cancellation)
│   │   │       ├── body_limit.go                    # Request body size limit middleware (DoS protection)
│   │   │       ├── gzip.go                          # Response compression middleware (optional)
│   │   │       ├── etag.go                          # ETag/If-None-Match support for cacheable GET endpoints
│   │   │       └── idempotency.go                   # Idempotency key support for safe retries on write endpoints
│   │   └── grpc/
│   │       ├── server.go                            # gRPC server setup + service registration (optional)
│   │       └── interceptors/
│   │           ├── recover.go                       # gRPC panic recovery interceptor
│   │           ├── error_mapper.go                  # Map AppError/domain errors to gRPC status codes/details
│   │           ├── access_log.go                    # gRPC access logging interceptor (method/latency/code)
│   │           ├── authn.go                         # gRPC authentication interceptor (JWT via metadata)
│   │           └── timeout.go                       # gRPC deadline enforcement interceptor
│   └── modules/
│       ├── user/
│       │   ├── domain/
│       │   │   ├── entity.go                        # User entity + core invariants
│       │   │   ├── value_objects.go                 # Domain value objects (UserID/Email/etc.) with validation
│       │   │   ├── repository.go                    # Domain repository interface (usecases depend on this)
│       │   │   └── errors.go                        # Domain-specific errors (not found, invalid state, etc.)
│       │   ├── usecase/
│       │   │   ├── create_user/
│       │   │   │   ├── handler.go                   # Use case orchestration for creating a user (kept thin)
│       │   │   │   ├── input.go                     # Use case input DTO
│       │   │   │   ├── output.go                    # Use case output DTO
│       │   │   │   └── validator.go                 # Use case input validation rules
│       │   │   ├── get_user/
│       │   │   │   ├── handler.go                   # Use case orchestration for fetching a user
│       │   │   │   ├── input.go                     # Use case input DTO (ID)
│       │   │   │   └── output.go                    # Use case output DTO (view model)
│       │   │   └── list_users/
│       │   │       ├── handler.go                   # Use case orchestration for listing users
│       │   │       ├── input.go                     # Use case input DTO (pagination/filter)
│       │   │       ├── output.go                    # Use case output DTO (items + meta)
│       │   │       └── filter.go                    # Filter/query composition helpers to keep handler small
│       │   ├── adapter/
│       │   │   └── http/
│       │   │       ├── routes.go                    # User route registration + per-route middleware binding
│       │   │       ├── handler_create.go            # HTTP handler: parse/validate -> call create_user use case
│       │   │       ├── handler_get.go               # HTTP handler: parse -> call get_user use case
│       │   │       ├── handler_list.go              # HTTP handler: parse -> call list_users use case
│       │   │       └── dto.go                       # HTTP DTOs (requests/responses), separate from domain/usecase
│       │   └── infra/
│       │       └── persistence/
│       │           └── postgres/
│       │               ├── repo.go                  # Repository implementation wrapping sqlc Store/Queries
│       │               └── mapper.go                # Mapping between sqlc models and domain entities
│       ├── auth/
│       │   ├── domain/
│       │   │   ├── entity.go                        # Auth-related entities (Session/RefreshToken) if needed
│       │   │   ├── repository.go                    # Auth repository interface (tokens/sessions persistence)
│       │   │   └── errors.go                        # Auth domain errors (invalid credentials, revoked token, etc.)
│       │   ├── usecase/
│       │   │   ├── login/
│       │   │   │   ├── handler.go                   # Login orchestration (verify credentials, issue tokens)
│       │   │   │   ├── input.go                     # Login input DTO
│       │   │   │   ├── output.go                    # Login output DTO (access/refresh tokens)
│       │   │   │   └── validator.go                 # Login validation rules
│       │   │   └── refresh_token/
│       │   │       ├── handler.go                   # Refresh orchestration (validate refresh token, rotate/issue)
│       │   │       ├── input.go                     # Refresh input DTO
│       │   │       └── output.go                    # Refresh output DTO
│       │   ├── adapter/
│       │   │   └── http/
│       │   │       ├── routes.go                    # Auth route registration
│       │   │       ├── handler_login.go             # HTTP handler for login
│       │   │       ├── handler_refresh.go           # HTTP handler for refresh
│       │   │       └── dto.go                       # HTTP DTOs for auth endpoints
│       │   └── infra/
│       │       └── persistence/
│       │           └── postgres/
│       │               ├── repo.go                  # Auth repository implementation wrapping sqlc
│       │               └── mapper.go                # Mapping between sqlc models and auth domain types
│       └── health/
│           ├── usecase/
│           │   └── ping/
│           │       ├── handler.go                   # Health/ping orchestration (optionally checks dependencies)
│           │       └── output.go                    # Health response DTO (status/version/build info)
│           └── adapter/
│               └── http/
│                   ├── routes.go                    # Health route registration (e.g., /healthz, /readyz)
│                   └── handler_ping.go              # HTTP handler for health checks
├── pkg/
│   ├── id/
│   │   └── uuid.go                                  # Public UUID helpers (generate/parse/validate)
│   └── timeutil/
│       └── timeutil.go                              # Public time utilities (UTC normalization, parsing)
├── test/                                             # OPTIONAL: test scaffolding (include only if you want a dedicated top-level test directory)
│   ├── integration/                                  # OPTIONAL: integration tests (DB + HTTP server)
│   └── contract/                                     # OPTIONAL: contract tests (OpenAPI/Postman/consumer-driven)
├── go.mod                                            # Go module definition and dependency constraints
└── go.sum                                            # Dependency checksums (auto-generated)
