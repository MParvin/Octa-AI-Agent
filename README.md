# Go AI Agent V1.0

A modular, extensible workflow orchestration engine built in Go with JSON-based inter-module communication, advanced templating, and robust error handling.

## 🚀 New in V1.0

- **JSON Communication Protocol**: All modules communicate via JSON stdin/stdout for better data handling
- **Go Template System**: Advanced templating using Go's `text/template` package with workflow data and node context
- **Initial Workflow Data**: Support for passing initial JSON data to workflows
- **Enhanced Error Handling**: Structured error responses with graceful degradation
- **New Action Modules**: `writefile-json`, `httprequest` with comprehensive features
- **CLI Interface**: Command-based interface with validation capabilities
- **Comprehensive Logging**: Detailed logging throughout the execution pipeline

## 📁 Project Structure

```
go-ai-agent-v1/
├── orchestrator/           # Main workflow orchestrator
│   ├── main.go
│   └── go.mod
├── cli/                   # Command-line interface
│   ├── main.go
│   └── go.mod
├── actions/               # Action modules
│   ├── echo-json/         # Echo action with JSON I/O
│   ├── writefile-json/    # File writing action
│   └── httprequest/       # HTTP request action
├── workflows/             # Sample workflows
│   ├── v1-demo.json       # Comprehensive demo workflow
│   └── simple-test.json   # Basic functionality test
├── bin/                   # Built binaries (created by build.sh)
├── build.sh              # Build script
└── README.md             # This file
```

## 🛠️ Installation & Setup

### Prerequisites
- Go 1.19 or later
- Unix-like system (Linux, macOS, WSL)

### Build

```bash
# Clone or navigate to the project directory
cd go-ai-agent-v1

# Build all components
./build.sh
```

This creates the following binaries in the `bin/` directory:
- `orchestrator` - Main workflow orchestrator
- `cli` - Command-line interface
- `echo-json` - Echo action module
- `writefile-json` - File writing action module
- `httprequest` - HTTP request action module

## 🎯 Usage

### Command Line Interface

The CLI provides two main commands:

#### Run a Workflow
```bash
./bin/cli run <workflow-file> [initial-data]
```

Examples:
```bash
# Run simple test workflow
./bin/cli run workflows/simple-test.json '{"agent_name":"Go-AI-Agent","test_id":"123"}'

# Run demo workflow with comprehensive data
./bin/cli run workflows/v1-demo.json '{"user_name":"Alice","user_id":"1","timestamp":"2024-01-20T10:30:00Z","output_dir":"/tmp/reports"}'
```

#### Validate a Workflow
```bash
./bin/cli validate <workflow-file>
```

### Direct Orchestrator Usage
```bash
# Run with workflow file and initial data
echo '{"user_name":"Alice","user_id":"1","timestamp":"2024-01-20T10:30:00Z","output_dir":"/tmp/reports"}' | ./bin/orchestrator workflows/v1-demo.json
```

## 📄 Workflow Format

V1 workflows use JSON format with enhanced templating capabilities:

```json
{
  "name": "Workflow Name",
  "description": "Workflow description",
  "version": "1.0",
  "steps": [
    {
      "id": "unique_step_id",
      "name": "Human-readable step name",
      "action": "action-module-name",
      "inputs": {
        "param1": "value with {{.WorkflowData.variable}} templating",
        "param2": "reference to {{.Nodes.previous_step.output_field}}"
      }
    }
  ]
}
```

### Templating System

V1 uses Go's `text/template` package with two main contexts:

1. **WorkflowData**: Initial data passed to the workflow
   ```json
   "message": "Hello {{.WorkflowData.user_name}}"
   ```

2. **Nodes**: Outputs from previous workflow steps
   ```json
   "content": "Status: {{.Nodes.api_call.status_code}}"
   ```

## 🔧 Action Modules

### echo-json
Echoes a message with optional formatting.

**Input:**
```json
{
  "message": "string"
}
```

**Output:**
```json
{
  "success": true,
  "message": "echoed message"
}
```

