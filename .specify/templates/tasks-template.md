---

description: "Task list template for feature implementation"
---

# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Testing Policy (repo)**: Do not write **unit test** code. Instead, each user story must have clear
manual verification steps (checklist/test steps). If automated tests are needed for special requirements:
only consider contract/integration/e2e (not unit tests) and must clearly state exception in plan.

**Organization**: Tasks are grouped by user story to enable independent implementation and verification of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Repo layout**: `backend/`, `frontend/`, `deploy/`, `scripts/`, `specs/`
- **MANDATORY**: comply with `STRUCTURE.md` in each root directory (especially `backend/STRUCTURE.md` and `frontend/STRUCTURE.md`)
- Tasks must specify actual file paths (Go/TS) instead of generic examples

<!-- 
  ============================================================================
  IMPORTANT: The tasks below are SAMPLE TASKS for illustration purposes only.
  
  The /speckit.tasks command MUST replace these with actual tasks based on:
  - User stories from spec.md (with their priorities P1, P2, P3...)
  - Feature requirements from plan.md
  - Entities from data-model.md
  - Endpoints from contracts/
  
  Tasks MUST be organized by user story so each story can be:
  - Implemented independently
  - Tested independently
  - Delivered as an MVP increment
  
  DO NOT keep these sample tasks in the generated tasks.md file.
  ============================================================================
-->

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per implementation plan
- [ ] T002 Initialize [language] project with [framework] dependencies
- [ ] T003 [P] Configure linting and formatting tools

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

Examples of foundational tasks (adjust based on your project):

- [ ] T004 Setup database schema and migrations framework
- [ ] T005 [P] Implement authentication/authorization framework
- [ ] T006 [P] Setup API routing and middleware structure
- [ ] T007 Create base models/entities that all stories depend on
- [ ] T008 Configure error handling and logging infrastructure
- [ ] T009 Setup environment configuration management

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - [Title] (Priority: P1) 🎯 MVP

**Goal**: [Brief description of what this story delivers]

**Independent Verification**: [How to verify this story works on its own]

### Verification for User Story 1 (MANDATORY) ✅

- [ ] T010 [P] [US1] Update manual test checklist in `specs/[###-feature-name]/manual-test-checklist.md`
- [ ] T011 [US1] Manual verify: run FE/BE, execute main flow, capture results/errors (if any)

### Implementation for User Story 1

- [ ] T012 [P] [US1] If API exists: update OpenAPI in `backend/docs/openapi/**` (request/response/status/error)
- [ ] T013 [US1] Backend: implement usecase/handler per `backend/STRUCTURE.md` (validate input, context/timeout)
- [ ] T014 [P] [US1] Frontend: implement UI/flow per `frontend/STRUCTURE.md` (loading/empty/error + basic a11y)
- [ ] T015 [US1] Standardize errors: return `{code,message,details?}` + `X-Request-ID` header
- [ ] T016 [US1] Structured logging with `request_id`

**Checkpoint**: At this point, User Story 1 should be fully functional and verifiable independently

---

## Phase 4: User Story 2 - [Title] (Priority: P2)

**Goal**: [Brief description of what this story delivers]

**Independent Verification**: [How to verify this story works on its own]

### Verification for User Story 2 (MANDATORY) ✅

- [ ] T018 [P] [US2] Update manual test checklist in `specs/[###-feature-name]/manual-test-checklist.md`
- [ ] T019 [US2] Manual verify: execute US2 user journey and ensure it doesn't break US1

### Implementation for User Story 2

- [ ] T020 [P] [US2] If API exists: update OpenAPI in `backend/docs/openapi/**`
- [ ] T021 [US2] Backend: implement per `backend/STRUCTURE.md`
- [ ] T022 [P] [US2] Frontend: implement per `frontend/STRUCTURE.md`
- [ ] T023 [US2] Integrate with US1 (if needed) but ensure US2 can be verified independently

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - [Title] (Priority: P3)

**Goal**: [Brief description of what this story delivers]

**Independent Verification**: [How to verify this story works on its own]

### Verification for User Story 3 (MANDATORY) ✅

- [ ] T024 [P] [US3] Update manual test checklist in `specs/[###-feature-name]/manual-test-checklist.md`
- [ ] T025 [US3] Manual verify: execute US3 user journey and ensure it doesn't break US1/US2

### Implementation for User Story 3

- [ ] T026 [P] [US3] If API exists: update OpenAPI in `backend/docs/openapi/**`
- [ ] T027 [US3] Backend: implement per `backend/STRUCTURE.md`
- [ ] T028 [P] [US3] Frontend: implement per `frontend/STRUCTURE.md`

**Checkpoint**: All user stories should now be independently functional

---

[Add more user story phases as needed, following the same pattern]

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] TXXX [P] Documentation updates in docs/
- [ ] TXXX Code cleanup and refactoring
- [ ] TXXX Performance optimization across all stories
- [ ] TXXX Update manual test checklist / regression checklist (if needed)
- [ ] TXXX Security hardening
- [ ] TXXX Run quickstart.md validation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently verifiable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - May integrate with US1/US2 but should be independently verifiable

### Within Each User Story

- Verification steps (acceptance scenarios + checklist) MUST be defined before implementation
- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Example parallel tasks (different files, minimal coupling):
Task: "Update OpenAPI in backend/docs/openapi/** for [endpoint]"
Task: "Implement backend handler/usecase in backend/internal/modules/[domain]/..."
Task: "Implement frontend UI in frontend/src/[area]/..."
Task: "Update manual verification checklist in specs/[###-feature-name]/manual-test-checklist.md"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Verify User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Verify independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Verify independently → Deploy/Demo
4. Add User Story 3 → Verify independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and verifiable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
