# Examples Directory

This directory contains example workflows for the Go AI Agent V1.0, ranging from simple to professional use cases. All examples are now available in YAML format for better readability and ease of editing.

## Quick Start

1. **Build the project:**
   ```bash
   cd ..
   ./build.sh
   ```

2. **Run any example:**
   ```bash
   # Hello World (Simple)
   ../bin/cli run hello-world.yaml 'name: "Alice"'
   
   # File Operations (Basic)
   ../bin/cli run file-operations.yaml 'user_id: "12345"
timestamp: "2024-01-20T10:30:00Z"'
   
   # API Integration (Intermediate)  
   mkdir -p /tmp/reports
   ../bin/cli run api-integration.yaml 'user_id: "3"
timestamp: "2024-01-20T10:30:00Z"
output_dir: "/tmp/reports"'
   
   # Data Pipeline (Advanced)
   mkdir -p /tmp/data-pipeline
   ../bin/cli run data-pipeline.yaml 'batch_id: "batch_001"
timestamp: "2024-01-20T10:30:00Z"
output_dir: "/tmp/data-pipeline"
record_count: "100"'
   
   # CI/CD Pipeline (Professional)
   mkdir -p /tmp/cicd/{artifacts,reports}
   ../bin/cli run cicd-pipeline.yaml 'project_name: "web-application"
version: "2.1.0"
branch: "main"
commit_sha: "abc123def456"
environment: "staging"
pipeline_id: "pipeline_001"
timestamp: "2024-01-20T10:30:00Z"
artifacts_dir: "/tmp/cicd/artifacts"
reports_dir: "/tmp/cicd/reports"
api_token: "sk-test-token"
deploy_token: "deploy-token-123"'
   
   # Claude AI Integration (AI-Powered)
   export CLAUDE_API_KEY="your-claude-api-key-here"
   ../bin/cli run claude-ai-integration.yaml '{}'
   
   # Git Repository Watcher (Monitoring)
   ../bin/cli run git-watch.yaml 'repo_url: "https://github.com/owner/repo.git"
branch: "main"
username: ""
password: ""
timestamp: "2024-01-20T10:30:00Z"'
   ```

3. **Validate workflows:**
   ```bash
   ../bin/cli validate hello-world.yaml
   ```

## Example Files

- `hello-world.yaml` - Simple greeting with templating
- `file-operations.yaml` - File creation and content templating  
- `api-integration.yaml` - External API integration with HTTP requests
- `data-pipeline.yaml` - Advanced data processing with validation
- `cicd-pipeline.yaml` - Professional CI/CD workflow
- `claude-ai-integration.yaml` - AI-powered text generation
- `git-watch.yaml` - Git repository monitoring

## Benefits of YAML Format

The YAML format provides several advantages over JSON:

1. **Human Readability**: Much easier to read and edit manually
2. **Multi-line Support**: Natural handling of multi-line content without escaping
3. **Comments**: Ability to add comments for documentation (future enhancement)
4. **Cleaner Syntax**: No need for excessive quotes and brackets
5. **Better Diffs**: More readable version control diffs

## Legacy JSON Support

All examples are also available in JSON format for backward compatibility. The system supports both formats simultaneously.  
- `api-integration.json` - HTTP API calls and data processing
- `data-pipeline.json` - Advanced pipeline with validation
- `cicd-pipeline.json` - Professional CI/CD workflow
- `git-watch.json` - Git repository monitoring and change detection

For detailed explanations of each example, see the main [EXAMPLE.md](../EXAMPLE.md) file.
