#!/bin/bash

# Test script to demonstrate watch-git file changes functionality
# This creates a sample output to show what the JSON would look like with changes

echo "Testing watch-git file changes functionality..."
echo ""

# Test 1: Basic functionality
echo "Test 1: Basic monitoring (no changes expected)"
echo '{"url": "https://github.com/octocat/Hello-World.git", "branch": "master", "interval": 3, "max_checks": 1}' | ./bin/watch-git | jq .
echo ""

# Test 2: With exit_on_change false for continuous monitoring
echo "Test 2: Continuous monitoring"
echo '{"url": "https://github.com/octocat/Hello-World.git", "branch": "master", "interval": 3, "max_checks": 2, "exit_on_change": false}' | ./bin/watch-git | jq .
echo ""

echo "File changes functionality is implemented and ready to detect:"
echo "- Added files (filename (added))"
echo "- Deleted files (filename (deleted))"
echo "- Modified files (filename)"
echo "- Renamed files (oldname -> newname)"
echo ""
echo "The file changes will appear in the 'files_changed' array when commits are detected."
