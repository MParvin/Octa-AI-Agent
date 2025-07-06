package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
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

// WorkflowV1 represents the V1 workflow definition structure
type WorkflowV1 struct {
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	WorkflowDataSchema map[string]string `json:"workflow_data_schema,omitempty"`
	Nodes              []NodeV1          `json:"nodes"`
}

// NodeV1 represents a V1 action node with JSON-based input
type NodeV1 struct {
	ID                 string                 `json:"id"`
	Type               string                 `json:"type"`
	InputsFromWorkflow map[string]interface{} `json:"inputs_from_workflow"`
}

// TemplateContext holds data available for templating
type TemplateContext struct {
	WorkflowData map[string]interface{} `json:"workflow_data"`
	Nodes        map[string]NodeOutput  `json:"nodes"`
}

// NodeOutput stores the JSON output from executed nodes
type NodeOutput struct {
	Output map[string]interface{} `json:"output"`
	Error  string                 `json:"error,omitempty"`
}

// ActionError represents an error response from an action
type ActionError struct {
	Error           string                 `json:"error"`
	OriginalRequest map[string]interface{} `json:"original_request,omitempty"`
}

func main() {
	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Check command-line arguments
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <workflow_file.json> [initial_data_json]\n", os.Args[0])
		os.Exit(1)
	}

	workflowFile := os.Args[1]
	var initialData map[string]interface{}

	// Parse optional initial data
	if len(os.Args) == 3 {
		initialDataStr := os.Args[2]
		if err := json.Unmarshal([]byte(initialDataStr), &initialData); err != nil {
			log.Fatalf("Error parsing initial data JSON: %v", err)
		}
	} else {
		initialData = make(map[string]interface{})
	}

	// Parse workflow definition
	workflow, err := parseWorkflowV1(workflowFile)
	if err != nil {
		log.Fatalf("Error parsing workflow file: %v", err)
	}

	// Update the workflow start message
	log.Printf("Starting workflow execution: %s %s", workflow.Name, statusINFO)
	log.Printf("Description: %s %s", workflow.Description, statusINFO)

	// Execute workflow
	if err := executeWorkflowV1(workflow, initialData); err != nil {
		// Update failure messages
		log.Printf("Workflow execution failed: %s %s", err, statusFAILED)
		log.Fatalf("Workflow execution failed: %v %s", err, statusFAILED)
	}

	// Update workflow completion message
	log.Printf("Workflow completed successfully %s", statusOK)
}

// parseWorkflowV1 reads and parses the V1 workflow JSON file
func parseWorkflowV1(filename string) (*WorkflowV1, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var workflow WorkflowV1
	if err := json.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &workflow, nil
}

// executeWorkflowV1 executes all nodes in the workflow
func executeWorkflowV1(workflow *WorkflowV1, initialData map[string]interface{}) error {
	// Initialize template context
	context := &TemplateContext{
		WorkflowData: initialData,
		Nodes:        make(map[string]NodeOutput),
	}

	// Execute each node sequentially
	for _, node := range workflow.Nodes {
		// Update node execution messages
		log.Printf("Executing node: %s (%s) %s", node.ID, node.Type, statusINFO)

		output, err := executeNodeV1(node, context)
		if err != nil {
			log.Printf("Node %s execution failed %s", node.ID, statusFAILED)
			return fmt.Errorf("error executing node %s: %w", node.ID, err)
		}

		// Store the output for future template resolution
		context.Nodes[node.ID] = NodeOutput{
			Output: output,
		}

		// Update node completion message
		log.Printf("Node %s completed successfully %s", node.ID, statusOK)
	}

	return nil
}

// executeNodeV1 executes a single V1 node
func executeNodeV1(node NodeV1, context *TemplateContext) (map[string]interface{}, error) {
	// Resolve templates in the input
	resolvedInput, err := resolveTemplates(node.InputsFromWorkflow, context)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve templates: %w", err)
	}

	// Convert resolved input to JSON
	inputJSON, err := json.Marshal(resolvedInput)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input JSON: %w", err)
	}

	log.Printf("Sending to action: %s", string(inputJSON))

	// Get the path to the orchestrator binary
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	// Construct action binary path in same directory
	actionPath := strings.Replace(execPath, "orchestrator", node.Type, 1)

	// Execute the action binary
	cmd := exec.Command(actionPath)
	cmd.Stdin = strings.NewReader(string(inputJSON))

	// Capture stdout and stderr separately
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Log stderr for debugging
		if stderr.Len() > 0 {
			// Update the stderr capture message
			log.Printf("Action stderr output: %s %s", stderr.String(), statusWARN)
		}
		return nil, fmt.Errorf("action failed: %w", err)
	}

	// Log stderr if present (for debugging)
	if stderr.Len() > 0 {
		log.Printf("Action stderr: %s %s", stderr.String(), statusINFO)
	}

	// Parse the output JSON
	var output map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		return nil, fmt.Errorf("failed to parse action output JSON: %w", err)
	}

	// Check for error in the output
	if errorMsg, exists := output["error"]; exists {
		return nil, fmt.Errorf("action returned error: %v", errorMsg)
	}

	log.Printf("Action output: %s %s", stdout.String(), statusINFO)
	return output, nil
}

// resolveTemplates processes templates in the input data
func resolveTemplates(input map[string]interface{}, context *TemplateContext) (map[string]interface{}, error) {
	// Convert input to JSON and back to handle nested structures
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// Apply template processing to the JSON string
	tmpl, err := template.New("input").Parse(string(inputJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Parse the resolved JSON back
	var resolved map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &resolved); err != nil {
		return nil, fmt.Errorf("failed to parse resolved JSON: %w", err)
	}

	return resolved, nil
}
