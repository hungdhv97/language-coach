# Tasks: Multilingual Dictionary with Vocabulary Game

**Input**: Design documents from `/specs/001-dictionary-vocab-game/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Manual testing discipline per constitution - no automated tests required.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `backend/internal/`, `backend/cmd/`, `backend/pkg/`
- **Frontend**: `frontend/src/`
- All paths are relative to repository root

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure aligned with the constitution

- [x] T001 Create backend project structure following `backend/STRUCTURE.md` in `backend/`
- [x] T002 Initialize Go module in `backend/go.mod` with Go 1.21+ dependencies
- [x] T003 [P] Create frontend directory structure following `frontend/STRUCTURE.md` in `frontend/src/`
- [x] T004 [P] Configure ESLint and Prettier for frontend in `frontend/.eslintrc.json` and `frontend/.prettierrc`
- [x] T005 [P] Configure Go linting and formatting tools in `backend/.golangci.yml` (no code may merge with lint/build errors)
- [x] T006 Create environment configuration files in `deploy/env/dev/backend.env` and `deploy/env/dev/frontend.env`
- [x] T007 [P] Create Dockerfiles in `deploy/docker/backend/Dockerfile` and `deploy/docker/frontend/Dockerfile`
- [x] T008 Create docker-compose configuration in `deploy/compose/docker-compose.yml`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T009 Setup MySQL database connection in `backend/internal/infrastructure/db/mysql.go` with connection pooling
- [x] T010 Create database migration runner in `backend/cmd/migration/main.go` to execute `0001_init.sql`
- [x] T011 Verify database schema by running migration `backend/internal/infrastructure/db/migrations/0001_init.sql`
- [x] T012 [P] Create unified error schema in `backend/internal/shared/response/error.go` with `{ code: string, message: string, details?: unknown }`
- [x] T013 [P] Implement centralized error handler middleware in `backend/internal/interface/http/middleware/error_handler.go`
- [x] T014 [P] Setup structured logging with zap in `backend/internal/infrastructure/logger/zap_logger.go`
- [x] T015 [P] Create request logging middleware in `backend/internal/interface/http/middleware/logger.go`
- [x] T016 [P] Setup HTTP router and server in `backend/internal/interface/http/server.go`
- [x] T017 [P] Create CORS middleware in `backend/internal/interface/http/middleware/cors.go`
- [x] T018 [P] Create input validation package in `backend/pkg/validator/validator.go` using go-playground/validator
- [x] T019 [P] Create HTTP client configuration in `frontend/src/shared/api/config.ts`
- [x] T020 [P] Create HTTP client wrapper in `frontend/src/shared/api/http-client.ts` with axios/fetch
- [x] T021 [P] Create error interceptor in `frontend/src/shared/api/interceptors/error-interceptor.ts`
- [x] T022 [P] Setup React Router in `frontend/src/app/router/routes.tsx` with route definitions
- [x] T023 [P] Create AppProviders component in `frontend/src/app/providers/AppProviders.tsx` for global providers
- [x] T024 [P] Create i18n configuration for Vietnamese error messages in `frontend/src/shared/lib/i18n/i18n.ts`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Landing Page Navigation (Priority: P1) 🎯 MVP

**Goal**: Users see a landing page with two action buttons ("Play Game" and "Dictionary Lookup") and can navigate to respective pages.

**Independent Test**: Navigate to homepage, verify both buttons are visible and clickable. Click each button and verify navigation works. This delivers immediate value by enabling feature discovery.

### Implementation for User Story 1

- [x] T025 [US1] Create LandingPage component in `frontend/src/pages/LandingPage.tsx` with two prominent action buttons
- [x] T026 [US1] Add route for landing page (/) in `frontend/src/app/router/routes.tsx`
- [x] T027 [US1] Style landing page buttons following shared design system in `frontend/src/pages/LandingPage.tsx`
- [x] T028 [US1] Add navigation handler for "Play Game" button to route `/games` in `frontend/src/pages/LandingPage.tsx`
- [x] T029 [US1] Add navigation handler for "Dictionary Lookup" button to route `/dictionary` in `frontend/src/pages/LandingPage.tsx`
- [x] T030 [US1] Ensure landing page loads within 2 seconds per SC-001 in `frontend/src/pages/LandingPage.tsx`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently - users can navigate from landing page to game list or dictionary lookup.

---

## Phase 4: User Story 2 - Game List Display (Priority: P1)

**Goal**: Users see a list of available games and can select the vocabulary game to proceed to configuration.

**Independent Test**: Click "Play Game" from landing page, verify game list displays with vocabulary game option. Click vocabulary game and verify navigation to configuration page.

### Implementation for User Story 2

- [x] T031 [US2] Create GameListPage component in `frontend/src/pages/game/GameListPage.tsx` to display available games
- [x] T032 [US2] Add route for game list page (/games) in `frontend/src/app/router/routes.tsx`
- [x] T033 [US2] Create game list item component showing vocabulary game option in `frontend/src/pages/game/GameListPage.tsx`
- [x] T034 [US2] Add navigation handler for vocabulary game selection to route `/games/vocab/config` in `frontend/src/pages/game/GameListPage.tsx`
- [x] T035 [US2] Style game list following shared design system in `frontend/src/pages/game/GameListPage.tsx`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - users can navigate to game list and select vocabulary game.

---

## Phase 5: User Story 3 - Vocabulary Game Configuration (Priority: P1)

**Goal**: Users configure game session by selecting source language, target language, and mode (topic or level), then start the game.

**Independent Test**: Select vocabulary game, configure languages (source ≠ target), select topic or level, click start. Verify game session is created. This delivers value by allowing users to customize their learning experience.

### Backend Implementation for User Story 3

- [x] T036 [P] [US3] Create Language domain model in `backend/internal/domain/dictionary/model/language.go`
- [x] T037 [P] [US3] Create Topic domain model in `backend/internal/domain/dictionary/model/topic.go`
- [x] T038 [P] [US3] Create Level domain model in `backend/internal/domain/dictionary/model/level.go`
- [x] T039 [P] [US3] Create Language repository interface in `backend/internal/domain/dictionary/port/repository.go`
- [x] T040 [P] [US3] Create Topic repository interface in `backend/internal/domain/dictionary/port/repository.go`
- [x] T041 [P] [US3] Create Level repository interface in `backend/internal/domain/dictionary/port/repository.go`
- [x] T042 [US3] Implement Language repository in `backend/internal/repository/dictionary_pg.go` for MySQL queries
- [x] T043 [US3] Implement Topic repository in `backend/internal/repository/dictionary_pg.go` for MySQL queries
- [x] T044 [US3] Implement Level repository in `backend/internal/repository/dictionary_pg.go` for MySQL queries
- [x] T045 [US3] Create GET /api/v1/reference/languages endpoint handler in `backend/internal/interface/http/handler/dictionary_handler.go`
- [x] T046 [US3] Create GET /api/v1/reference/topics endpoint handler in `backend/internal/interface/http/handler/dictionary_handler.go`
- [x] T047 [US3] Create GET /api/v1/reference/levels endpoint handler with languageId query param in `backend/internal/interface/http/handler/dictionary_handler.go`
- [x] T048 [US3] Register dictionary routes in `backend/internal/app/router.go` for reference data endpoints
- [x] T049 [US3] Create GameSession domain model in `backend/internal/domain/game/model/game_session.go`
- [x] T050 [US3] Create game session repository interface in `backend/internal/domain/game/port/repository.go`
- [x] T051 [US3] Create CreateGameSessionRequest DTO in `backend/internal/domain/game/dto/create_session.go`
- [x] T052 [US3] Implement validation for CreateGameSessionRequest (source ≠ target, topic XOR level) in `backend/internal/domain/game/dto/create_session.go`
- [x] T053 [US3] Create CreateGameSession use case in `backend/internal/domain/game/usecase/command/create_session.go`
- [x] T054 [US3] Implement game session repository in `backend/internal/repository/game_pg.go` for MySQL queries
- [x] T055 [US3] Create POST /api/v1/games/sessions endpoint handler in `backend/internal/interface/http/handler/game_handler.go`
- [x] T056 [US3] Register game routes in `backend/internal/app/router.go` for session creation endpoint
- [x] T057 [US3] Add error handling for insufficient words scenario (FR-026) in `backend/internal/domain/game/usecase/command/create_session.go`
- [x] T058 [US3] Add Vietnamese error messages for validation errors (FR-025) in `backend/internal/interface/http/handler/game_handler.go`

### Frontend Implementation for User Story 3

- [x] T059 [P] [US3] Create Language type definitions in `frontend/src/entities/dictionary/model/dictionary.types.ts`
- [x] T060 [P] [US3] Create Topic type definitions in `frontend/src/entities/dictionary/model/dictionary.types.ts`
- [x] T061 [P] [US3] Create Level type definitions in `frontend/src/entities/dictionary/model/dictionary.types.ts`
- [x] T062 [P] [US3] Create dictionary API endpoints in `frontend/src/entities/dictionary/api/dictionary.endpoints.ts`
- [x] T063 [P] [US3] Create reference data queries in `frontend/src/entities/dictionary/api/dictionary.queries.ts` for languages, topics, levels
- [x] T064 [P] [US3] Create game session types in `frontend/src/entities/game/model/game.types.ts`
- [x] T065 [P] [US3] Create game API endpoints in `frontend/src/entities/game/api/game.endpoints.ts`
- [x] T066 [US3] Create GameConfigPage component in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T067 [US3] Create language selection dropdowns (source and target) in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T068 [US3] Create mode selector (topic or level) in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T069 [US3] Create topic selection interface that appears when topic mode is selected in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T070 [US3] Create level selection interface that appears when level mode is selected in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T071 [US3] Implement validation logic: source language ≠ target language (FR-010) in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T072 [US3] Implement validation logic: topic XOR level required (FR-011) in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T073 [US3] Create game configuration form component in `frontend/src/features/game/components/GameConfigForm.tsx`
- [x] T074 [US3] Create hook for game configuration in `frontend/src/features/game/hooks/useGameConfig.ts`
- [x] T075 [US3] Add game session creation mutation in `frontend/src/features/game/api/game.mutations.ts`
- [x] T076 [US3] Add start button handler that creates game session in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T077 [US3] Navigate to game play page after session creation in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T078 [US3] Add route for game config page (/games/vocab/config) in `frontend/src/app/router/routes.tsx`
- [x] T079 [US3] Ensure configuration can be completed in under 30 seconds per SC-002 in `frontend/src/pages/game/GameConfigPage.tsx`
- [x] T080 [US3] Display Vietnamese error messages for validation errors (FR-025) in `frontend/src/pages/game/GameConfigPage.tsx`

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should work independently - users can configure and start a game session.

---

## Phase 6: User Story 4 - Playing Vocabulary Game (Priority: P1)

**Goal**: Users play vocabulary game by answering multiple-choice questions (A, B, C, D) with word in source language and translations as options.

**Independent Test**: Start configured game session, answer multiple-choice questions, complete all questions. Verify questions load within 1 second, answers are recorded, completion screen appears. This delivers immediate learning value through interactive vocabulary practice.

### Backend Implementation for User Story 4

- [x] T081 [P] [US4] Create Word domain model in `backend/internal/domain/dictionary/model/word.go`
- [x] T082 [P] [US4] Create Sense domain model in `backend/internal/domain/dictionary/model/sense.go`
- [x] T083 [US4] Create Word repository interface in `backend/internal/domain/dictionary/port/repository.go`
- [x] T084 [US4] Implement Word repository with search methods in `backend/internal/repository/dictionary_pg.go`
- [x] T085 [US4] Create GameQuestion domain model in `backend/internal/domain/game/model/game_question.go`
- [x] T086 [US4] Create GameQuestionOption domain model in `backend/internal/domain/game/model/game_question.go`
- [x] T087 [US4] Create question generator interface in `backend/internal/domain/game/port/question_generator.go`
- [x] T088 [US4] Implement question generator service in `backend/internal/domain/game/service/question_generator_service.go` to generate 4 options with 1 correct answer
- [x] T089 [US4] Update CreateGameSession use case to generate questions upfront in `backend/internal/domain/game/usecase/command/create_session.go`
- [x] T090 [US4] Create game question repository interface in `backend/internal/domain/game/port/repository.go`
- [x] T091 [US4] Implement game question repository in `backend/internal/repository/game_pg.go` for MySQL queries
- [x] T092 [US4] Create GET /api/v1/games/sessions/{sessionId} endpoint handler in `backend/internal/interface/http/handler/game_handler.go`
- [x] T093 [US4] Create GameAnswer domain model in `backend/internal/domain/game/model/game_answer.go`
- [x] T094 [US4] Create SubmitAnswerRequest DTO in `backend/internal/domain/game/dto/submit_answer.go`
- [x] T095 [US4] Create SubmitAnswer use case in `backend/internal/domain/game/usecase/command/answer_question.go`
- [x] T096 [US4] Create game answer repository interface in `backend/internal/domain/game/port/repository.go`
- [x] T097 [US4] Implement game answer repository in `backend/internal/repository/game_pg.go` for MySQL queries
- [x] T098 [US4] Create POST /api/v1/games/sessions/{sessionId}/answers endpoint handler in `backend/internal/interface/http/handler/game_handler.go`
- [x] T099 [US4] Create EndGameSession use case in `backend/internal/domain/game/usecase/command/end_session.go`
- [x] T100 [US4] Update game session repository to support ending session in `backend/internal/repository/game_pg.go`
- [x] T101 [US4] Add structured logging for game session start/end in `backend/internal/domain/game/usecase/command/create_session.go` and `end_session.go`
- [x] T102 [US4] Add structured logging for question answers in `backend/internal/domain/game/usecase/command/answer_question.go`
- [x] T103 [US4] Register game routes for session retrieval and answer submission in `backend/internal/app/router.go`
- [x] T104 [US4] Optimize question generation queries to avoid N+1 problem in `backend/internal/domain/game/service/question_generator_service.go`
- [x] T105 [US4] Ensure question generation completes within 1 second per SC-003 in `backend/internal/domain/game/service/question_generator_service.go`

### Frontend Implementation for User Story 4

- [x] T106 [P] [US4] Create GameQuestion type definitions in `frontend/src/entities/game/model/game.types.ts`
- [x] T107 [P] [US4] Create GameAnswer type definitions in `frontend/src/entities/game/model/game.types.ts`
- [x] T108 [P] [US4] Create game session queries in `frontend/src/features/game/api/game.queries.ts` for session retrieval
- [x] T109 [P] [US4] Create answer submission mutation in `frontend/src/features/game/api/game.mutations.ts`
- [x] T110 [US4] Create GamePlayPage component in `frontend/src/pages/game/GamePlayPage.tsx`
- [x] T111 [US4] Create GameQuestion component to display question with 4 options (A, B, C, D) in `frontend/src/features/game/components/GameQuestion.tsx`
- [x] T112 [US4] Display source word prominently in `frontend/src/features/game/components/GameQuestion.tsx`
- [x] T113 [US4] Display four answer options as A, B, C, D buttons in `frontend/src/features/game/components/GameQuestion.tsx`
- [x] T114 [US4] Create hook for game session management in `frontend/src/features/game/hooks/useGameSession.ts`
- [x] T115 [US4] Implement answer selection handler in `frontend/src/features/game/components/GameQuestion.tsx`
- [x] T116 [US4] Track response time for each answer in `frontend/src/features/game/hooks/useGameSession.ts`
- [x] T117 [US4] Submit answer to backend API in `frontend/src/features/game/hooks/useGameSession.ts`
- [x] T118 [US4] Progress to next question after answer submission in `frontend/src/pages/game/GamePlayPage.tsx`
- [x] T119 [US4] Create game completion screen component in `frontend/src/pages/game/GamePlayPage.tsx`
- [x] T120 [US4] Display completion screen when all questions are answered in `frontend/src/pages/game/GamePlayPage.tsx`
- [x] T121 [US4] Add "View Statistics" button on completion screen (FR-018) in `frontend/src/pages/game/GamePlayPage.tsx`
- [x] T122 [US4] Add route for game play page (/games/vocab/play/:sessionId) in `frontend/src/app/router/routes.tsx`
- [x] T123 [US4] Ensure questions load within 1 second per SC-003 in `frontend/src/pages/game/GamePlayPage.tsx`
- [x] T124 [US4] Implement loading states during question transitions in `frontend/src/pages/game/GamePlayPage.tsx`

**Checkpoint**: At this point, User Stories 1-4 should work independently - users can play a complete game session from start to finish.

---

## Phase 7: User Story 5 - Viewing Game Statistics (Priority: P2)

**Goal**: Users can view game statistics including total questions, correct answers, accuracy percentage, and session duration after completing a game.

**Independent Test**: Complete a game session, click "View Statistics", verify statistics display correctly. This delivers value by providing learning feedback and progress tracking.

### Backend Implementation for User Story 5

- [x] T125 [P] [US5] Create SessionStatistics DTO in `backend/internal/domain/game/dto/session_statistics.go`
- [x] T126 [US5] Create GetSessionStatistics use case in `backend/internal/domain/game/usecase/query/get_statistics.go`
- [x] T127 [US5] Implement statistics calculation logic (accuracy, duration, average response time) in `backend/internal/domain/game/usecase/query/get_statistics.go`
- [x] T128 [US5] Create GET /api/v1/statistics/sessions/{sessionId} endpoint handler in `backend/internal/interface/http/handler/statistics_handler.go`
- [x] T129 [US5] Register statistics routes in `backend/internal/app/router.go` for session statistics endpoint
- [x] T130 [US5] Ensure statistics display complete and accurate information per SC-008 in `backend/internal/domain/game/usecase/query/get_statistics.go`

### Frontend Implementation for User Story 5

- [x] T131 [P] [US5] Create SessionStatistics type definitions in `frontend/src/entities/game/model/game.types.ts`
- [x] T132 [P] [US5] Create statistics API endpoints in `frontend/src/features/game/api/game.endpoints.ts`
- [x] T133 [P] [US5] Create statistics queries in `frontend/src/features/game/api/game.queries.ts`
- [x] T134 [US5] Create GameStatistics component to display statistics in `frontend/src/features/game/components/GameStatistics.tsx`
- [x] T135 [US5] Display total questions, correct answers, accuracy percentage in `frontend/src/features/game/components/GameStatistics.tsx`
- [x] T136 [US5] Display session duration in `frontend/src/features/game/components/GameStatistics.tsx`
- [x] T137 [US5] Create GameStatisticsPage component in `frontend/src/pages/game/GameStatisticsPage.tsx`
- [x] T138 [US5] Add navigation from completion screen to statistics page in `frontend/src/pages/game/GamePlayPage.tsx`
- [x] T139 [US5] Add route for statistics page (/games/vocab/statistics/:sessionId) in `frontend/src/app/router/routes.tsx`
- [x] T140 [US5] Add navigation options (play again, return to game list) in `frontend/src/pages/game/GameStatisticsPage.tsx`

**Checkpoint**: At this point, User Stories 1-5 should work independently - users can view statistics after completing games.

---

## Phase 8: User Story 6 - Dictionary Lookup (Priority: P2)

**Goal**: Users can search for words in the multilingual dictionary and view detailed information including definitions, translations, examples, and pronunciation.

**Independent Test**: Click "Dictionary Lookup" from landing page, enter a word, verify search results display. Click a word, verify detailed information shows. This delivers immediate value by enabling word lookup functionality.

### Backend Implementation for User Story 6

- [x] T141 [P] [US6] Create Sense repository interface methods in `backend/internal/domain/dictionary/port/repository.go`
- [x] T142 [US6] Implement word search with multiple strategies (lemma, normalized, search_key) in `backend/internal/repository/dictionary_pg.go`
- [x] T143 [US6] Create dictionary service for word lookup in `backend/internal/domain/dictionary/service/dictionary_service.go`
- [x] T144 [US6] Create WordDetail DTO with senses, translations, examples in `backend/internal/domain/dictionary/dto/word_detail.go`
- [x] T145 [US6] Create GET /api/v1/dictionary/search endpoint handler with pagination in `backend/internal/interface/http/handler/dictionary_handler.go`
- [x] T146 [US6] Create GET /api/v1/dictionary/words/{wordId} endpoint handler in `backend/internal/interface/http/handler/dictionary_handler.go`
- [x] T147 [US6] Register dictionary routes for search and word details in `backend/internal/app/router.go`
- [x] T148 [US6] Optimize word search queries using existing indexes in `backend/internal/repository/dictionary_pg.go`
- [x] T149 [US6] Ensure dictionary lookup returns results within 1 second for 95% of searches per SC-005 in `backend/internal/repository/dictionary_pg.go`
- [x] T150 [US6] Implement pagination for search results (limit, offset) in `backend/internal/interface/http/handler/dictionary_handler.go`
- [x] T151 [US6] Handle empty search results gracefully in `backend/internal/interface/http/handler/dictionary_handler.go`

### Frontend Implementation for User Story 6

- [x] T152 [P] [US6] Create Word type definitions with detailed structure in `frontend/src/entities/dictionary/model/dictionary.types.ts`
- [x] T153 [P] [US6] Create WordDetail type definitions in `frontend/src/entities/dictionary/model/dictionary.types.ts`
- [x] T154 [P] [US6] Create dictionary search API endpoints in `frontend/src/entities/dictionary/api/dictionary.endpoints.ts`
- [x] T155 [P] [US6] Create dictionary search queries in `frontend/src/entities/dictionary/api/dictionary.queries.ts`
- [x] T156 [US6] Create DictionaryLookupPage component in `frontend/src/pages/dictionary/DictionaryLookupPage.tsx`
- [x] T157 [US6] Create DictionarySearch component with search input in `frontend/src/features/dictionary/components/DictionarySearch.tsx`
- [x] T158 [US6] Display search results list in `frontend/src/features/dictionary/components/DictionarySearch.tsx`
- [x] T159 [US6] Implement search debouncing for performance in `frontend/src/features/dictionary/components/DictionarySearch.tsx`
- [x] T160 [US6] Create WordDetail component to display word information in `frontend/src/features/dictionary/components/WordDetail.tsx`
- [x] T161 [US6] Display definitions, translations, examples, pronunciation in `frontend/src/features/dictionary/components/WordDetail.tsx`
- [x] T162 [US6] Add route for dictionary lookup page (/dictionary) in `frontend/src/app/router/routes.tsx`
- [x] T163 [US6] Add route for word detail page (/dictionary/words/:wordId) in `frontend/src/app/router/routes.tsx`
- [x] T164 [US6] Handle empty search results with user-friendly message in `frontend/src/features/dictionary/components/DictionarySearch.tsx`
- [x] T165 [US6] Handle invalid/empty search queries gracefully in `frontend/src/features/dictionary/components/DictionarySearch.tsx`

**Checkpoint**: At this point, all User Stories 1-6 should work independently - users can lookup words and play vocabulary games.

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories and ensure constitution compliance

- [x] T166 [P] Add structured logging for all critical flows (game sessions, dictionary lookups) across all handlers
- [x] T167 [P] Add request ID/correlation ID to all log entries in `backend/internal/interface/http/middleware/logger.go`
- [x] T168 [P] Review and add error handling for network errors during game play in `frontend/src/features/game/hooks/useGameSession.ts`
- [x] T169 [P] Add loading states for all async operations in frontend components
- [x] T170 [P] Add error states with Vietnamese error messages throughout frontend (FR-025)
- [x] T171 [P] Implement navigation back to game list or landing page from any page (FR-027) in `frontend/src/app/router/routes.tsx`
- [x] T172 [P] Review database query performance and add indexes if needed beyond existing schema
- [x] T173 [P] Add pagination UI components for dictionary search results in `frontend/src/features/dictionary/components/DictionarySearch.tsx`
- [ ] T174 [P] Add pagination UI components for game session history (if implemented) in frontend
- [x] T175 [P] Update quickstart.md with manual testing steps for all P1 user journeys in `specs/001-dictionary-vocab-game/quickstart.md`
- [x] T176 [P] Create manual test checklist document for P1 user stories in `specs/001-dictionary-vocab-game/`
- [ ] T177 [P] Perform manual smoke tests for critical flows: landing page, game configuration, game play, dictionary lookup
- [x] T178 [P] Review all error messages are in Vietnamese and user-friendly (FR-025) across backend and frontend
- [x] T179 [P] Verify all external inputs are validated (FR-023) in all handlers and forms
- [x] T180 [P] Verify unified error schema is used consistently (FR-024) across all endpoints
- [x] T181 [P] Review code for magic values and extract to constants (game question count, timeouts)
- [x] T182 [P] Code cleanup and refactoring to maintain code quality & consistency per constitution
- [x] T183 [P] Update README.md with setup instructions if needed in repository root

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-8)**: All depend on Foundational phase completion
  - User stories can then proceed sequentially in priority order (US1 → US2 → US3 → US4 → US5 → US6)
  - US3, US4, US6 have backend dependencies that must complete before frontend work
  - Frontend components can be built in parallel within each story if backend APIs are ready
- **Polish (Phase 9)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Depends on US1 for navigation from landing page
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - Requires backend reference data endpoints first
- **User Story 4 (P1)**: Depends on US3 for game session creation - Requires backend question generation
- **User Story 5 (P2)**: Depends on US4 for completed game sessions
- **User Story 6 (P2)**: Can start after Foundational (Phase 2) - Requires backend dictionary endpoints

### Within Each User Story

- Backend models before repository interfaces
- Repository interfaces before implementations
- Services before use cases
- Use cases before handlers
- Backend APIs before frontend API clients
- Frontend types before components
- Components before pages
- Core implementation before integration

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Backend domain models within a story marked [P] can run in parallel
- Frontend type definitions marked [P] can run in parallel
- Frontend components can be built in parallel once backend APIs are ready
- Polish tasks marked [P] can run in parallel after user stories complete

---

## Parallel Example: User Story 3

```bash
# Backend models can be created in parallel:
T036: Create Language domain model
T037: Create Topic domain model
T038: Create Level domain model

