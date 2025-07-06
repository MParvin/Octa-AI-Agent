#!/bin/bash

# Go AI Agent V1.0 - Examples Validation Script
# Quick validation of all example workflows

echo "🔍 Validating Go AI Agent V1.0 Examples..."
echo "=============================================="
echo ""

# Check if project is built
if [ ! -f "./bin/cli" ]; then
    echo "❌ Project not built. Please run ./build.sh first"
    exit 1
fi

echo "✅ Project is built"
echo ""

# Validate each example workflow
examples_dir="examples"
all_valid=true

echo "📝 Validating workflow syntax..."

for workflow in "$examples_dir"/*.json; do
    if [ -f "$workflow" ]; then
        filename=$(basename "$workflow")
        echo -n "   $filename ... "
        
        if ./bin/cli validate "$workflow" >/dev/null 2>&1; then
            echo "✅ Valid"
        else
            echo "❌ Invalid"
            all_valid=false
        fi
    fi
done

echo ""

if [ "$all_valid" = true ]; then
    echo "🎉 All examples are valid!"
    echo ""
    echo "🚀 Ready to run examples:"
    echo "   ./bin/cli run examples/hello-world.json '{\"name\":\"Alice\"}'"
    echo "   ./bin/cli run examples/file-operations.json '{\"user_id\":\"12345\",\"timestamp\":\"2025-06-07T17:30:00Z\"}'"
    echo "   mkdir -p /tmp/reports && ./bin/cli run examples/api-integration.json '{\"user_id\":\"3\",\"timestamp\":\"2025-06-07T17:30:00Z\",\"output_dir\":\"/tmp/reports\"}'"
    echo ""
else
    echo "❌ Some examples have validation errors"
    exit 1
fi
