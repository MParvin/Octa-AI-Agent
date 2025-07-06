#!/bin/bash

# Test script to demonstrate watch-git YAML functionality
# This shows the human-readable YAML format for both input and output

echo "=== Watch-Git YAML Functionality Demo ==="
echo ""

echo "1. Creating human-readable YAML input file:"
cat > demo-input.yaml << 'EOF'
url: "https://github.com/octocat/Hello-World.git"
branch: "master"
interval: 3
max_checks: 2
exit_on_change: true
EOF

echo "Input (YAML):"
cat demo-input.yaml
echo ""

echo "2. Running watch-git with YAML input:"
echo "$ cat demo-input.yaml | ./bin/watch-git"
echo ""

cat demo-input.yaml | ./bin/watch-git
echo ""

echo "3. Testing with authentication example:"
cat > demo-auth.yaml << 'EOF'
url: "https://github.com/golang/example.git"
branch: "master"
username: ""
password: ""
interval: 2
max_checks: 1
exit_on_change: false
EOF

echo "Input with authentication (YAML):"
cat demo-auth.yaml
echo ""

echo "Output (YAML):"
cat demo-auth.yaml | ./bin/watch-git

# Cleanup
rm -f demo-input.yaml demo-auth.yaml

echo ""
echo "=== YAML Benefits ==="
echo "✅ Human-readable format"
echo "✅ Easy to edit and maintain"
echo "✅ Natural structure for configuration"
echo "✅ No escaping quotes needed"
echo "✅ Better for complex data structures"
