# Examples Directory

This directory   # Git Repository Watcher (Monitoring)
   ../bin/cli run git-watch.json '{"repo_url":"https://github.com/octocat/Hello-World.git","branch":"master","username":"","password":"","timestamp":"2024-01-20T10:30:00Z"}'ontains example workflows for the Go AI Agent V1.0, ranging from simple to professional use cases.

## Quick Start

1. **Build the project:**
   ```bash
   cd ..
   ./build.sh
   ```

2. **Run any example:**
   ```bash
   # Hello World (Simple)
   ../bin/cli run hello-world.json '{"name":"Alice"}'
   
   # File Operations (Basic)
   ../bin/cli run file-operations.json '{"user_id":"12345","timestamp":"2024-01-20T10:30:00Z"}'
   
   # API Integration (Intermediate)  
   mkdir -p /tmp/reports
   ../bin/cli run api-integration.json '{"user_id":"3","timestamp":"2024-01-20T10:30:00Z","output_dir":"/tmp/reports"}'
   
   # Data Pipeline (Advanced)
   mkdir -p /tmp/data-pipeline
   ../bin/cli run data-pipeline.json '{"batch_id":"batch_001","timestamp":"2024-01-20T10:30:00Z","output_dir":"/tmp/data-pipeline","record_count":"100"}'
   
   # CI/CD Pipeline (Professional)
   mkdir -p /tmp/cicd/{artifacts,reports}
   ../bin/cli run cicd-pipeline.json '{"project_name":"web-application","version":"2.1.0","branch":"main","commit_sha":"abc123def456","environment":"staging","pipeline_id":"pipeline_001","timestamp":"2024-01-20T10:30:00Z","artifacts_dir":"/tmp/cicd/artifacts","reports_dir":"/tmp/cicd/reports","api_token":"sk-test-token","deploy_token":"deploy-token-123"}'
   
   # Claude AI Integration (AI-Powered)
   export CLAUDE_API_KEY="your-claude-api-key-here"
   ../bin/cli run claude-ai-integration.json '{}'
   
   # Git Repository Watcher (Monitoring)
   ../bin/cli run git-watch.json '{"repo_url":"https://github.com/owner/repo.git","branch":"main","username":"","password":"","interval":"30","max_checks":"5","timestamp":"2024-01-20T10:30:00Z"}'
   ```

3. **Validate workflows:**
   ```bash
   ../bin/cli validate hello-world.json
   ```

## Example Files

- `hello-world.json` - Simple greeting with templating
- `file-operations.json` - File creation and content templating  
- `api-integration.json` - HTTP API calls and data processing
- `data-pipeline.json` - Advanced pipeline with validation
- `cicd-pipeline.json` - Professional CI/CD workflow
- `git-watch.json` - Git repository monitoring and change detection

For detailed explanations of each example, see the main [EXAMPLE.md](../EXAMPLE.md) file.
