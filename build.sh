#!/bin/bash

# Build script for Go AI Agent V1.0
# This script builds all components of the V1 agent

set -e

echo "Building Go AI Agent V1.0..."

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"

# Create bin directory if it doesn't exist
mkdir -p "$PROJECT_ROOT/bin"

echo "Building orchestrator..."
cd "$PROJECT_ROOT/orchestrator"
go mod tidy
go build -o "../bin/orchestrator" .

echo "Building CLI..."
cd "$PROJECT_ROOT/cli"
go mod tidy
go build -o "../bin/cli" .

echo "Building action modules..."

# Build echo-json action
echo "  - echo-json"
cd "$PROJECT_ROOT/actions/echo-json"
go mod tidy
go build -o "../../bin/echo-json" .

# Build writefile-json action
echo "  - writefile-json"
cd "$PROJECT_ROOT/actions/writefile-json"
go mod tidy
go build -o "../../bin/writefile-json" .

# Build httprequest action
echo "  - httprequest"
cd "$PROJECT_ROOT/actions/httprequest"
go mod tidy
go build -o "../../bin/httprequest" .

# Build claude-api action
echo "  - claude-api"
cd "$PROJECT_ROOT/actions/claude-api"
go mod tidy
go build -o "../../bin/claude-api" .

# Build watch-git action
echo "  - watch-git"
cd "$PROJECT_ROOT/actions/watch-git"
go mod tidy
go build -o "../../bin/watch-git" .

echo ""
echo "✅ Build completed successfully!"
echo ""
echo "Built components:"
echo "  - bin/orchestrator    - Main workflow orchestrator"
echo "  - bin/cli            - Command-line interface"
echo "  - bin/echo-json      - Echo action with JSON I/O"
echo "  - bin/writefile-json - File writing action"
echo "  - bin/httprequest    - HTTP request action"
echo "  - bin/claude-api     - Claude API integration action"
echo "  - bin/watch-git      - Git repository watcher action"
echo ""
echo "To test the V1 agent:"
echo "  ./bin/cli run workflows/simple-test.yaml '{\"agent_name\":\"Go-AI-Agent\",\"test_id\":\"123\"}'"
echo ""
