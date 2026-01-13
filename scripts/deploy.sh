#!/bin/bash

# Deployment Script with Rollback Support
# This script handles deployment with automatic rollback on failure
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

# Get version from first argument (required - passed from workflow)
DEPLOY_VERSION="${1:-}"

# Validate version is provided
if [ -z "$DEPLOY_VERSION" ]; then
    echo -e "${RED}❌ Error: Version is required${NC}"
    echo "Usage: $0 <version>"
    echo "Version should be passed from GitHub Actions workflow"
    exit 1
fi

# Validate version format
if [[ ! "$DEPLOY_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ Error: Invalid version format: $DEPLOY_VERSION${NC}"
    echo "Expected format: x.y.z (e.g., 1.0.0)"
    exit 1
fi

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
BACKUP_TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_PATH="$BACKUP_DIR/$ENV-$BACKUP_TIMESTAMP"

# Compose file
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

# Trap to handle errors
cleanup_on_error() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        log_error "Deployment failed with exit code $exit_code"
        log_error "Attempting rollback..."
        
        # Try to rollback
        if [ -d "$BACKUP_PATH" ]; then
            log_info "Rolling back to previous version..."
            bash "$SCRIPT_DIR/rollback.sh" "$BACKUP_TIMESTAMP" || {
                log_error "Rollback failed! Manual intervention required."
            }
        else
            log_error "No backup found for rollback!"
        fi
        
        exit $exit_code
    fi
}

trap cleanup_on_error ERR

log_info "Starting deployment for $ENV environment"
log_info "Deployment version: $DEPLOY_VERSION"

# Check if compose file exists
if [ ! -f "$COMPOSE_FILE" ]; then
    log_error "Docker compose file not found at $COMPOSE_FILE"
    exit 1
fi

# Step 1: Create backup of current deployment
log_info "Creating backup of current deployment..."
mkdir -p "$BACKUP_DIR"

# Backup docker-compose file (needed for rollback)
if [ -f "$COMPOSE_FILE" ]; then
    mkdir -p "$BACKUP_PATH"
    cp "$COMPOSE_FILE" "$BACKUP_PATH/docker-compose.${ENV}.yml"
    log_success "Backed up docker-compose file"
    
    # Save version info to backup (for reference)
    echo "Version: $DEPLOY_VERSION" > "$BACKUP_PATH/version.txt"
    echo "Environment: $ENV" >> "$BACKUP_PATH/version.txt"
    echo "Timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$BACKUP_PATH/version.txt"
fi

# Step 2: Build new images
log_info "Building new Docker images with version: $DEPLOY_VERSION..."
cd "$PROJECT_ROOT"

# Build images
$DOCKER_COMPOSE -f "$COMPOSE_FILE" build --no-cache 2>&1 | tee -a "$LOG_FILE" || {
    log_error "Failed to build Docker images"
    exit 1
}

log_success "Docker images built successfully"

# Step 3: Stop current containers (gracefully)
log_info "Stopping current containers..."
if $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps -q | grep -q .; then
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" stop 2>&1 | tee -a "$LOG_FILE" || {
        log_warning "Some containers may not have stopped gracefully"
    }
    log_success "Containers stopped"
else
    log_info "No running containers to stop"
fi

# Step 4: Start new containers
log_info "Starting new containers..."
$DOCKER_COMPOSE -f "$COMPOSE_FILE" up -d 2>&1 | tee -a "$LOG_FILE" || {
    log_error "Failed to start containers"
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
    log_error "Services did not become healthy within timeout period"
    
    # Show container status
    log_info "Container status:"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps 2>&1 | tee -a "$LOG_FILE"
    
    # Show logs for failed containers
    log_info "Recent logs from containers:"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs --tail=50 2>&1 | tee -a "$LOG_FILE"
    
    exit 1
fi

log_success "All services are healthy"

# Step 6: Cleanup old backups (keep last 5)
log_info "Cleaning up old backups..."
cd "$BACKUP_DIR"
ls -t | grep "^${ENV}-" | tail -n +6 | xargs -r rm -rf
log_success "Old backups cleaned up"

# Step 7: Deployment successful
log_success "Deployment completed successfully!"
log_info "Deployment version: $DEPLOY_VERSION"
log_info "Backup saved at: $BACKUP_PATH"

# Show container status
log_info "Final container status:"
$DOCKER_COMPOSE -f "$COMPOSE_FILE" ps

echo ""
echo -e "${GREEN}🎉 Deployment to $ENV completed successfully!${NC}"
echo -e "${GREEN}📦 Version: $DEPLOY_VERSION${NC}"
