# Quickstart Guide: Multilingual Dictionary with Vocabulary Game

**Created**: 2025-01-27  
**Feature**: [spec.md](./spec.md)

This guide helps you set up and run the multilingual dictionary and vocabulary game feature locally for development.

## Prerequisites

- **Go**: Version 1.21 or later
- **Node.js**: Version 18 or later
- **npm** or **yarn**: Package manager
- **MySQL/MariaDB**: Version 8.0 or later (or compatible)
- **Git**: Version control

## Backend Setup

### 1. Navigate to Backend Directory

```bash
cd backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Database Setup

Create a MySQL database:

```sql
CREATE DATABASE english_coach CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. Configure Environment

Create a configuration file or set environment variables:

```bash
# Example: backend/.env or backend/config.yaml
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=english_coach
DB_DSN=your_username:your_password@tcp(localhost:3306)/english_coach?charset=utf8mb4&parseTime=True&loc=Local
```

### 5. Run Migrations

Run the database migration to create tables:

```bash
# If migration command exists:
go run cmd/migration/main.go

# Or manually run:
mysql -u your_username -p english_coach < internal/infrastructure/db/migrations/0001_init.sql
```

### 6. Start Backend Server

```bash
go run cmd/api/main.go
```

The backend API should be running on `http://localhost:8080` (or port specified in config).

## Frontend Setup

### 1. Navigate to Frontend Directory

```bash
cd frontend
```

### 2. Install Dependencies

```bash
npm install
# or
yarn install
```

### 3. Configure Environment

Create a `.env` file:

```bash
# frontend/.env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 4. Start Development Server

```bash
npm run dev
# or
yarn dev
```

The frontend should be running on `http://localhost:5173` (default Vite port).

## Docker Setup (Alternative)

### 1. Navigate to Repository Root

```bash
cd /path/to/english-coach
```

### 2. Configure Environment Files

Edit environment files in `deploy/env/dev/`:
- `backend.env`: Backend configuration
- `frontend.env`: Frontend configuration

### 3. Start Services with Docker Compose

```bash
cd deploy/compose
docker-compose --profile dev up -d
```

Services will start:
- Backend API: `http://localhost:8080`
- Frontend: `http://localhost:3000`
- MySQL: `localhost:3306`

## Verify Installation

### 1. Test Backend API

```bash
# Test health endpoint (if available)
curl http://localhost:8080/health

# Test languages endpoint
curl http://localhost:8080/api/v1/reference/languages
```

### 2. Test Frontend

1. Open `http://localhost:5173` in your browser
2. You should see the landing page with "Play Game" and "Dictionary Lookup" buttons

### 3. Test Database Connection

Verify tables were created:

```sql
USE english_coach;
SHOW TABLES;
```

You should see tables like:
- `languages`
- `words`
- `vocab_game_sessions`
- `vocab_game_questions`
- etc.

## Manual Testing Steps

For detailed manual testing procedures, see [Manual Test Checklist](./manual-test-checklist.md).

### Quick Test: Complete User Journey

#### Test Dictionary Lookup (User Story 6)

1. Navigate to Dictionary Lookup page (`/dictionary`)
2. Enter a word in search box (e.g., "hello" or "xin chào")
3. Wait for debounced search (500ms delay)
4. Verify search results appear with pagination if needed
5. Click on a word to see detailed information
6. Verify word detail shows: definitions, translations, examples, pronunciation
7. Click "Quay lại tìm kiếm" to return to search

#### Test Vocabulary Game (User Stories 1-5)

**Step 1: Landing Page (US1)**
1. Navigate to Landing Page (`/`)
2. Verify "Play Game" and "Dictionary Lookup" buttons are visible
3. Click "Play Game" to navigate to game list

**Step 2: Game List (US2)**
1. Verify game list page (`/games`) displays
2. Verify "Học Từ Vựng" (Vocabulary Game) card is visible
3. Click on vocabulary game card

**Step 3: Game Configuration (US3)**
1. Verify configuration page (`/games/vocab/config`) loads
2. Select source language (e.g., English)
3. Select target language (e.g., Vietnamese) - must be different
4. Select mode: "Topic" or "Level"
5. Select a topic or level based on mode
6. Click "Bắt Đầu Chơi" (Start Game)
7. Verify session is created and navigates to play page

