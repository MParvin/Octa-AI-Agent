# Migration Guide: JSON to YAML

This document outlines the migration from JSON to YAML format in the Octa-AI-Agent system.

## Overview

The Octa-AI-Agent system has been fully migrated from JSON to YAML format for improved readability, maintainability, and human-friendliness. This migration affects:

- **Workflow definitions** (examples/ and workflows/ directories)
- **Action input/output** (all actions in actions/ directory)
- **CLI and orchestrator** (parsing and marshaling)
- **Documentation and examples**

## What Changed

### 1. File Extensions
- All `.json` files have been converted to `.yaml` files
- Original JSON files are preserved for backward compatibility

### 2. Action Input/Output
All actions now use YAML for input parsing and output formatting:

#### Before (JSON):
```go
type Input struct {
    Message string `json:"message"`
    Count   int    `json:"count"`
}
```

#### After (YAML):
```go
type Input struct {
    Message string `yaml:"message"`
    Count   int    `yaml:"count"`
}
```

### 3. Workflow Definitions
Workflow files are now in YAML format, making them more readable:

#### Before (hello-world.json):
```json
{
  "name": "Hello World Workflow",
  "description": "Simple greeting workflow with templating",
  "nodes": {
    "greeting": {
      "action": "echo-json",
      "input": {
        "message": "Hello, {{.WorkflowData.name}}! Welcome to Go AI Agent V1.0"
      }
    }
  }
}
```

#### After (hello-world.yaml):
```yaml
name: "Hello World Workflow"
description: "Simple greeting workflow with templating"
nodes:
  greeting:
    action: "echo-json"
    input:
      message: "Hello, {{.WorkflowData.name}}! Welcome to Go AI Agent V1.0"
```

### 4. CLI Commands
Commands remain the same, but now reference YAML files:

#### Before:
```bash
./bin/cli run examples/hello-world.json '{"name":"Alice"}'
```

#### After:
```bash
./bin/cli run examples/hello-world.yaml '{"name":"Alice"}'
```

## Migration Benefits

### 1. Improved Readability
- No curly braces and commas cluttering the content
- Natural indentation structure
- Comments support (using `#`)

### 2. Better Multi-line Support
YAML's natural multi-line string support makes complex content easier to manage:

```yaml
# Multi-line strings without escaping
content: |
  This is a multi-line string
  that preserves line breaks
  and doesn't require escaping quotes.
```

### 3. Reduced Syntax Errors
- No missing commas or brackets
- More forgiving format
- Better error messages

### 4. Enhanced Developer Experience
- Easier to read and modify workflows
- Better integration with configuration management tools
- Improved version control diffs

## Backward Compatibility

The system maintains backward compatibility:
- JSON files still exist alongside YAML files
- Both formats are supported by the CLI and orchestrator
- Existing JSON workflows continue to work

## File Mapping

| Original JSON File | New YAML File |
|-------------------|---------------|
| `examples/hello-world.json` | `examples/hello-world.yaml` |
| `examples/file-operations.json` | `examples/file-operations.yaml` |
| `examples/api-integration.json` | `examples/api-integration.yaml` |
| `examples/data-pipeline.json` | `examples/data-pipeline.yaml` |
| `examples/cicd-pipeline.json` | `examples/cicd-pipeline.yaml` |
| `examples/git-watch.json` | `examples/git-watch.yaml` |
| `examples/claude-ai-integration.json` | `examples/claude-ai-integration.yaml` |
| `workflows/simple-test.json` | `workflows/simple-test.yaml` |
| `workflows/v1-demo.json` | `workflows/v1-demo.yaml` |
| `workflows/error-test.json` | `workflows/error-test.yaml` |
| `workflows/v1-simple-demo.json` | `workflows/v1-simple-demo.yaml` |

## Updated Commands

### Validation
```bash
# Before
./bin/cli validate examples/hello-world.json

# After  
./bin/cli validate examples/hello-world.yaml
```

### Execution
```bash
# Before
./bin/cli run examples/file-operations.json '{"user_id":"12345","timestamp":"2024-01-20T10:30:00Z"}'

# After
./bin/cli run examples/file-operations.yaml '{"user_id":"12345","timestamp":"2024-01-20T10:30:00Z"}'
```

### Running Examples
Use the updated `run-examples.sh` script which now uses YAML files:

```bash
./run-examples.sh
```

## Best Practices for YAML Workflows

### 1. Use Consistent Indentation
Always use 2 spaces for indentation:

```yaml
nodes:
  step1:
    action: "echo-json"
    input:
      message: "Hello World"
```

### 2. Quote Template Variables
Always quote template variables to prevent parsing issues:

```yaml
message: "Hello, {{.WorkflowData.name}}!"
path: "{{.WorkflowData.output_dir}}/file.txt"
```

### 3. Use Multi-line Strings for Complex Content
```yaml
content: |
  Line 1
  Line 2
  Line 3
```

### 4. Add Comments for Clarity
```yaml
# This node fetches data from the API
fetch_data:
  action: "httprequest"
  input:
    url: "https://api.example.com/data"
    method: "GET"
```

## Troubleshooting

### Common Issues

1. **Indentation Errors**: YAML is indent-sensitive. Use 2 spaces consistently.

2. **Quote Issues**: Always quote strings containing special characters or template variables.

3. **Multi-line Strings**: Use `|` for literal block scalars or `>` for folded block scalars.

### Validation
Always validate your YAML workflows before running:

```bash
./bin/cli validate path/to/workflow.yaml
```

## Next Steps

1. **Update your workflows**: Convert any custom JSON workflows to YAML format
2. **Update documentation**: Reference YAML files in your documentation
3. **Update scripts**: Modify any automation scripts to use YAML files
4. **Remove JSON files**: Once confident in the migration, consider removing JSON files

For more information, see:
- `YAML-COMPARISON.md` - Detailed comparison between JSON and YAML formats
- `examples/README.md` - Updated examples documentation
- `README.md` - Main project documentation
