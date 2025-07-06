# Watch Git Action

The `watch-git` action monitors a Git repository for changes over a specified period of time. It uses **YAML format** for both input and output, making it highly human-readable and easy to configure.

## Features

- Monitor any Git repository (public or private)
- Watch specific branches
- Optional authentication with username/password
- Configurable check intervals and maximum checks
- **Exit on first change**: By default, exits immediately when a change is detected
- **Continuous monitoring**: Can be configured to monitor for multiple changes
- **File change tracking**: Reports which files were added, modified, or deleted in each commit
- **YAML format**: Human-readable input and output using YAML instead of JSON
- Reports detected changes with commit details
- Memory-efficient operation (uses in-memory cloning)

## Input Parameters

| Parameter       | Type   | Required | Default | Description                                    |
|----------------|--------|----------|---------|------------------------------------------------|
| `url`          | string | Yes      | -       | Git repository URL (https or ssh)             |
| `username`     | string | No       | ""      | Git username for authentication                |
| `password`     | string | No       | ""      | Git password/token for authentication         |
| `branch`       | string | No       | "main"  | Branch to monitor                              |
| `interval`     | int    | No       | 60      | Check interval in seconds                      |
| `max_checks`   | int    | No       | 10      | Maximum number of checks to perform            |
| `local_dir`    | string | No       | ""      | Local directory to clone to (optional)        |
| `exit_on_change`| bool   | No       | true    | Exit immediately when first change is detected |

## Output

| Field        | Type     | Description                                    |
|-------------|----------|------------------------------------------------|
| `success`   | bool     | Whether the operation was successful           |
| `message`   | string   | Human-readable status message                  |
| `url`       | string   | The monitored repository URL                   |
| `branch`    | string   | The monitored branch                           |
| `last_commit`| string  | Hash of the last commit found                  |
| `changes`   | Change[] | Array of detected changes                      |
| `check_count`| int     | Number of checks performed                     |
| `error`     | string   | Error message if operation failed              |

### Change Object

| Field         | Type     | Description                    |
|--------------|----------|--------------------------------|
| `commit_hash` | string   | Hash of the commit             |
| `author`     | string   | Author name                    |
| `message`    | string   | Commit message                 |
| `timestamp`  | string   | Commit timestamp (RFC3339)     |
| `files_changed`| string[] | List of files changed in commit |

## Examples

### Exit on First Change (Default Behavior)

Create a YAML file `input.yaml`:
```yaml
url: "https://github.com/owner/repo.git"
branch: "main"
interval: 30
max_checks: 5
```

Run the action:
```bash
cat input.yaml | ./bin/watch-git
```

### Continuous Monitoring (Multiple Changes)

```yaml
url: "https://github.com/owner/repo.git"
branch: "main"
interval: 30
max_checks: 10
exit_on_change: false
```

### With Authentication (Private Repository)

```yaml
url: "https://github.com/owner/private-repo.git"
username: "your-username"
password: "your-token"
branch: "develop"
interval: 60
max_checks: 10
```

### In a Workflow

See `examples/git-watch.yaml` for a complete workflow example that demonstrates how to integrate the watch-git action with other actions to create reports and notifications.

## Authentication

- For private repositories, provide `username` and `password`
- For GitHub, use a Personal Access Token as the password
- For GitLab, use a Deploy Token or Personal Access Token
- SSH URLs are supported but require key-based authentication setup

## Use Cases

1. **CI/CD Monitoring**: Watch for new commits to trigger builds
2. **Security Monitoring**: Monitor repositories for unauthorized changes
3. **Release Tracking**: Watch release branches for new versions
4. **Development Coordination**: Monitor feature branches for team updates
5. **Compliance**: Track changes for audit purposes

## Example Output with Changes

When a change is detected, the output includes detailed file information in YAML format:

```yaml
success: true
message: "Change detected in repository https://github.com/owner/repo.git on branch main - exited early"
url: "https://github.com/owner/repo.git"
branch: "main"
last_commit: "abc123def456..."
changes:
  - commit_hash: "abc123def456..."
    author: "John Doe"
    message: "Fix critical bug in authentication"
    timestamp: "2024-07-06T15:30:00Z"
    files_changed:
      - "src/auth.go"
      - "tests/auth_test.go"
      - "docs/README.md (added)"
      - "old_file.txt (deleted)"
check_count: 2
```

## File Change Annotations

- `filename` - File was modified
- `filename (added)` - File was newly created
- `filename (deleted)` - File was removed
- `oldname -> newname` - File was renamed

## Error Handling

The action gracefully handles various error conditions:
- Invalid repository URLs
- Authentication failures
- Network connectivity issues
- Branch not found
- Repository access permissions

All errors are returned in the JSON output with detailed error messages.