**Step 4: Game Play (US4)**
1. Verify game play page loads with first question
2. Verify question loads within 1 second (SC-003)
3. Answer question by clicking option A, B, C, or D
4. Verify answer is submitted and next question appears
5. Complete all questions in session
6. Verify completion screen appears with "View Statistics" button

**Step 5: Statistics (US5)**
1. Click "View Statistics" button
2. Verify statistics page displays:
   - Total questions
   - Correct answers
   - Accuracy percentage
   - Session duration
   - Average response time
3. Test navigation buttons:
   - "Chơi lại" → navigates to config page
   - "Quay lại danh sách game" → navigates to game list

### Expected Behavior

- **Landing Page (US1)**: Two buttons visible and functional, loads within 2 seconds (SC-001)
- **Game List (US2)**: Shows vocabulary game option, back button works
- **Game Configuration (US3)**: Validates inputs, completes in under 30 seconds (SC-002), error messages in Vietnamese
- **Game Play (US4)**: Questions load within 1 second (SC-003), answers recorded, completion screen appears
- **Statistics (US5)**: All statistics display correctly, navigation works
- **Dictionary Lookup (US6)**: Search works, results within 1 second (SC-005), word detail shows complete information

## Common Issues

### Backend Issues

**Problem**: Database connection error
- **Solution**: Check database credentials in config
- **Solution**: Ensure MySQL is running: `sudo systemctl status mysql`

**Problem**: Migration fails
- **Solution**: Ensure database exists and user has permissions
- **Solution**: Check SQL syntax in migration file

### Frontend Issues

**Problem**: Cannot connect to backend API
- **Solution**: Check `VITE_API_BASE_URL` in `.env`
- **Solution**: Ensure backend is running
- **Solution**: Check CORS configuration in backend

**Problem**: Build errors
- **Solution**: Clear node_modules and reinstall: `rm -rf node_modules && npm install`
- **Solution**: Check TypeScript version compatibility

## Next Steps

1. **Seed Data**: Add sample languages, words, and topics to database
2. **API Testing**: Use API endpoints documented in [OpenAPI spec](./contracts/openapi.yaml)
3. **Feature Development**: Follow structure in [plan.md](./plan.md)
4. **Review Documentation**:
   - [Data Model](./data-model.md) for entity relationships
   - [Specification](./spec.md) for requirements
   - [Plan](./plan.md) for implementation details

## Development Scripts

### Backend

```bash
# Run linter
golangci-lint run

# Run formatter
go fmt ./...

# Build
go build -o bin/api cmd/api/main.go
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

## Project Structure Overview

```
english-coach/
├── backend/              # Go backend application
│   ├── cmd/             # Application entry points
│   ├── internal/        # Internal packages
│   │   ├── domain/      # Domain logic (DDD)
│   │   ├── repository/  # Data access implementations
│   │   └── interface/   # HTTP handlers, etc.
│   └── pkg/             # Shared packages
├── frontend/            # React + TypeScript frontend
│   ├── src/
│   │   ├── pages/       # Page components
│   │   ├── features/    # Feature modules
│   │   ├── entities/    # Business entities
│   │   └── shared/      # Shared utilities
├── deploy/              # Docker and deployment configs
└── scripts/             # Development scripts
```

## API Documentation

Interactive API documentation is available via OpenAPI specification:
- File: `contracts/openapi.yaml`
- View with: Swagger UI, Postman, or similar tools

Key endpoints:
- `GET /api/v1/dictionary/search` - Search words
- `GET /api/v1/dictionary/words/{id}` - Get word details
- `POST /api/v1/games/sessions` - Create game session
- `POST /api/v1/games/sessions/{id}/answers` - Submit answer
- `GET /api/v1/statistics/sessions/{id}` - Get session statistics

## Getting Help

- Review [Specification](./spec.md) for feature requirements
- Check [Data Model](./data-model.md) for database structure
- See [Plan](./plan.md) for implementation architecture
- Check backend/STRUCTURE.md and frontend/STRUCTURE.md for code organization

## Troubleshooting Database

### Reset Database

```sql
DROP DATABASE english_coach;
CREATE DATABASE english_coach CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Then re-run migrations.

### Check Tables

```sql
USE english_coach;
SHOW TABLES;
DESCRIBE languages;
DESCRIBE words;
```

### Verify Foreign Keys

```sql
SELECT 
    TABLE_NAME,
    CONSTRAINT_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM
    INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE
    TABLE_SCHEMA = 'english_coach'
    AND REFERENCED_TABLE_NAME IS NOT NULL;
```

