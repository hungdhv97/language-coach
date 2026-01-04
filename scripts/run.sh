#!/bin/bash

# Environment Startup Script
# Starts services using docker-compose with specified environment (dev/prod)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Get environment from argument (default to dev)
ENV="${1:-dev}"

# Validate environment
if [ "$ENV" != "dev" ] && [ "$ENV" != "prod" ]; then
    echo -e "${RED}❌ Error: Invalid environment '${ENV}'. Must be 'dev' or 'prod'${NC}"
    echo ""
    echo "Usage: $0 [dev|prod]"
    echo "  $0 dev    # Start development environment (default)"
    echo "  $0 prod   # Start production environment"
    exit 1
fi

# Set compose file based on environment
COMPOSE_FILE="$PROJECT_ROOT/deploy/compose/docker-compose.${ENV}.yml"

# Set environment display names
if [ "$ENV" = "dev" ]; then
    ENV_NAME="Development"
elif [ "$ENV" = "prod" ]; then
    ENV_NAME="Production"
fi

echo -e "${GREEN}🚀 Starting LexiGo ${ENV_NAME} Environment${NC}"
echo ""

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null && ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Error: docker-compose or docker is not installed${NC}"
    exit 1
fi

# Use docker compose (v2) if available, otherwise fall back to docker-compose (v1)
if docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
elif command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
else
    echo -e "${RED}❌ Error: docker compose is not available${NC}"
    exit 1
fi

# Change to project root
cd "$PROJECT_ROOT"

# Check if compose file exists
if [ ! -f "$COMPOSE_FILE" ]; then
    echo -e "${RED}❌ Error: Docker compose file not found at $COMPOSE_FILE${NC}"
    exit 1
fi

# Load docker-compose environment variables (for COMPOSE_PROJECT_NAME if needed)
ENV_FILE="$PROJECT_ROOT/deploy/env/${ENV}/docker-compose.env"
if [ -f "$ENV_FILE" ]; then
    # Load environment variables from docker-compose.env
    set -a  # automatically export all variables
    source "$ENV_FILE"
    set +a
fi

# Start services
echo -e "${YELLOW}📦 Starting services for ${ENV} environment...${NC}"

$DOCKER_COMPOSE -f "$COMPOSE_FILE" up --build -d