# Frontend types can be created in parallel:
T059: Create Language type definitions
T060: Create Topic type definitions
T061: Create Level type definitions

# API endpoints can be created in parallel (after models):
T045: GET /api/v1/reference/languages endpoint
T046: GET /api/v1/reference/topics endpoint
T047: GET /api/v1/reference/levels endpoint
```

---

## Parallel Example: User Story 4

```bash
# Domain models can be created in parallel:
T081: Create Word domain model
T082: Create Sense domain model
T085: Create GameQuestion domain model
T086: Create GameQuestionOption domain model

# Frontend types can be created in parallel:
T106: Create GameQuestion type definitions
T107: Create GameAnswer type definitions
```

---

## Implementation Strategy

### MVP First (User Stories 1-4 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Landing Page)
4. Complete Phase 4: User Story 2 (Game List)
5. Complete Phase 5: User Story 3 (Game Configuration)
6. Complete Phase 6: User Story 4 (Playing Game)
7. **STOP and VALIDATE**: Test complete flow from landing page through game completion
8. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (Navigation works)
3. Add User Story 2 → Test independently → Deploy/Demo (Game selection works)
4. Add User Story 3 → Test independently → Deploy/Demo (Configuration works)
5. Add User Story 4 → Test independently → Deploy/Demo (MVP - Full game playable!)
6. Add User Story 5 → Test independently → Deploy/Demo (Statistics available)
7. Add User Story 6 → Test independently → Deploy/Demo (Dictionary available)
8. Add Polish phase → Final release

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Frontend only)
   - Developer B: User Story 3 Backend (Reference data APIs)
   - Developer C: User Story 6 Backend (Dictionary APIs)
3. After User Stories 1-3 complete:
   - Developer A: User Story 2 (Frontend)
   - Developer B: User Story 3 Frontend (Configuration)
   - Developer C: User Story 4 Backend (Game logic)
4. After User Story 4 Backend complete:
   - Developer A: User Story 4 Frontend (Game play)
   - Developer B: User Story 5 (Statistics)
   - Developer C: User Story 6 Frontend (Dictionary)
5. All complete → Polish phase together

---

## Notes

- **[P] tasks** = different files, no dependencies - can run in parallel
- **[Story] label** maps task to specific user story for traceability (US1-US6)
- Each user story should be independently completable and testable
- Manual testing required per constitution - no automated tests
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Backend must be ready before frontend can consume APIs
- All error messages must be in Vietnamese (FR-025)
- All inputs must be validated (FR-023)
- All errors must follow unified schema (FR-024)
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence

---

## Task Summary

**Total Tasks**: 183

**Tasks by Phase**:
- Phase 1 (Setup): 8 tasks
- Phase 2 (Foundational): 16 tasks
- Phase 3 (US1 - Landing Page): 6 tasks
- Phase 4 (US2 - Game List): 5 tasks
- Phase 5 (US3 - Game Configuration): 45 tasks (25 backend + 20 frontend)
- Phase 6 (US4 - Playing Game): 44 tasks (25 backend + 19 frontend)
- Phase 7 (US5 - Statistics): 16 tasks (6 backend + 10 frontend)
- Phase 8 (US6 - Dictionary): 25 tasks (11 backend + 14 frontend)
- Phase 9 (Polish): 18 tasks

**Tasks by User Story**:
- US1: 6 tasks
- US2: 5 tasks
- US3: 45 tasks
- US4: 44 tasks
- US5: 16 tasks
- US6: 25 tasks

**Parallel Opportunities**: 80+ tasks marked with [P] can run in parallel

**Suggested MVP Scope**: User Stories 1-4 (Landing Page through Complete Game Play) = 100 tasks

**Independent Test Criteria**:
- US1: Navigate to homepage, verify buttons, click and navigate
- US2: Navigate to game list, verify vocabulary game, click to config
- US3: Configure languages, mode, topic/level, start session
- US4: Play game, answer questions, complete session
- US5: View statistics after game completion
- US6: Search dictionary, view word details

