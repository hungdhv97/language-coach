<!--
Sync Impact Report

- Version change: N/A (template placeholders) → 1.0.0
- Modified principles:
  - Template placeholders → I. Monorepo & Module Boundaries (STRUCTURE.md)
  - Template placeholders → II. API Contract & Documentation
  - Template placeholders → III. Quality & Safe Changes
  - Template placeholders → IV. Frontend (Vite + shadcn/ui)
  - Template placeholders → V. Backend (Go REST/JSON)
- Added sections:
  - Testing, CI/CD, Tooling & Secrets
  - Project Documentation
- Templates requiring updates:
  - ✅ .specify/templates/plan-template.md
  - ✅ .specify/templates/spec-template.md
  - ✅ .specify/templates/tasks-template.md
  - ⚠ N/A: .specify/templates/commands/*.md (folder does not exist in this repo)
-->

# LexiGo Constitution

## Core Principles

### I. Monorepo & Module Boundaries (STRUCTURE.md)

- **MANDATORY**: All file/folder changes MUST comply with the `STRUCTURE.md` of the relevant root directory:
  - `frontend/STRUCTURE.md`, `backend/STRUCTURE.md`, `deploy/STRUCTURE.md`, `scripts/STRUCTURE.md`,
    `ops/STRUCTURE.md`, `.github/STRUCTURE.md`
- **No arbitrary large restructuring**: Do not change large layouts / move files in bulk without clear
  justification + migration plan + `STRUCTURE.md` update.
- **Clear boundaries**:
  - `frontend/` and `backend/` are separate apps; communicate **only** via HTTP API (OpenAPI).
  - Shared code **within each app** resides in that app's shared area (FE: `frontend/src/shared/`,
    BE: `backend/internal/shared/` or `backend/pkg/` when public API is needed).
  - Avoid circular dependencies; prioritize minimal dependencies; new dependencies must have specific justification.

### II. API Contract & Documentation (OpenAPI / Postman)

- **All API changes must have clear contracts**: request/response, status code, headers, error codes,
  example payloads.
- **Source of truth**:
  - OpenAPI: `backend/docs/openapi/**` (create/change/delete endpoints must update immediately).
  - If feature-specific checklist/spec exists: update `specs/<feature>/contracts/` synchronously when applying.
- **Backward-compatible by default**: breaking changes must use versioning (e.g., `/api/v2`) and have
  migration/deprecation plans.
- **Standardized error response**:
  - Body: `{"code": "SOME_CODE", "message": "…", "details": {...}}` (details optional)
  - Correlation: requests **always** have `X-Request-ID` (received from client or generated) and server **always**
    returns `X-Request-ID` header. Contract docs must describe this header for every endpoint.

### III. Quality & Safe Changes

- **Priority**: Correctness > Clarity > Performance. Optimize only when there is data or clear bottlenecks.
- **No "magic"**: clear naming, short functions, readable flow; avoid over-engineering.
- **Small, reviewable changes**: split PRs by slice, avoid bundling multiple topics; prefer small refactors
  alongside behavior changes within the same scope.

### IV. Frontend (Vite + TS + shadcn/ui)

- **shadcn/ui**: always prioritize usage according to the latest version; when adding new components, use
  official shadcn tooling to sync correct API and current style (reference `frontend/components.json`).
- **Consistent UX**:
  - use shared layouts/components before creating new ones
  - have clear and consistent loading/empty/error states
- **Basic accessibility**: reasonable labels/aria, focus/keyboard navigation for main flows.
- **Sufficient performance**: avoid unnecessary re-renders; lazy-load routes/chunks when appropriate; avoid heavy dependencies.

### V. Backend (Go REST/JSON API)

- **Go convention**: `gofmt` mandatory; clear package/module structure; do not overuse interfaces (only use when
  there is benefit for replacing implementations or clear boundaries).
- **Input & context**: validate input; propagate `context.Context`; set reasonable timeouts/deadlines for DB
  and outbound calls.
- **Error & logging**:
  - Usecases return standardized errors (e.g., `AppError`); adapters/middleware map to HTTP + unified error shape
    (code/message/details) and do not leak internal causes.
  - Structured logging; always include `request_id` (from `X-Request-ID`) in log context.

## Testing, CI/CD, Tooling & Secrets

- **Testing (MANDATORY)**: Do not write **unit test** code.
  - Instead: maintain/update manual test checklists (e.g., `specs/<feature>/manual-test-checklist.md`)
    and describe "Independent Verification" in spec.
  - If automated tests are needed for special requirements: only consider contract/integration/e2e (not unit tests)
    and must clearly state the reason in plan.
- **Tooling**:
  - Frontend: lint/format according to current config (eslint, ts).
  - Backend: gofmt + linter (golangci-lint if available).
  - Prefer using repo scripts (e.g., `scripts/lint-all.sh`) to lint FE/BE synchronously.
- **CI/CD**: no need to run tests; focus on lint/format + necessary builds.
- **Secrets & environment**:
  - Do not commit secrets.
  - When adding/modifying/deleting environment variables: immediately update corresponding env example files:
    `deploy/env/example/backend.env` and/or `deploy/env/example/frontend.env`.

## Project Documentation

- **No need to write/maintain README as mandatory requirement**; operational documentation for features resides in
  `specs/<feature>/` (spec/plan/contracts/quickstart/manual-test-checklist).
- **Docs must sync with code**: changes to API/contract/env/structure must update corresponding docs
  in the same PR.

## Governance

- **This constitution is supreme**: conflicts between other docs and this file → this file wins.
- **Amendments**:
  - All amendments must update: content + **Version** (SemVer) + **Last Amended** + Sync Impact Report
  - Sync templates in `.specify/templates/` when principles change
- **SemVer for constitution**:
  - MAJOR: principle changes that break current practices (or remove principles)
  - MINOR: add principles/sections or significantly expand guidance
  - PATCH: clarify wording/typos, no meaning change
- **Compliance expectation**: all new spec/plan/tasks must have "Constitution Check" section and pass
  related gates (structure, contract, error shape, env example, lint/format).

**Version**: 1.0.0 | **Ratified**: 2026-01-08 | **Last Amended**: 2026-01-08
