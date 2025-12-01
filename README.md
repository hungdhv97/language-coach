# English Coach

A multilingual dictionary application with vocabulary learning games for Vietnamese users.

## Features

- **Dictionary Lookup**: Search and view detailed word information including definitions, translations, examples, and pronunciation
- **Vocabulary Game**: Practice vocabulary through multiple-choice questions filtered by topic or level
- **Game Statistics**: View performance statistics after completing game sessions
- **Multi-language Support**: Support for Vietnamese, English, Chinese, and other languages

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL (migrated from MySQL)
- **Architecture**: Clean Architecture / Domain-Driven Design
- **Logging**: Zap (structured logging)

### Frontend
- **Framework**: React 19 + TypeScript
- **Build Tool**: Vite
- **State Management**: TanStack Query (React Query)
- **Routing**: React Router
- **UI Components**: shadcn/ui
- **Architecture**: Feature-Sliced Design

## Quick Start

### Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- PostgreSQL 12 or later
- Git

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment:
   - Copy `deploy/env/dev/backend.env` and update database credentials
   - Or set environment variables directly

4. Run database migrations:
```bash
go run cmd/migration/main.go
```

5. Start backend server:
```bash
go run cmd/api/main.go
```

Backend API will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Configure environment:
   - Create `.env` file with `VITE_API_BASE_URL=http://localhost:8080/api/v1`

4. Start development server:
```bash
npm run dev
```

Frontend will be available at `http://localhost:5173`

### Docker Setup (Alternative)

1. Configure environment files in `deploy/env/dev/`

2. Start services:
```bash
cd deploy/compose
docker-compose up -d
```

## Project Structure

```
english-coach/
├── backend/              # Go backend application
│   ├── cmd/             # Application entry points
│   │   ├── api/         # API server
│   │   └── migration/   # Migration runner
│   ├── internal/        # Internal packages
│   │   ├── domain/      # Domain logic (DDD)
│   │   ├── repository/  # Data access implementations
│   │   ├── interface/   # HTTP handlers, middleware
│   │   └── shared/      # Shared utilities
│   └── pkg/             # Public packages
├── frontend/            # React + TypeScript frontend
│   └── src/
│       ├── pages/       # Page components
│       ├── features/    # Feature modules
│       ├── entities/    # Business entities
│       └── shared/      # Shared utilities
├── deploy/              # Docker and deployment configs
├── specs/               # Feature specifications
└── scripts/             # Development scripts
```

## API Endpoints

### Reference Data
- `GET /api/v1/reference/languages` - Get all languages
- `GET /api/v1/reference/topics` - Get all topics
- `GET /api/v1/reference/levels?languageId={id}` - Get levels (optionally filtered by language)

### Dictionary
- `GET /api/v1/dictionary/search?q={query}&languageId={id}&limit={limit}&offset={offset}` - Search words
- `GET /api/v1/dictionary/words/{wordId}` - Get word detail

### Game
- `POST /api/v1/games/sessions` - Create game session
- `GET /api/v1/games/sessions/{sessionId}` - Get game session with questions
- `POST /api/v1/games/sessions/{sessionId}/answers` - Submit answer

### Statistics
- `GET /api/v1/statistics/sessions/{sessionId}` - Get session statistics

## Development

### Backend

```bash
# Run linter
golangci-lint run

# Format code
go fmt ./...

# Build
go build -o bin/api cmd/api/main.go

# Run tests (if any)
go test ./...
```

### Frontend

```bash
# Run linter
npm run lint

# Build for production
npm run build

# Preview production build
npm run preview
```

## Testing

See [Manual Test Checklist](specs/001-dictionary-vocab-game/manual-test-checklist.md) for detailed testing procedures.

Quick test:
1. Start backend and frontend
2. Navigate to `http://localhost:5173`
3. Test dictionary lookup and vocabulary game flows

## Documentation

- [Feature Specification](specs/001-dictionary-vocab-game/spec.md)
- [Implementation Plan](specs/001-dictionary-vocab-game/plan.md)
- [Data Model](specs/001-dictionary-vocab-game/data-model.md)
- [Quickstart Guide](specs/001-dictionary-vocab-game/quickstart.md)
- [Manual Test Checklist](specs/001-dictionary-vocab-game/manual-test-checklist.md)
- [API Contract](specs/001-dictionary-vocab-game/contracts/openapi.yaml)

## License

[Add license information]

## Contributing

[Add contributing guidelines]

