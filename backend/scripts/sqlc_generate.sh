#!/bin/bash

# SQLC Generate Script
# This script generates sqlc code from db/queries and db/schema

set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

cd "$BACKEND_DIR"

# Check if sqlc is installed
if ! command -v sqlc &> /dev/null; then
	echo "Error: sqlc is not installed. Please install it first:"
	echo "  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
	exit 1
fi

echo "Generating sqlc code..."
sqlc generate

echo "SQLC code generation completed successfully!"

