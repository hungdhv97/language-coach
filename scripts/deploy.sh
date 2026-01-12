#!/bin/bash

# Deployment Script with Rollback Support
# This script handles deployment with automatic rollback on failure

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

# Get environment from first argument (default to prod)
ENV="${1:-prod}"

# Get version from second argument (optional)
DEPLOY_VERSION="${2:-}"

# Validate environment - only prod is supported for deployment
if [ "$ENV" != "prod" ]; then
    echo -e "${RED}❌ Error: Deployment only supports 'prod' environment${NC}"
    exit 1
fi

# Function to bump version based on commit message
bump_version() {
    local current_version=$1
    local commit_message=$2
    
    # Parse version components
    IFS='.' read -r -a parts <<< "$current_version"
    local major=${parts[0]:-0}
    local minor=${parts[1]:-0}
    local patch=${parts[2]:-0}
    
    # Convert commit message to lowercase for case-insensitive matching
    local msg_lower=$(echo "$commit_message" | tr '[:upper:]' '[:lower:]')
    
    # Determine bump type from commit message
    if echo "$msg_lower" | grep -qi "major"; then
        major=$((major + 1))
        minor=0
        patch=0
        echo "$major.$minor.$patch"
    elif echo "$msg_lower" | grep -qi "minor"; then
        minor=$((minor + 1))
        patch=0
        echo "$major.$minor.$patch"
    elif echo "$msg_lower" | grep -qi "patch"; then
        patch=$((patch + 1))
        echo "$major.$minor.$patch"
    else
        # Default to patch if no keyword found
        patch=$((patch + 1))
        echo "$major.$minor.$patch"
    fi
}

# Get version from deploy/VERSION file and bump if needed
VERSION_FILE="$PROJECT_ROOT/deploy/VERSION"
if [ -f "$VERSION_FILE" ]; then
    CURRENT_VERSION=$(cat "$VERSION_FILE" | tr -d ' \n')
    # Validate version format (x.y.z)
    if [[ ! "$CURRENT_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo -e "${RED}❌ Error: Invalid version format in deploy/VERSION: $CURRENT_VERSION${NC}"
        echo "Expected format: x.y.z (e.g., 1.0.0)"
        exit 1
    fi
    
    # If version not provided, bump based on commit message
    if [ -z "$DEPLOY_VERSION" ]; then
        # Get commit message (last commit)
        if [ -d "$PROJECT_ROOT/.git" ]; then
            COMMIT_MESSAGE=$(git log -1 --pretty=%B 2>/dev/null || echo "")
            if [ -n "$COMMIT_MESSAGE" ]; then
                DEPLOY_VERSION=$(bump_version "$CURRENT_VERSION" "$COMMIT_MESSAGE")
                # Update VERSION file
                echo "$DEPLOY_VERSION" > "$VERSION_FILE"
            else
                # Fallback: bump patch if no commit message
                DEPLOY_VERSION=$(bump_version "$CURRENT_VERSION" "patch")
                echo "$DEPLOY_VERSION" > "$VERSION_FILE"
            fi
        else
            # Fallback: bump patch if not in git repo
            DEPLOY_VERSION=$(bump_version "$CURRENT_VERSION" "patch")
            echo "$DEPLOY_VERSION" > "$VERSION_FILE"
        fi
    fi
else
    if [ -z "$DEPLOY_VERSION" ]; then
        echo -e "${RED}❌ Error: deploy/VERSION file not found${NC}"
        echo "Please create deploy/VERSION file with semantic version (x.y.z)"
        exit 1
    fi
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
            bash "$SCRIPT_DIR/rollback.sh" "$ENV" "$BACKUP_TIMESTAMP" || {
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

# Backup docker-compose file
if [ -f "$COMPOSE_FILE" ]; then
    mkdir -p "$BACKUP_PATH"
    cp "$COMPOSE_FILE" "$BACKUP_PATH/docker-compose.${ENV}.yml"
    log_success "Backed up docker-compose file"
fi

# Save version info to backup
if [ -d "$BACKUP_PATH" ]; then
    echo "Version: $DEPLOY_VERSION" > "$BACKUP_PATH/version.txt"
    echo "Environment: $ENV" >> "$BACKUP_PATH/version.txt"
    echo "Timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$BACKUP_PATH/version.txt"
    if [ -d "$PROJECT_ROOT/.git" ]; then
        echo "Commit: $(git rev-parse HEAD 2>/dev/null || echo 'unknown')" >> "$BACKUP_PATH/version.txt"
        echo "Branch: $(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')" >> "$BACKUP_PATH/version.txt"
    fi
fi

# Backup current images (if containers are running)
log_info "Backing up current container images..."
if $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps -q | grep -q .; then
    # Save current image tags
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" config --images > "$BACKUP_PATH/previous-images.txt" 2>/dev/null || true
    log_success "Backed up current image information"
fi

# Step 2: Pull latest code (if in git repo)
if [ -d "$PROJECT_ROOT/.git" ]; then
    log_info "Pulling latest code..."
    cd "$PROJECT_ROOT"
    git pull origin main || log_warning "Could not pull latest code (may not be in git repo)"
fi

# Step 3: Build new images
log_info "Building new Docker images with version: $DEPLOY_VERSION..."
cd "$PROJECT_ROOT"

# Build images
$DOCKER_COMPOSE -f "$COMPOSE_FILE" build --no-cache 2>&1 | tee -a "$LOG_FILE" || {
    log_error "Failed to build Docker images"
    rm -f "$TEMP_COMPOSE"
    exit 1
}

# Tag images with version after build
log_info "Tagging images with version: $DEPLOY_VERSION..."
SERVICES=$($DOCKER_COMPOSE -f "$COMPOSE_FILE" config --services)
for SERVICE in $SERVICES; do
    # Skip postgres and redis (external images)
    if [ "$SERVICE" = "postgres" ] || [ "$SERVICE" = "redis" ]; then
        continue
    fi
    
    # Get the image that was just built
    IMAGE_ID=$($DOCKER_COMPOSE -f "$COMPOSE_FILE" images -q "$SERVICE" 2>/dev/null | head -n 1)
    if [ -n "$IMAGE_ID" ]; then
        # Create versioned image name
        VERSIONED_IMAGE="lexigo-${SERVICE}:${DEPLOY_VERSION}"
        docker tag "$IMAGE_ID" "$VERSIONED_IMAGE" 2>/dev/null && {
            log_success "Tagged $SERVICE as $VERSIONED_IMAGE"
        } || log_warning "Failed to tag $SERVICE"
    fi
done

log_success "Docker images built and tagged successfully"

# Step 5: Stop current containers (gracefully)
log_info "Stopping current containers..."
if $DOCKER_COMPOSE -f "$COMPOSE_FILE" ps -q | grep -q .; then
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" stop 2>&1 | tee -a "$LOG_FILE" || {
        log_warning "Some containers may not have stopped gracefully"
    }
    log_success "Containers stopped"
else
    log_info "No running containers to stop"
fi

# Step 6: Start new containers
log_info "Starting new containers..."
$DOCKER_COMPOSE -f "$COMPOSE_FILE" up -d 2>&1 | tee -a "$LOG_FILE" || {
    log_error "Failed to start containers"
    exit 1
}
log_success "Containers started"

# Step 7: Wait for services to be healthy
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

# Step 8: Cleanup old backups (keep last 5)
log_info "Cleaning up old backups..."
cd "$BACKUP_DIR"
ls -t | grep "^${ENV}-" | tail -n +6 | xargs -r rm -rf
log_success "Old backups cleaned up"

# Step 9: Save deployment version info
VERSION_FILE="$PROJECT_ROOT/deploy/versions/current-${ENV}.txt"
mkdir -p "$(dirname "$VERSION_FILE")"
echo "Version: $DEPLOY_VERSION" > "$VERSION_FILE"
echo "Environment: $ENV" >> "$VERSION_FILE"
echo "Timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$VERSION_FILE"
if [ -d "$PROJECT_ROOT/.git" ]; then
    echo "Commit: $(git rev-parse HEAD 2>/dev/null || echo 'unknown')" >> "$VERSION_FILE"
    echo "Branch: $(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')" >> "$VERSION_FILE"
fi

# Step 10: Deployment successful
log_success "Deployment completed successfully!"
log_info "Deployment version: $DEPLOY_VERSION"
log_info "Backup saved at: $BACKUP_PATH"

# Show container status
log_info "Final container status:"
$DOCKER_COMPOSE -f "$COMPOSE_FILE" ps

echo ""
echo -e "${GREEN}🎉 Deployment to $ENV completed successfully!${NC}"
echo -e "${GREEN}📦 Version: $DEPLOY_VERSION${NC}"
