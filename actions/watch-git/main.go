package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"gopkg.in/yaml.v3"
)

// ActionInput represents the input structure for the watch-git action
type ActionInput struct {
	URL          string `yaml:"url"`                      // Git repository URL
	Username     string `yaml:"username,omitempty"`       // Git username (optional)
	Password     string `yaml:"password,omitempty"`       // Git password/token (optional)
	Branch       string `yaml:"branch,omitempty"`         // Branch to watch (default: main)
	Interval     int    `yaml:"interval,omitempty"`       // Check interval in seconds (default: 60)
	MaxChecks    int    `yaml:"max_checks,omitempty"`     // Max number of checks (default: 10)
	LocalDir     string `yaml:"local_dir,omitempty"`      // Local directory to clone to (optional)
	ExitOnChange *bool  `yaml:"exit_on_change,omitempty"` // Exit immediately when first change is detected (default: true)
}

// ActionOutput represents the output structure for the watch-git action
type ActionOutput struct {
	Success    bool     `yaml:"success"`
	Message    string   `yaml:"message"`
	URL        string   `yaml:"url"`
	Branch     string   `yaml:"branch"`
	LastCommit string   `yaml:"last_commit,omitempty"`
	Changes    []Change `yaml:"changes,omitempty"`
	CheckCount int      `yaml:"check_count"`
	Error      string   `yaml:"error,omitempty"`
}

// Change represents a detected change in the repository
type Change struct {
	CommitHash   string   `yaml:"commit_hash"`
	Author       string   `yaml:"author"`
	Message      string   `yaml:"message"`
	Timestamp    string   `yaml:"timestamp"`
	FilesChanged []string `yaml:"files_changed"`
}

func main() {
	// Read YAML input from stdin
	var input ActionInput

	// Read all input from stdin
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		sendErrorResponse("Failed to read input from stdin", err.Error())
		return
	}

	// Parse YAML input
	if err := yaml.Unmarshal(inputData, &input); err != nil {
		sendErrorResponse("Failed to parse YAML input", err.Error())
		return
	}

	// Validate required fields
	if input.URL == "" {
		sendErrorResponse("Missing required field", "url is required")
		return
	}

	// Set defaults
	if input.Branch == "" {
		input.Branch = "main"
	}
	if input.Interval == 0 {
		input.Interval = 60
	}
	if input.MaxChecks == 0 {
		input.MaxChecks = 10
	}
	// Default to exit on first change (true by default)
	exitOnChange := true
	if input.ExitOnChange != nil {
		exitOnChange = *input.ExitOnChange
	}

	// Setup authentication if credentials are provided
	var auth *http.BasicAuth
	if input.Username != "" && input.Password != "" {
		auth = &http.BasicAuth{
			Username: input.Username,
			Password: input.Password,
		}
	}

	// Watch the repository
	changes, lastCommit, actualChecks, err := watchRepository(input.URL, input.Branch, auth, input.Interval, input.MaxChecks, exitOnChange)
	if err != nil {
		sendErrorResponse("Failed to watch repository", err.Error())
		return
	}

	// Determine the message based on whether changes were found
	var message string
	if len(changes) > 0 {
		if exitOnChange {
			message = fmt.Sprintf("Change detected in repository %s on branch %s - exited early", input.URL, input.Branch)
		} else {
			message = fmt.Sprintf("Successfully watched repository %s on branch %s - found %d changes", input.URL, input.Branch, len(changes))
		}
	} else {
		message = fmt.Sprintf("No changes detected in repository %s on branch %s", input.URL, input.Branch)
	}

	// Send success response
	output := ActionOutput{
		Success:    true,
		Message:    message,
		URL:        input.URL,
		Branch:     input.Branch,
		LastCommit: lastCommit,
		Changes:    changes,
		CheckCount: actualChecks,
	}

	outputYAML, err := yaml.Marshal(output)
	if err != nil {
		sendErrorResponse("Failed to marshal output YAML", err.Error())
		return
	}

	fmt.Print(string(outputYAML))
}

