#!/bin/bash

# Go AI Agent V1.0 - Examples Test Script
# This script demonstrates how to run all the examples

echo "==========================================="
echo "Go AI Agent V1.0 - Examples Test Runner"
echo "==========================================="
echo ""

# Build the project first
echo "🔨 Building Go AI Agent V1.0..."
./build.sh
echo ""

# Create necessary directories
echo "📁 Creating test directories..."
mkdir -p /tmp/{reports,data-pipeline,cicd/{artifacts,reports}}
echo ""

echo "🚀 Running Examples..."
echo ""

# Example 1: Hello World
echo "1️⃣  Running Hello World Example..."
echo "   Command: ./bin/cli run examples/hello-world.json '{\"name\":\"Alice\"}'"
./bin/cli run examples/hello-world.json '{"name":"Alice"}'
echo ""

# Example 2: File Operations
echo "2️⃣  Running File Operations Example..."
echo "   Command: ./bin/cli run examples/file-operations.json '{\"user_id\":\"12345\",\"timestamp\":\"2024-01-20T10:30:00Z\"}'"
./bin/cli run examples/file-operations.json '{"user_id":"12345","timestamp":"2024-01-20T10:30:00Z"}'
echo ""

# Example 3: API Integration
echo "3️⃣  Running API Integration Example..."
echo "   Command: ./bin/cli run examples/api-integration.json '{\"user_id\":\"3\",\"timestamp\":\"2024-01-20T10:30:00Z\",\"output_dir\":\"/tmp/reports\"}'"
./bin/cli run examples/api-integration.json '{"user_id":"3","timestamp":"2024-01-20T10:30:00Z","output_dir":"/tmp/reports"}'
echo ""

# Example 4: Data Pipeline
echo "4️⃣  Running Data Pipeline Example..."
echo "   Command: ./bin/cli run examples/data-pipeline.json '{\"batch_id\":\"batch_001\",\"timestamp\":\"2024-01-20T10:30:00Z\",\"output_dir\":\"/tmp/data-pipeline\",\"record_count\":\"100\"}'"
./bin/cli run examples/data-pipeline.json '{"batch_id":"batch_001","timestamp":"2024-01-20T10:30:00Z","output_dir":"/tmp/data-pipeline","record_count":"100"}'
echo ""

# Example 5: CI/CD Pipeline
echo "5️⃣  Running CI/CD Pipeline Example..."
echo "   Command: ./bin/cli run examples/cicd-pipeline.json '...(complex JSON)...'"
./bin/cli run examples/cicd-pipeline.json '{
  "project_name": "web-application",
  "version": "2.1.0",
  "branch": "main",
  "commit_sha": "abc123def456",
  "environment": "staging",
  "pipeline_id": "pipeline_001",
  "timestamp": "2024-01-20T10:30:00Z",
  "artifacts_dir": "/tmp/cicd/artifacts",
  "reports_dir": "/tmp/cicd/reports",
  "api_token": "sk-test-token",
  "deploy_token": "deploy-token-123"
}'
echo ""

echo "✅ All examples completed!"
echo ""
echo "📁 Generated files can be found in:"
echo "   - /tmp/user_12345_report.txt"
echo "   - /tmp/reports/"
echo "   - /tmp/data-pipeline/"
echo "   - /tmp/cicd/"
echo ""
echo "🔍 To validate any workflow, use:"
echo "   ./bin/cli validate examples/[workflow-name].json"
echo ""
