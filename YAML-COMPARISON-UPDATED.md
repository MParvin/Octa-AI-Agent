# YAML vs JSON Comparison

This document demonstrates the benefits of using YAML over JSON for the Go AI Agent workflow system.

## Key Benefits of YAML

1. **Human Readability**: YAML is more readable and easier to edit manually
2. **Multi-line Support**: Natural support for multi-line strings without escaping
3. **No Quote Requirements**: Simpler syntax without requiring quotes for most strings
4. **Comments Support**: Ability to add comments for documentation
5. **Nested Structure**: More intuitive representation of nested data

## Side-by-Side Comparison

### Hello World Workflow

**JSON Version:**
```json
{
  "name": "Hello World Workflow",
  "description": "Simple greeting workflow with templating",
  "version": "1.0",
  "nodes": [
    {
      "id": "greeting",
      "type": "echo-json",
      "inputs_from_workflow": {
        "message": "Hello, {{.WorkflowData.name}}! Welcome to Go AI Agent V1.0"
      }
    }
  ]
}
```

**YAML Version:**
```yaml
name: "Hello World Workflow"
description: "Simple greeting workflow with templating"
version: "1.0"
nodes:
  - id: "greeting"
    type: "echo-json"
    inputs_from_workflow:
      message: "Hello, {{.WorkflowData.name}}! Welcome to Go AI Agent V1.0"
```

### File Operations with Multi-line Content

**JSON Version:**
```json
{
  "name": "File Operations Workflow",
  "description": "Demonstrates file creation and content templating",
  "version": "1.0",
  "nodes": [
    {
      "id": "create_user_report",
      "type": "writefile-json",
      "inputs_from_workflow": {
        "path": "/tmp/user_{{.WorkflowData.user_id}}_report.txt",
        "mode": "create",
        "content": "User Report\\n=============\\n\\nUser ID: {{.WorkflowData.user_id}}\\nGenerated: {{.WorkflowData.timestamp}}\\n\\nReport complete."
      }
    }
  ]
}
```

**YAML Version:**
```yaml
name: "File Operations Workflow"
description: "Demonstrates file creation and content templating"
version: "1.0"
nodes:
  - id: "create_user_report"
    type: "writefile-json"
    inputs_from_workflow:
      path: "/tmp/user_{{.WorkflowData.user_id}}_report.txt"
      mode: "create"
      content: |
        User Report
        =============
        
        User ID: {{.WorkflowData.user_id}}
        Generated: {{.WorkflowData.timestamp}}
        
        Report complete.
```

### API Integration with Headers

**JSON Version:**
```json
{
  "name": "API Integration Workflow",
  "nodes": [
    {
      "id": "fetch_data",
      "type": "httprequest",
      "inputs_from_workflow": {
        "url": "https://api.example.com/users/{{.WorkflowData.user_id}}",
        "method": "GET",
        "headers": {
          "Accept": "application/json",
          "User-Agent": "Go-AI-Agent-V1",
          "Authorization": "Bearer {{.WorkflowData.token}}"
        },
        "timeout": 10
      }
    }
  ]
}
```

**YAML Version:**
```yaml
name: "API Integration Workflow"
nodes:
  - id: "fetch_data"
    type: "httprequest"
    inputs_from_workflow:
      url: "https://api.example.com/users/{{.WorkflowData.user_id}}"
      method: "GET"
      headers:
        Accept: "application/json"
        User-Agent: "Go-AI-Agent-V1"
        Authorization: "Bearer {{.WorkflowData.token}}"
      timeout: 10
```

## Action Input/Output Examples

### Echo Action

**Input (YAML):**
```yaml
message: "Hello, World!"
```

**Output (YAML):**
```yaml
echoed_message: "Hello, World!"
original_input:
  message: "Hello, World!"
```

### File Writing Action

**Input (YAML):**
```yaml
path: "/tmp/example.txt"
content: |
  This is a multi-line file content
  with proper line breaks and formatting.
  
  Template variables work here: {{.WorkflowData.user_name}}
mode: "create"
mkdir_all: true
```

**Output (YAML):**
```yaml
success: true
message: "Successfully wrote 134 bytes to /tmp/example.txt"
path: "/tmp/example.txt"
size: 134
```

### HTTP Request Action

**Input (YAML):**
```yaml
url: "https://httpbin.org/post"
method: "POST"
headers:
  Content-Type: "application/json"
  Authorization: "Bearer token123"
body: |
  {
    "user_id": "{{.WorkflowData.user_id}}",
    "timestamp": "{{.WorkflowData.timestamp}}"
  }
timeout: 30
```

**Output (YAML):**
```yaml
success: true
message: "HTTP POST request to https://httpbin.org/post completed successfully"
status_code: 200
headers:
  Content-Type: "application/json"
  Server: "gunicorn/19.9.0"
body: |
  {
    "json": {
      "user_id": "123",
      "timestamp": "2024-01-20T10:30:00Z"
    }
  }
```

## Migration Benefits

1. **Easier Configuration Management**: YAML workflows are much easier to read and modify
2. **Better Version Control**: YAML diffs are more readable in Git
3. **Reduced Syntax Errors**: Less prone to JSON syntax errors (missing commas, quotes)
4. **Natural Multi-line Support**: No need to escape newlines or manage JSON strings
5. **Self-Documenting**: The structure itself serves as documentation

## Backward Compatibility

The system maintains backward compatibility with existing JSON workflows while encouraging migration to YAML for new workflows. Both formats are supported simultaneously.