// watchRepository watches a git repository for changes
func watchRepository(url, branch string, auth *http.BasicAuth, interval, maxChecks int, exitOnChange bool) ([]Change, string, int, error) {
	var changes []Change
	var lastKnownCommit string
	var actualChecks int

	log.Printf("Starting to watch repository: %s, branch: %s", url, branch)

	// Initial clone to get the current state
	initialCommit, err := getLatestCommit(url, branch, auth)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to get initial commit: %v", err)
	}

	lastKnownCommit = initialCommit.Hash
	log.Printf("Initial commit: %s", lastKnownCommit)

	// Watch for changes
	for i := 0; i < maxChecks; i++ {
		actualChecks = i + 1

		if i > 0 {
			time.Sleep(time.Duration(interval) * time.Second)
		}

		log.Printf("Checking for changes... (check %d/%d)", i+1, maxChecks)

		latestCommit, err := getLatestCommit(url, branch, auth)
		if err != nil {
			log.Printf("Error checking for changes: %v", err)
			continue
		}

		if latestCommit.Hash != lastKnownCommit {
			log.Printf("New commit detected: %s", latestCommit.Hash)

			change := Change{
				CommitHash:   latestCommit.Hash,
				Author:       latestCommit.Author,
				Message:      latestCommit.Message,
				Timestamp:    latestCommit.Timestamp,
				FilesChanged: latestCommit.FilesChanged,
			}
			changes = append(changes, change)
			lastKnownCommit = latestCommit.Hash

			// Exit early if configured to do so
			if exitOnChange {
				log.Printf("Exiting early due to change detection (exit_on_change=true)")
				break
			}
		}
	}

	return changes, lastKnownCommit, actualChecks, nil
}

// CommitInfo represents basic commit information
type CommitInfo struct {
	Hash         string
	Author       string
	Message      string
	Timestamp    string
	FilesChanged []string
}

// getLatestCommit gets the latest commit from the repository
func getLatestCommit(url, branch string, auth *http.BasicAuth) (*CommitInfo, error) {
	// Clone the repository in memory with depth 2 to get parent for comparison
	cloneOptions := &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         2, // Need parent commit to compare changes
	}

	if auth != nil {
		cloneOptions.Auth = auth
	}

	repo, err := git.Clone(memory.NewStorage(), nil, cloneOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %v", err)
	}

	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %v", err)
	}

	// Get the commit object
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %v", err)
	}

	// Get the files changed in this commit
	filesChanged, err := getCommitFileChanges(repo, commit)
	if err != nil {
		log.Printf("Warning: could not get file changes for commit %s: %v", commit.Hash.String(), err)
		filesChanged = []string{} // Continue without file changes
	}

	return &CommitInfo{
		Hash:         commit.Hash.String(),
		Author:       commit.Author.Name,
		Message:      strings.TrimSpace(commit.Message),
		Timestamp:    commit.Author.When.Format(time.RFC3339),
		FilesChanged: filesChanged,
	}, nil
}

// getCommitFileChanges gets the list of files changed in a commit
func getCommitFileChanges(repo *git.Repository, commit *object.Commit) ([]string, error) {
	var changedFiles []string

	// Get the tree for this commit
	commitTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit tree: %v", err)
	}

	// If this commit has no parents, it's the initial commit - return all files
	if commit.NumParents() == 0 {
		err = commitTree.Files().ForEach(func(file *object.File) error {
			changedFiles = append(changedFiles, file.Name)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to iterate files in initial commit: %v", err)
		}
		return changedFiles, nil
	}

	// Get the parent commit
	parent, err := commit.Parent(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent commit: %v", err)
	}

	// Get the parent tree
	parentTree, err := parent.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get parent tree: %v", err)
	}

	// Compare trees to find changes
	changes, err := parentTree.Diff(commitTree)
	if err != nil {
		return nil, fmt.Errorf("failed to diff trees: %v", err)
	}

	// Extract file names from changes
	for _, change := range changes {
		from, to := change.From, change.To

		// Handle different types of changes
		if from.Name != "" && to.Name != "" {
			// File renamed
			if from.Name != to.Name {
				changedFiles = append(changedFiles, from.Name+" -> "+to.Name)
			} else {
				// File modified
				changedFiles = append(changedFiles, to.Name)
			}
		} else if from.Name != "" {
			// File deleted
			changedFiles = append(changedFiles, from.Name+" (deleted)")
		} else if to.Name != "" {
			// File added
			changedFiles = append(changedFiles, to.Name+" (added)")
		}
	}

	return changedFiles, nil
}

// sendErrorResponse sends an error response and exits
func sendErrorResponse(message string, errorDetail string) {
	output := ActionOutput{
		Success: false,
		Message: message,
		Error:   errorDetail,
	}

	outputYAML, err := yaml.Marshal(output)
	if err != nil {
		// Fallback error output in YAML format
		fmt.Printf("success: false\nmessage: \"Failed to marshal error response\"\nerror: \"%s\"\n", message)
		os.Exit(1)
	}

	fmt.Print(string(outputYAML))
	os.Exit(0) // Exit 0 for graceful error reporting
}
