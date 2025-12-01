#!/bin/bash

# Development Environment Startup Script
# Starts all services using docker-compose with dev profile

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
COMPOSE_FILE="$PROJECT_ROOT/deploy/compose/docker-compose.yml"

echo -e "${GREEN}🚀 Starting English Coach Development Environment${NC}"
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

# Start services with dev profile
echo -e "${YELLOW}📦 Starting services with dev profile...${NC}"

# Use dev Dockerfile for backend with live reload
export BACKEND_DOCKERFILE="Dockerfile.dev"

$DOCKER_COMPOSE -f "$COMPOSE_FILE" --profile dev up --build
