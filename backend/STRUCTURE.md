.
├── cmd/
│   ├── api/
│   │   ├── main.go
│   │   ├── config.yaml
│   │   └── wiring.go
│   ├── worker/
│   │   └── main.go
│   └─── migration/
|       ├── schema/
│       │   └── main.go
│       └── data/
│           └── main.go
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   ├── router.go
│   │   └── container.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/
│   │   ├── user/
│   │   │   ├── model/
│   │   │   │   ├── user.go
│   │   │   │   ├── user_id.go
│   │   │   │   ├── email.go
│   │   │   │   ├── phone.go
│   │   │   │   ├── role.go
│   │   │   │   ├── permission.go
│   │   │   │   ├── address.go
│   │   │   │   └── metadata.go
│   │   │   │
│   │   │   ├── valueobject/
│   │   │   │   ├── password_hash.go
│   │   │   │   ├── created_at.go
│   │   │   │   └── updated_at.go
│   │   │   │
│   │   │   ├── aggregate/
│   │   │   │   └── user_aggregate.go
│   │   │   │
│   │   │   ├── port/
│   │   │   │   ├── repository.go
│   │   │   │   ├── cache.go
│   │   │   │   ├── token.go
│   │   │   │   ├── hasher.go
│   │   │   │   ├── otp_sender.go
│   │   │   │   ├── email_sender.go
│   │   │   │   ├── sms_sender.go
│   │   │   │   └── audit_logger.go
│   │   │   │
│   │   │   ├── policy/
│   │   │   │   ├── password_policy.go
│   │   │   │   ├── username_policy.go
│   │   │   │   └── permission_policy.go
│   │   │   │
│   │   │   ├── service/
│   │   │   │   ├── user_domain_service.go
│   │   │   │   └── role_domain_service.go
│   │   │   │
│   │   │   ├── specification/
│   │   │   │   ├── user_filter.go
│   │   │   │   └── user_query_builder.go
│   │   │   │
│   │   │   ├── factory/
│   │   │   │   └── user_factory.go
│   │   │   │
│   │   │   ├── event/
│   │   │   │   ├── user_registered.go
│   │   │   │   ├── user_login.go
│   │   │   │   ├── user_updated.go
│   │   │   │   ├── user_disabled.go
│   │   │   │   └── user_password_changed.go
│   │   │   │
│   │   │   ├── dto/
│   │   │   │   ├── user_dto.go
│   │   │   │   ├── profile_dto.go
│   │   │   │   ├── filter_dto.go
│   │   │   │   └── pagination.go
│   │   │   │
│   │   │   ├── usecase/
│   │   │   │   ├── command/
│   │   │   │   │   ├── register.go
│   │   │   │   │   ├── login.go
│   │   │   │   │   ├── logout.go
│   │   │   │   │   ├── change_password.go
│   │   │   │   │   ├── change_email.go
│   │   │   │   │   ├── update_profile.go
│   │   │   │   │   ├── assign_role.go
│   │   │   │   │   ├── disable_user.go
│   │   │   │   │   └── update_address.go
│   │   │   │   │
│   │   │   │   └── query/
│   │   │   │       ├── get_profile.go
│   │   │   │       ├── get_user.go
│   │   │   │       ├── list_users.go
│   │   │   │       └── search_users.go
│   │   │   │
│   │   │   ├── error/
│   │   │   │   ├── user_errors.go
│   │   │   │   ├── policy_errors.go
│   │   │   │   └── repository_errors.go
│   │   │   │
│   │   │   └── doc.go
│   │   │
│   │   └── ... (other domains such as order/, product/, payment/, auth/, billing/, ...)
│   │
│   ├── infrastructure/
│   │   ├── db/
│   │   │   ├── postgres.go        # Init pgxpool.DB / pgxpool.Pool
│   │   │   ├── transaction.go
│   │   │   ├── migrations/
│   │   │   │   ├── schema/
│   │   │   │   │   ├── 0001_init.sql
│   │   │   │   │   ├── 0002_add_user_profile.sql
│   │   │   │   │   └── 0003_add_role_table.sql
│   │   │   │   │
│   │   │   │   └── data/
│   │   │   │       ├── 0001_seed_roles.jsonl
│   │   │   │       ├── 0002_seed_permissions.jsonl
│   │   │   │       └── 0003_seed_default_admin.jsonl
│   │   │   │
│   │   │   └── sqlc/
│   │   │       ├── query/
│   │   │       │   ├── user/
│   │   │       │   │   ├── user_crud.sql
│   │   │       │   │   ├── user_profile.sql
│   │   │       │   │   ├── user_auth.sql
│   │   │       │   │   ├── user_role.sql
│   │   │       │   │   ├── user_address.sql
│   │   │       │   │   └── user_search.sql
│   │   │       │   │
│   │   │       │   ├── role/
│   │   │       │   │   ├── role_crud.sql
│   │   │       │   │   └── role_permission.sql
│   │   │       │   │
│   │   │       │   └── common/
│   │   │       │       ├── pagination.sql
│   │   │       │       └── audit.sql
│   │   │       │
│   │   │       └── gen/
│   │   │           ├── user/
│   │   │           │   ├── user_crud.sql.go
│   │   │           │   ├── user_profile.sql.go
│   │   │           │   ├── user_auth.sql.go
│   │   │           │   ├── user_role.sql.go
│   │   │           │   ├── user_address.sql.go
│   │   │           │   └── user_search.sql.go
│   │   │           │
│   │   │           ├── role/
│   │   │           │   ├── role_crud.sql.go
│   │   │           │   └── role_permission.sql.go
│   │   │           │
│   │   │           └── common/
│   │   │               ├── pagination.sql.go
│   │   │               └── audit.sql.go
│   │   │
│   │   ├── repository/          
│   │   │   ├── user/
│   │   │   │   ├── user_pg.go       # implements domain/user/port.Repository by sqlc + pgx
│   │   │   │   ├── user_search_pg.go
│   │   │   │   └── user_repository_helpers.go
│   │   │   │
│   │   │   ├── role/
│   │   │   │   └── role_pg.go
│   │   │   │
│   │   │   └── common/
│   │   │       └── repository_helpers.go  # isUniqueViolation, map PG error → domain error
│   │   │
│   │   ├── cache/
│   │   │   ├── redis.go
│   │   │   └── user_cache.go
│   │   │
│   │   ├── mq/
│   │   │   ├── kafka_producer.go
│   │   │   ├── kafka_consumer.go
│   │   │   └── event_dispatcher.go
│   │   │
│   │   ├── email/
│   │   │   └── sendgrid_adapter.go
│   │   │
│   │   ├── sms/
│   │   │   └── twilio_adapter.go
│   │   │
│   │   ├── token/
│   │   │   └── jwt_provider.go
│   │   │
│   │   ├── logger/
│   │   │   └── zap_logger.go
│   │   │
│   │   └── storage/
│   │       └── s3_adapter.go
│   │
│   ├── interface/
│   │   ├── http/
│   │   │   ├── server.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── logger.go
│   │   │   │   └── cors.go
│   │   │   └── handler/
│   │   │       ├── user_handler.go
│   │   │       └── auth_handler.go
│   │   │
│   │   ├── grpc/
│   │   │   ├── user_service.go
│   │   │   └── server.go
│   │   │
│   │   ├── worker/
│   │   │   ├── user_consumer.go
│   │   │   └── email_consumer.go
│   │   │
│   │   └── translator/
│   │       ├── user_converter.go
│   │       └── dto_to_model.go
│   │
│   ├── shared/
│   │   ├── util/
│   │   ├── validator/
│   │   ├── response/
│   │   ├── pagination/
│   │   └── constants/
│   │
│   └── security/
│       ├── rbac.go
│       ├── permission_map.go
│       └── policy.go
│
├── pkg/
│   ├── logger/
│   ├── retry/
│   ├── pagination/
│   └── middleware/
│
├── proto/
│   └── user.proto
│
├── Makefile
├── sqlc.yaml
└── go.mod
