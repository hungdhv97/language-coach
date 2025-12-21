#!/bin/bash

# Migration Script
# This script runs database migrations (schema + data) using Go

set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

cd "$BACKEND_DIR"

# Run schema migration
echo "Running schema migration..."
go run cmd/migration/schema/main.go

# Run data migration (if needed)
if [ "$1" != "--schema-only" ]; then
	echo "Running data migration..."
	go run cmd/migration/data/main.go
fi

echo "Migration completed successfully!"

