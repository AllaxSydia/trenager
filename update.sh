#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Log function
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Main update function
update_project() {
    local PROJECT_DIR="/home/user1/trenager"
    
    log "Starting project update..."
    
    # Check if project directory exists
    if [ ! -d "$PROJECT_DIR" ]; then
        error "Project directory not found: $PROJECT_DIR"
        return 1
    fi
    
    cd "$PROJECT_DIR" || {
        error "Failed to enter project directory"
        return 1
    }
    
    # Stop current services
    log "Stopping current services..."
    if ! docker-compose down; then
        error "Failed to stop services"
        return 1
    fi
    
    # Backup current version (optional)
    log "Backing up current version..."
    if [ -d ".git" ]; then
        git tag "backup-$(date +'%Y%m%d-%H%M%S')"
    fi
    
    # Pull latest changes
    log "Pulling latest changes from GitHub..."
    if ! git fetch origin; then
        error "Failed to fetch from GitHub"
        return 1
    fi
    
    if ! git pull origin main; then
        error "Failed to pull changes from GitHub"
        return 1
    fi
    
    # Rebuild and start services
    log "Rebuilding and starting services..."
    if ! docker-compose up -d --build; then
        error "Failed to rebuild and start services"
        return 1
    fi
    
    # Wait for services to start
    log "Waiting for services to initialize..."
    sleep 10
    
    # Health checks
    log "Performing health checks..."
    
    # Check backend
    if curl -f http://localhost:3000/health > /dev/null 2>&1; then
        success "Backend is healthy"
    else
        error "Backend health check failed"
    fi
    
    # Check frontend
    if curl -f http://localhost:3001 > /dev/null 2>&1; then
        success "Frontend is responding"
    else
        error "Frontend health check failed"
    fi
    
    # Show container status
    log "Current container status:"
    docker-compose ps
    
    # Show recent logs
    log "Recent logs (last 10 lines):"
    docker-compose logs --tail=10
    
    success "Project update completed successfully!"
    log "Frontend: http://localhost:3001"
    log "Backend API: http://localhost:3000/health"
}

# Rollback function
rollback() {
    log "Starting rollback..."
    
    cd /home/user1/trenager || return 1
    
    # Stop services
    docker-compose down
    
    # Get last backup tag
    local LAST_TAG
    LAST_TAG=$(git tag --list "backup-*" | sort -r | head -n1)
    
    if [ -n "$LAST_TAG" ]; then
        log "Rolling back to: $LAST_TAG"
        git checkout "$LAST_TAG"
        docker-compose up -d --build
        success "Rollback to $LAST_TAG completed"
    else
        error "No backup tags found for rollback"
    fi
}

# Status function
status() {
    log "Project status:"
    cd /home/user1/trenager && docker-compose ps
    echo ""
    log "Recent logs:"
    docker-compose logs --tail=5
}

# Help function
show_help() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  update    - Update project to latest version (default)"
    echo "  rollback  - Rollback to last backup"
    echo "  status    - Show project status"
    echo "  help      - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 update    # Update project"
    echo "  $0 status    # Check status"
    echo "  $0 rollback  # Rollback to previous version"
}

# Main script logic
case "${1:-update}" in
    "update")
        update_project
        ;;
    "rollback")
        rollback
        ;;
    "status")
        status
        ;;
    "help"|"-h"|"--help")
        show_help
        ;;
    *)
        error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac