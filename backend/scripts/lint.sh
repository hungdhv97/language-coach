#!/bin/bash

# Lint Script
# This script runs linting and formatting checks

set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

cd "$BACKEND_DIR"

echo "Running gofmt..."
gofmt -s -w .

echo "Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
	golangci-lint run
else
	echo "Warning: golangci-lint is not installed. Skipping..."
fi

echo "Linting completed!"

