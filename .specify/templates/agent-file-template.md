# [PROJECT NAME] Development Guidelines

Auto-generated from all feature plans. Last updated: [DATE]

> This document summarizes active, technology-specific guidance. It MUST be read
> together with the project constitution defined in `.specify/memory/constitution.md`.

## Active Technologies

[EXTRACTED FROM ALL PLAN.MD FILES]

## Project Structure

```text
[ACTUAL STRUCTURE FROM PLANS]
```

## Commands

[ONLY COMMANDS FOR ACTIVE TECHNOLOGIES]

## Code Style & Quality Gates

- Follow language- and framework-specific conventions summarized here.
- Ensure code passes all configured linting, formatting, and type-checking tools
  before merging.
- Respect the architectural boundaries (domain, application/service,
  infrastructure/interface) described in the constitution and feature plans.

## Error Handling, Security & Observability

- Use the shared error-handling and logging infrastructure for all new features.
- Apply the security expectations from the constitution (auth, RBAC, validation,
  secrets, and basic protections against XSS/SQL injection/CSRF).
- For critical flows, prefer structured logs and, where available, basic metrics
  or tracing hooks.

## Testing & UX

- Favor clear, reproducible manual test scenarios derived from user stories; record
  smoke test coverage for critical flows before each release.
- Keep UX consistent with the design system and ensure user feedback (loading,
  error, success states) is implemented for important interactions.

## Recent Changes

[LAST 3 FEATURES AND WHAT THEY ADDED]

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
