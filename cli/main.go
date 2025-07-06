package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// Status indicators
const (
	statusOK     = "[" + colorGreen + "OK" + colorReset + "]"
	statusFAILED = "[" + colorRed + "FAILED" + colorReset + "]"
	statusINFO   = "[" + colorBlue + "INFO" + colorReset + "]"
	statusWARN   = "[" + colorYellow + "WARN" + colorReset + "]"
)

func main() {
	// Check command-line arguments
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> <args...>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  run <workflow_file.json> [initial_data_json]\n")
		fmt.Fprintf(os.Stderr, "  validate <workflow_file.json>\n")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		runWorkflow()
	case "validate":
		validateWorkflow()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// runWorkflow executes a workflow using the orchestrator
func runWorkflow() {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s run <workflow_file.json> [initial_data_json]\n", os.Args[0])
		os.Exit(1)
	}

	workflowFile := os.Args[2]

	// Find orchestrator binary in the same directory as CLI
	cliPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting CLI path: %v\n", err)
		os.Exit(1)
	}

	orchestratorPath := strings.Replace(cliPath, "cli", "orchestrator", 1)

	// Prepare orchestrator command
	args := []string{workflowFile}

	// Add initial data if provided
	if len(os.Args) == 4 {
		initialData := os.Args[3]

		// Validate that initial data is valid JSON
		var testData map[string]interface{}
		if err := json.Unmarshal([]byte(initialData), &testData); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid initial data JSON: %v\n", err)
			os.Exit(1)
		}

		args = append(args, initialData)
	}

	// Execute the orchestrator
	cmd := exec.Command(orchestratorPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error running orchestrator: %v\n", err)
		os.Exit(1)
	}
}

// validateWorkflow validates the syntax of a workflow JSON file
func validateWorkflow() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s validate <workflow_file.json>\n", os.Args[0])
		os.Exit(1)
	}

	workflowFile := os.Args[2]

	// Read and parse the workflow file
	data, err := os.ReadFile(workflowFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading workflow file: %v\n", err)
		os.Exit(1)
	}

	// Basic JSON validation
	var workflow map[string]interface{}
	if err := json.Unmarshal(data, &workflow); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid JSON syntax: %v\n", err)
		os.Exit(1)
	}

	// Validate required fields
	validationErrors := validateWorkflowStructure(workflow)
	if len(validationErrors) > 0 {
		fmt.Fprintf(os.Stderr, "Workflow validation failed:\n")
		for _, err := range validationErrors {
			fmt.Fprintf(os.Stderr, "  - %s\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("✅ Workflow file '%s' is valid\n", workflowFile)
}

// validateWorkflowStructure performs basic validation of workflow structure
func validateWorkflowStructure(workflow map[string]interface{}) []string {
	var errors []string

	// Check required top-level fields
	requiredFields := []string{"name", "description", "nodes"}
	for _, field := range requiredFields {
		if _, exists := workflow[field]; !exists {
			errors = append(errors, fmt.Sprintf("Missing required field: %s", field))
		}
	}

	// Validate nodes array
	if nodesInterface, exists := workflow["nodes"]; exists {
		if nodes, ok := nodesInterface.([]interface{}); ok {
			if len(nodes) == 0 {
				errors = append(errors, "Nodes array cannot be empty")
			}

			nodeIds := make(map[string]bool)
			for i, nodeInterface := range nodes {
				if node, ok := nodeInterface.(map[string]interface{}); ok {
					// Check required node fields
					nodeRequiredFields := []string{"id", "type", "inputs_from_workflow"}
					for _, field := range nodeRequiredFields {
						if _, exists := node[field]; !exists {
							errors = append(errors, fmt.Sprintf("Node %d missing required field: %s", i, field))
						}
					}

					// Check for duplicate node IDs
					if idInterface, exists := node["id"]; exists {
						if id, ok := idInterface.(string); ok {
							if nodeIds[id] {
								errors = append(errors, fmt.Sprintf("Duplicate node ID: %s", id))
							}
							nodeIds[id] = true

							// Validate ID format
							if strings.TrimSpace(id) == "" {
								errors = append(errors, fmt.Sprintf("Node %d has empty ID", i))
							}
						}
					}

					// Validate type field
					if typeInterface, exists := node["type"]; exists {
						if typeStr, ok := typeInterface.(string); ok {
							if strings.TrimSpace(typeStr) == "" {
								errors = append(errors, fmt.Sprintf("Node %d has empty type", i))
							}
						}
					}
				} else {
					errors = append(errors, fmt.Sprintf("Node %d is not a valid object", i))
				}
			}
		} else {
			errors = append(errors, "Nodes must be an array")
		}
	}

	return errors
}
