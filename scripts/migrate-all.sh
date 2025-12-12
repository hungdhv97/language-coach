#!/bin/bash

# Migration Script
# This script runs database migrations (schema + data) using Go
# Default: Runs schema migration then all data migrations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

# Default database URL (can be overridden by DATABASE_URL env var)
DEFAULT_DATABASE_URL="postgres://postgres:postgres@localhost:5500/english_coach?sslmode=disable"

# Parse flags
SCHEMA_ONLY=false
DATA_ONLY=false
DATA_INIT=false
DATA_WORD_EN=false
DATA_WORD_VI=false
DATA_WORD_ZH=false
HELP=false

show_help() {
    echo -e "${BLUE}🗄️  Database Migration Script${NC}"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --schema-only          Run schema migration only"
    echo "  --data-only            Run data migration only (skip schema)"
    echo "  --data-init            Run data migration with --init flag (initial metadata only)"
    echo "  --data-word-en         Run data migration with --word-en flag (English words only)"
    echo "  --data-word-vi         Run data migration with --word-vi flag (Vietnamese words only)"
    echo "  --data-word-zh         Run data migration with --word-zh flag (Chinese words only)"
    echo "  --help, -h             Show this help message"
    echo ""
    echo "Default behavior (no flags):"
    echo "  - Run schema migration"
    echo "  - Then run all data migrations (init + word-en + word-vi + word-zh)"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Schema + all data (default)"
    echo "  $0 --schema-only                      # Schema only"
    echo "  $0 --data-only                        # All data only (skip schema)"
    echo "  $0 --data-init                        # Data init only"
    echo "  $0 --data-word-en                     # English words only"
    echo "  $0 --data-init --data-word-en         # Init + English words"
    echo ""
    echo "Environment:"
    echo "  DATABASE_URL                          Database connection string (optional)"
    echo "                                        Default: $DEFAULT_DATABASE_URL"
    echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --schema-only)
            SCHEMA_ONLY=true
            shift
            ;;
        --data-only)
            DATA_ONLY=true
            shift
            ;;
        --data-init)
            DATA_INIT=true
            shift
            ;;
        --data-word-en)
            DATA_WORD_EN=true
            shift
            ;;
        --data-word-vi)
            DATA_WORD_VI=true
            shift
            ;;
        --data-word-zh)
            DATA_WORD_ZH=true
            shift
            ;;
        --help|-h)
            HELP=true
            shift
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
done

# Show help if requested
if [ "$HELP" = true ]; then
    show_help
    exit 0
fi

# Set database URL if not already set
if [ -z "$DATABASE_URL" ]; then
    export DATABASE_URL="$DEFAULT_DATABASE_URL"
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🗄️  Database Migration Script${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}Database URL:${NC} $DATABASE_URL"
echo ""

# Check if backend directory exists
if [ ! -d "$BACKEND_DIR" ]; then
    echo -e "${RED}Error: Backend directory not found at $BACKEND_DIR${NC}"
    exit 1
fi

# Determine what to run
RUN_SCHEMA=false
RUN_DATA=false
DATA_FLAGS=""

if [ "$SCHEMA_ONLY" = true ]; then
    RUN_SCHEMA=true
    RUN_DATA=false
elif [ "$DATA_ONLY" = true ]; then
    RUN_SCHEMA=false
    RUN_DATA=true
    # Build data flags
    if [ "$DATA_INIT" = true ]; then
        DATA_FLAGS="$DATA_FLAGS --init"
    fi
    if [ "$DATA_WORD_EN" = true ]; then
        DATA_FLAGS="$DATA_FLAGS --word-en"
    fi
    if [ "$DATA_WORD_VI" = true ]; then
        DATA_FLAGS="$DATA_FLAGS --word-vi"
    fi
    if [ "$DATA_WORD_ZH" = true ]; then
        DATA_FLAGS="$DATA_FLAGS --word-zh"
    fi
    # If no data flags specified, run all
    if [ -z "$DATA_FLAGS" ]; then
        DATA_FLAGS=""  # Empty means run all (default behavior of data migration)
    fi
else
    # Default: run schema then all data
    RUN_SCHEMA=true
    RUN_DATA=true
    # Check if specific data flags were provided
    if [ "$DATA_INIT" = true ] || [ "$DATA_WORD_EN" = true ] || [ "$DATA_WORD_VI" = true ] || [ "$DATA_WORD_ZH" = true ]; then
        if [ "$DATA_INIT" = true ]; then
            DATA_FLAGS="$DATA_FLAGS --init"
        fi
        if [ "$DATA_WORD_EN" = true ]; then
            DATA_FLAGS="$DATA_FLAGS --word-en"
        fi
        if [ "$DATA_WORD_VI" = true ]; then
            DATA_FLAGS="$DATA_FLAGS --word-vi"
        fi
        if [ "$DATA_WORD_ZH" = true ]; then
            DATA_FLAGS="$DATA_FLAGS --word-zh"
        fi
    else
        # No specific data flags, run all (empty DATA_FLAGS)
        DATA_FLAGS=""
    fi
fi

# Run Schema Migration
if [ "$RUN_SCHEMA" = true ]; then
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}Step 1: Running Schema Migration${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    cd "$BACKEND_DIR"
    
    if ! go run cmd/migration/schema/main.go; then
        echo -e "${RED}Schema migration failed!${NC}"
        exit 1
    fi
    
    echo ""
    echo -e "${GREEN}✓ Schema migration completed successfully!${NC}"
    echo ""
fi

# Run Data Migration
if [ "$RUN_DATA" = true ]; then
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}Step 2: Running Data Migration${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    cd "$BACKEND_DIR"
    
    # Build the command
    DATA_CMD="go run cmd/migration/data/main.go"
    if [ -n "$DATA_FLAGS" ]; then
        DATA_CMD="$DATA_CMD $DATA_FLAGS"
    fi
    
    echo -e "${YELLOW}Running: $DATA_CMD${NC}"
    echo ""
    
    if ! eval "$DATA_CMD"; then
        echo -e "${RED}Data migration failed!${NC}"
        exit 1
    fi
    
    echo ""
    echo -e "${GREEN}✓ Data migration completed successfully!${NC}"
    echo ""
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✓ All migrations completed successfully!${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