### writefile-json
Writes content to files with various modes.

**Input:**
```json
{
  "path": "/path/to/file",
  "content": "file content",
  "mode": "create|append|overwrite",
  "mkdir_all": true
}
```

**Output:**
```json
{
  "success": true,
  "message": "Success message",
  "path": "/path/to/file",
  "size": 1234
}
```

### httprequest
Makes HTTP requests with full feature support.

**Input:**
```json
{
  "url": "https://api.example.com/data",
  "method": "GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS",
  "headers": {"key": "value"},
  "body": "request body",
  "timeout": 30
}
```

**Output:**
```json
{
  "success": true,
  "message": "Request completed",
  "status_code": 200,
  "headers": {"content-type": "application/json"},
  "body": "response body"
}
```

## 🧪 Testing

### Quick Test
```bash
# Build and run simple test
./build.sh
./bin/cli run workflows/simple-test.json '{"agent_name":"Go-AI-Agent","test_id":"test-'$(date +%s)'"}'
```

### Comprehensive Demo
```bash
# Run the full demo workflow
mkdir -p /tmp/reports
./bin/cli run workflows/v1-demo.json '{"user_name":"Alice","user_id":"1","timestamp":"'$(date -Iseconds)'","output_dir":"/tmp/reports"}'
```

### Manual Testing
```bash
# Test individual action modules
echo '{"message":"Hello V1!"}' | ./bin/echo-json
echo '{"path":"/tmp/test.txt","content":"Test content"}' | ./bin/writefile-json
echo '{"url":"https://httpbin.org/get","method":"GET"}' | ./bin/httprequest
```

## 🏗️ Architecture

### Communication Protocol
- All inter-module communication uses JSON via stdin/stdout
- Structured error responses with consistent format
- Graceful error handling with fallback mechanisms

### Templating Engine
- Go `text/template` package for robust template processing
- Support for complex expressions and functions
- Context-aware variable resolution

### Error Handling
- Structured JSON error responses
- Non-zero exit codes for failures
- Comprehensive logging throughout execution

### Modularity
- Each action module is a standalone executable
- Orchestrator discovers and executes actions dynamically
- Easy to extend with new action modules

## 🔍 Debugging

### Enable Verbose Logging
Set environment variable for detailed logs:
```bash
export GO_AI_AGENT_DEBUG=true
./bin/cli run workflows/simple-test.json '{"test_id":"debug-test"}'
```

### Common Issues

1. **Action Module Not Found**: Ensure the action binary exists in the same directory as the orchestrator
2. **Template Errors**: Check template syntax and available variables
3. **JSON Parsing Errors**: Validate JSON format of workflows and initial data
4. **Permission Errors**: Ensure proper file/directory permissions for file operations

## 🚀 Development

### Adding New Action Modules

1. Create new directory under `actions/`
2. Implement JSON stdin/stdout protocol
3. Follow the error handling pattern
4. Add to `build.sh`
5. Update documentation

### Example Action Module Template
```go
package main

import (
    "encoding/json"
    "os"
)

type ActionInput struct {
    // Define input structure
}

type ActionOutput struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}

func main() {
    var input ActionInput
    decoder := json.NewDecoder(os.Stdin)
    if err := decoder.Decode(&input); err != nil {
        sendErrorResponse("Failed to parse JSON input", err.Error())
        return
    }
    
    // Process action logic
    
    output := ActionOutput{
        Success: true,
        Message: "Action completed successfully",
    }
    
    encoder := json.NewEncoder(os.Stdout)
    encoder.Encode(output)
}

func sendErrorResponse(message, errorDetail string) {
    output := ActionOutput{
        Success: false,
        Message: message,
        Error:   errorDetail,
    }
    encoder := json.NewEncoder(os.Stdout)
    encoder.Encode(output)
    os.Exit(1)
}
```

## 📝 License

This project is part of the octo-agent repository. See the main repository for license information.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📞 Support

For issues and questions:
1. Check the debugging section above
2. Review the workflow examples
3. Create an issue in the main repository
