#!/bin/bash

# Development Environment Startup Script
# Starts development services using docker-compose

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Set environment to dev
ENV="dev"

# Set compose file
COMPOSE_FILE="$PROJECT_ROOT/deploy/compose/docker-compose.${ENV}.yml"

echo -e "${GREEN}🚀 Starting LexiGo Development Environment${NC}"
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

# Start services
echo -e "${YELLOW}📦 Starting services for development environment...${NC}"

$DOCKER_COMPOSE -f "$COMPOSE_FILE" up --build -d

echo ""
echo -e "${GREEN}✅ Development environment started successfully!${NC}"
