<!--
Sync Impact Report
- Version change: N/A → 1.0.0
- Modified principles:
  - [PRINCIPLE_1_NAME] → Code Quality & Consistency
  - [PRINCIPLE_2_NAME] → Type Safety
  - [PRINCIPLE_3_NAME] → Error Handling & Observability
  - [PRINCIPLE_4_NAME] → API Design & Separation of Concerns
  - [PRINCIPLE_5_NAME] → Security
- Added sections:
  - Performance & Scalability, Manual Testing
  - Documentation, Specification & UX Discipline
- Removed sections:
  - Template-only placeholders for unnamed sections
- Templates requiring updates:
  - ✅ .specify/templates/plan-template.md
  - ✅ .specify/templates/spec-template.md
  - ✅ .specify/templates/tasks-template.md
  - ✅ .specify/templates/checklist-template.md
  - ✅ .specify/templates/agent-file-template.md
- Follow-up TODOs:
  - TODO(RATIFICATION_DATE): Set the original ratification date for this constitution when known
-->

# English Coach Constitution

## Core Principles

### Code Quality & Consistency

- All code MUST comply with the repository’s unified linting and formatting tooling
  (ESLint/Prettier or language-appropriate equivalents).
- No code with lint errors, failing builds, or unresolved type errors MAY be merged to
  main or released branches.
- Names for variables, functions, methods, classes, and modules MUST be explicit and
  intention-revealing; unclear abbreviations and single-letter names (except trivial
  loop indices) are NOT allowed.
- “Magic values” (e.g., raw status codes, timeouts, configuration knobs) MUST be
  extracted into named constants or configuration with clear semantics.
- Functions, methods, and modules MUST stay small and focused; deeply nested logic
  and long functions MUST be refactored into smaller units.
- The backend and frontend MUST follow clean architecture principles: clear separation
  between domain, application/service, and infrastructure/interface layers, with
  dependencies flowing inward (infrastructure depends on domain, never the reverse).

**Rationale**: Consistent, readable, and well-structured code reduces defects, speeds up
code reviews, and makes the system maintainable as the project and team grow.

### Type Safety

- Both backend and frontend MUST use a static type system (TypeScript for JS/TS-based
  code or a language-native equivalent such as Go types).
- The use of `any` (or equivalent untyped escape hatches) is strongly discouraged and
  MUST be treated as an exception:
  - Each `any` MUST include a comment explaining why it is necessary and, where
    applicable, a TODO indicating how to remove it later.
- Data models and API contracts MUST have explicit types and schemas and SHOULD be
  shared across frontend/backend when technically feasible (e.g., shared TypeScript
  types, OpenAPI-generated clients, protobufs).
- When changing any data model or API contract, the corresponding types, schemas, and
  generated artifacts MUST be updated first, before or together with implementation
  changes.

**Rationale**: Strong typing shifts many classes of bugs from runtime to compile time
and keeps contracts between services and components explicit and reliable.

### Error Handling & Observability

- All APIs MUST return errors in a unified schema, for example:
  - `{ code: string; message: string; details?: unknown }`
  - Implementations MAY extend this schema, but the core shape MUST stay consistent.
- There MUST be a shared error-handling layer (middleware/interceptor) in each API
  surface; uncaught exceptions MUST be captured, normalized into the unified error
  schema, and logged.
- Logging MUST capture sufficient context for debugging:
  - Request identifiers or correlation IDs (when available),
  - Key input parameters (redacted for sensitive data),
  - Error codes and stack traces for failures.
- Critical flows (e.g., authentication, payments, and core data operations such as
  creating/updating primary entities) MUST produce structured logs at key steps
  (start, important decisions, success, failure).
- When monitoring/observability tooling is available, important APIs MUST emit basic
  metrics (e.g., request count, error count, latency buckets) and/or traces to
  support troubleshooting and performance analysis.

**Rationale**: Unified error handling and observability turn production issues from
guesswork into diagnosable, actionable problems.

### API Design & Separation of Concerns

- All APIs MUST be designed against a clear spec (OpenAPI or equivalent REST contract)
  prior to or in parallel with implementation.
- Controllers/route handlers MAY only perform orchestration:
  - Validate/parse input (or delegate to a validation layer),
  - Call application/service/use case layer,
  - Map domain results to HTTP/transport responses.
- Business rules and domain logic MUST live in dedicated service or use-case layers
  and MUST NOT be embedded in controllers, database adapters, or UI code.
- Data access (databases, caches, message queues, external APIs) MUST reside in
  infrastructure/repository/adapters; domain layers MUST depend only on interfaces,
  not concrete infrastructure implementations.
- APIs MUST be consistent in:
  - Naming (resource-oriented URIs, verbs aligned with actions),
  - Status codes (2xx for success, 4xx for client errors, 5xx for server errors),
  - Error schema and pagination/filters patterns.

**Rationale**: Clear separation of concerns and consistent API design simplify change,
testing, and the ability to evolve the system without breaking consumers.

### Security

- All communication that carries sensitive data (authentication tokens, credentials,
  personal data, payment details) MUST be protected with TLS (HTTPS or equivalent).
- Authentication MUST use a unified mechanism across the system (e.g., JWT,
  session-based auth, or OAuth2) with centralized configuration and documentation.
