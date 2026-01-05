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

# Get environment from first argument (default to dev)
# Check if first argument is dev or prod
if [ "$1" = "dev" ] || [ "$1" = "prod" ]; then
    ENV="$1"
    shift  # Remove ENV argument
else
    ENV="dev"  # Default to dev
fi

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
    echo "Usage: $0 [ENV] [OPTIONS]"
    echo ""
    echo "Environment:"
    echo "  ENV                                    Environment (dev|prod), default: dev"
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
    echo "  $0 dev                                 # Schema + all data for dev (default)"
    echo "  $0 prod                                # Schema + all data for prod"
    echo "  $0 dev --schema-only                   # Schema only for dev"
    echo "  $0 prod --data-only                    # All data only for prod (skip schema)"
    echo "  $0 dev --data-init                     # Data init only for dev"
    echo "  $0 prod --data-word-en                 # English words only for prod"
    echo "  $0 dev --data-init --data-word-en      # Init + English words for dev"
    echo ""
    echo "Note:"
    echo "  Database connection is configured from deploy/env/{ENV}/backend.env"
    echo "  PostgreSQL port is mapped from docker-compose.{ENV}.yml (dev: 5500, prod: 5501)"
    echo "  DATABASE_URL env var can override the default connection"
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

# Show help if requested (before validation)
if [ "$HELP" = true ]; then
    show_help
    exit 0
fi

# Validate environment
if [ "$ENV" != "dev" ] && [ "$ENV" != "prod" ]; then
    echo -e "${RED}❌ Error: Invalid environment '${ENV}'. Must be 'dev' or 'prod'${NC}"
    echo ""
    show_help
    exit 1
fi

# Load environment configuration files
BACKEND_ENV_FILE="$PROJECT_ROOT/deploy/env/${ENV}/backend.env"

if [ ! -f "$BACKEND_ENV_FILE" ]; then
    echo -e "${RED}❌ Error: Backend env file not found at $BACKEND_ENV_FILE${NC}"
    exit 1
fi

# Load environment variables from backend.env
set -a  # automatically export all variables
source "$BACKEND_ENV_FILE"
set +a

# Set PostgreSQL host port based on environment
# These ports match docker-compose.{env}.yml port mappings
if [ "$ENV" = "dev" ]; then
    POSTGRES_PORT=5500
elif [ "$ENV" = "prod" ]; then
    POSTGRES_PORT=5501
fi

# Build DATABASE_URL from environment variables if not already set
if [ -z "$DATABASE_URL" ]; then
    # Use POSTGRES_PORT (host port from docker-compose)
    # Use DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE from backend.env
    # Host is localhost because script runs on host machine
    export DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@localhost:${POSTGRES_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
fi

# Set environment display names
if [ "$ENV" = "dev" ]; then
    ENV_NAME="Development"
elif [ "$ENV" = "prod" ]; then
    ENV_NAME="Production"
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🗄️  Database Migration Script (${ENV_NAME})${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}Environment:${NC} ${ENV} (${ENV_NAME})"
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
