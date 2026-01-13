#!/bin/bash

# Rollback Script
# This script rolls back to a previous deployment version
# NOTE: Only supports production environment

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

# Hardcode environment to production
ENV="prod"

# Get backup timestamp from first argument (optional - uses latest if not provided)
BACKUP_TIMESTAMP="${1:-}"

# Log file in deploy/logs/
LOG_DIR="$PROJECT_ROOT/deploy/logs"
LOG_FILE="$LOG_DIR/deploy.log"
mkdir -p "$LOG_DIR"

# Function to log messages
log() {
    local level=$1
    shift
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_FILE"
}

# Function to log errors
log_error() {
    log "ERROR" "$@"
    echo -e "${RED}❌ $@${NC}" >&2
}

# Function to log success
log_success() {
    log "SUCCESS" "$@"
    echo -e "${GREEN}✅ $@${NC}"
}

# Function to log info
log_info() {
    log "INFO" "$@"
    echo -e "${BLUE}ℹ️  $@${NC}"
}

# Function to log warning
log_warning() {
    log "WARNING" "$@"
    echo -e "${YELLOW}⚠️  $@${NC}"
}

# Backup directory
BACKUP_DIR="$PROJECT_ROOT/deploy/backups"
COMPOSE_FILE="$PROJECT_ROOT/deploy/compose/docker-compose.${ENV}.yml"

# Use docker compose (v2) if available, otherwise fall back to docker-compose (v1)
if docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
elif command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
else
    log_error "docker compose is not available"
    exit 1
fi

log_info "Starting rollback for $ENV environment"

# Find backup to restore
if [ -z "$BACKUP_TIMESTAMP" ]; then
    # Find latest backup
    LATEST_BACKUP=$(ls -t "$BACKUP_DIR" | grep "^${ENV}-" | head -n 1)
    if [ -z "$LATEST_BACKUP" ]; then
        log_error "No backup found for $ENV environment"
        exit 1
    fi
    BACKUP_PATH="$BACKUP_DIR/$LATEST_BACKUP"
    log_info "Using latest backup: $LATEST_BACKUP"
else
    BACKUP_PATH="$BACKUP_DIR/${ENV}-${BACKUP_TIMESTAMP}"
    if [ ! -d "$BACKUP_PATH" ]; then
        log_error "Backup not found: $BACKUP_PATH"
        log_info "Available backups:"
        ls -1 "$BACKUP_DIR" | grep "^${ENV}-" || echo "  (none)"
        exit 1
    fi
    log_info "Using backup: ${ENV}-${BACKUP_TIMESTAMP}"
fi

# Read version from backup if available
ROLLBACK_VERSION="unknown"
if [ -f "$BACKUP_PATH/version.txt" ]; then
    ROLLBACK_VERSION=$(grep "^Version:" "$BACKUP_PATH/version.txt" | cut -d' ' -f2- || echo "unknown")
    log_info "Rolling back to version: $ROLLBACK_VERSION"
fi

# Step 1: Stop current containers
log_info "Stopping current containers..."
if $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps -q | grep -q .; then
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" stop 2>&1 | tee -a "$LOG_FILE" || {
        log_warning "Some containers may not have stopped gracefully"
    }
    log_success "Containers stopped"
else
    log_info "No running containers to stop"
fi

# Step 2: Remove current containers
log_info "Removing current containers..."
$DOCKER_COMPOSE -f "$COMPOSE_FILE" down 2>&1 | tee -a "$LOG_FILE" || {
    log_warning "Some containers may not have been removed"
}

# Step 3: Restore docker-compose file if backed up
if [ -f "$BACKUP_PATH/docker-compose.${ENV}.yml" ]; then
    log_info "Restoring docker-compose file..."
    cp "$BACKUP_PATH/docker-compose.${ENV}.yml" "$COMPOSE_FILE"
    log_success "Docker-compose file restored"
fi

# Step 4: Rebuild and start containers
log_info "Rebuilding and starting containers..."
cd "$PROJECT_ROOT"
$DOCKER_COMPOSE -f "$COMPOSE_FILE" up -d --build 2>&1 | tee -a "$LOG_FILE" || {
    log_error "Failed to start containers after rollback"
    exit 1
}
log_success "Containers started"

# Step 5: Wait for services to be healthy
log_info "Waiting for services to be healthy..."
sleep 5

# Check if containers are running
MAX_RETRIES=30
RETRY_COUNT=0
ALL_HEALTHY=false

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    # Check if containers are running (using simple ps check)
    RUNNING_COUNT=$($DOCKER_COMPOSE -f "$COMPOSE_FILE" ps -q | wc -l)
    EXPECTED_COUNT=$($DOCKER_COMPOSE -f "$COMPOSE_FILE" config --services | wc -l)
    
    if [ "$RUNNING_COUNT" -ge "$EXPECTED_COUNT" ]; then
        # Check if any containers exited
        EXITED_COUNT=$($DOCKER_COMPOSE -f "$COMPOSE_FILE" ps | grep -c "Exit" || echo "0")
        
        if [ "$EXITED_COUNT" -eq 0 ]; then
            ALL_HEALTHY=true
            break
        fi
    fi
    
    RETRY_COUNT=$((RETRY_COUNT + 1))
    log_info "Waiting for services... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

if [ "$ALL_HEALTHY" = false ]; then
    log_error "Services did not become healthy after rollback"
    
    # Show container status
    log_info "Container status:"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps 2>&1 | tee -a "$LOG_FILE"
    
    # Show logs for failed containers
    log_info "Recent logs from containers:"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs --tail=50 2>&1 | tee -a "$LOG_FILE"
    
    exit 1
fi

log_success "All services are healthy after rollback"

# Step 6: Save rollback version info
VERSION_FILE="$PROJECT_ROOT/deploy/versions/current-${ENV}.txt"
mkdir -p "$(dirname "$VERSION_FILE")"
echo "Version: $ROLLBACK_VERSION (rolled back)" > "$VERSION_FILE"
echo "Environment: $ENV" >> "$VERSION_FILE"
echo "Timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$VERSION_FILE"
echo "Backup: $BACKUP_PATH" >> "$VERSION_FILE"

# Step 7: Rollback successful
log_success "Rollback completed successfully!"
log_info "Restored to version: $ROLLBACK_VERSION"
log_info "Restored from backup: $BACKUP_PATH"

# Show container status
log_info "Final container status:"
$DOCKER_COMPOSE -f "$COMPOSE_FILE" ps

echo ""
echo -e "${GREEN}🎉 Rollback to $ENV completed successfully!${NC}"
echo -e "${GREEN}📦 Version: $ROLLBACK_VERSION${NC}"