- Authorization MUST follow Role-Based Access Control (RBAC):
  - Roles and permissions MUST be modeled explicitly,
  - Permission checks MUST NOT be hard-coded ad hoc throughout the codebase; they
    MUST be routed through a small number of well-defined authorization components.
- All external inputs (body, query, params, headers) MUST be validated before
  reaching business logic, using Zod or an equivalent validation library.
- The system MUST implement reasonable protection against common vulnerabilities:
  - XSS (encoding, CSP where appropriate),
  - SQL injection (parameterized queries/ORM),
  - CSRF where applicable (especially for browser-based sessions),
  - Rate limiting for authentication and other sensitive endpoints.
- Secrets (tokens, API keys, database passwords, encryption keys) MUST NOT be
  committed to the repository; they MUST be supplied via environment variables or a
  secrets manager appropriate to the deployment environment.

**Rationale**: By treating security as a first-class concern and centralizing
authentication, authorization, and validation, the system reduces exposure to
critical vulnerabilities and compliance risks.

## Performance & Scalability, Manual Testing

- The frontend MUST implement basic performance hygiene:
  - Use lazy loading or code splitting for non-essential routes/features,
  - Avoid loading unnecessary bundles or large libraries when not needed.
- Static assets SHOULD leverage browser and server-side caching, with cache headers
  configured appropriately for their expected update cadence.
- Backend services SHOULD be designed as stateless where practical, so that
  horizontal scaling (adding more instances) is straightforward.
- APIs returning large datasets MUST use pagination, cursor-based navigation, or
  streaming; endpoints MUST NOT return unbounded result sets.
- Database access patterns MUST be reviewed for indexes on primary query paths and to
  avoid N+1 query patterns.
- Automated tests are OPTIONAL; the project relies primarily on disciplined manual
  testing:
  - Before releasing any feature, core happy paths and basic error scenarios MUST be
    manually tested.
  - Critical flows (login, registration, payments, CRUD on core entities) MUST be
    manually smoke-tested before each release.
  - When bugs are found, reproduction steps and expected behavior MUST be documented
    (e.g., in an issue tracker or project board) for shared visibility.
- Features with known blocking bugs in primary flows MUST NOT be released.

**Rationale**: Thoughtful performance and scalability practices prevent bottlenecks,
while explicit manual testing discipline compensates for the absence of mandatory
automated tests.

## Documentation, Specification & UX Discipline

- Every new feature MUST start from a written specification:
  - Problem statement and goals,
  - Expected behavior and user journeys,
  - Key business rules and constraints.
- Any significant change to data models, APIs, or business rules MUST be reflected in
  the appropriate docs/specs (e.g., contracts, data-model docs, README/quickstart).
- README and onboarding materials MUST remain sufficient for a new contributor to:
  - Understand the overall architecture (frontend, backend, infrastructure),
  - Set up the project locally,
  - Run the main services and perform basic manual tests.
- Important architectural or product decisions (including rejected alternatives)
  SHOULD be captured with:
  - Context/problem,
  - Options considered,
  - Chosen approach and rationale.
- Frontend implementations MUST respect a shared design system (colors, typography,
  spacing, and standard components) and prioritize reuse of components over one-off
  variants.
- All user-critical interactions MUST provide clear feedback:
  - Loading states for long operations,
  - Error states with actionable messages,
  - Success confirmations.
- Basic accessibility MUST be upheld:
  - Text contrast sufficient for readability,
  - Clickable areas and controls large enough for comfortable use,
  - Forms with explicit labels and visible focus states.

**Rationale**: High-quality documentation and consistent UX make the system easier
to extend, reduce onboarding cost, and improve end-user satisfaction.

## Governance

- This constitution defines non-negotiable engineering and product quality standards
  for the English Coach project and supersedes informal or undocumented practices.
- All feature work MUST reference this constitution during:
  - Specification (`spec.md`),
  - Planning (`plan.md`, including “Constitution Check”),
  - Task breakdown (`tasks.md`),
  - Checklists (`checklist` artifacts for releases or features).
- Code reviews and merge approvals MUST explicitly verify:
  - Linting/build/type checks are clean,
  - Type safety rules are respected (no unexplained `any` or untyped contracts),
  - Error handling, logging, and validation follow shared patterns,
  - Security, performance, and UX requirements relevant to the change are addressed,
  - Specs and docs have been updated when models/APIs/business rules change.
- Constitution versions follow semantic versioning:
  - MAJOR: Backward-incompatible governance changes or removal/redefinition of
    principles.
  - MINOR: New principles/sections or material expansion of guidance.
  - PATCH: Clarifications, wordings, or non-semantic refinements.
- Amendments MUST:
  - Be proposed in writing (e.g., PR modifying this file),
  - Include an updated Sync Impact Report at the top of this file,
  - Update related templates and docs when the change affects workflows,
  - Set `LAST AMENDED` to the date the amendment is accepted.
- Compliance SHOULD be periodically reviewed (at least once per quarter or major
  release) to ensure that the actual practices in the codebase stay aligned with this
  constitution; gaps SHOULD be captured as tasks and prioritized.

**Version**: 1.0.0 | **Ratified**: TODO(RATIFICATION_DATE): Set initial ratification date | **Last Amended**: 2025-11-29
