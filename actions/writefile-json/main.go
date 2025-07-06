package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ActionInput represents the input structure for the writefile action
type ActionInput struct {
	Path     string `yaml:"path"`
	Content  string `yaml:"content"`
	Mode     string `yaml:"mode,omitempty"`      // "create", "append", "overwrite" (default: create)
	MkdirAll bool   `yaml:"mkdir_all,omitempty"` // Create parent directories if they don't exist
}

// ActionOutput represents the output structure for the writefile action
type ActionOutput struct {
	Success bool   `yaml:"success"`
	Message string `yaml:"message"`
	Path    string `yaml:"path"`
	Size    int64  `yaml:"size,omitempty"`
	Error   string `yaml:"error,omitempty"`
}

func main() {
	// Read YAML input from stdin
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		sendErrorResponse("Failed to read input from stdin", err.Error())
		return
	}

	var input ActionInput
	if err := yaml.Unmarshal(inputData, &input); err != nil {
		sendErrorResponse("Failed to parse YAML input", err.Error())
		return
	}

	// Validate required fields
	if input.Path == "" {
		sendErrorResponse("Missing required field", "path is required")
		return
	}

	if input.Content == "" && input.Mode != "create" {
		log.Printf("Warning: Content is empty for path: %s", input.Path)
	}

	// Set default mode
	if input.Mode == "" {
		input.Mode = "create"
	}

	// Validate mode
	validModes := map[string]bool{
		"create":    true,
		"append":    true,
		"overwrite": true,
	}
	if !validModes[input.Mode] {
		sendErrorResponse("Invalid mode", fmt.Sprintf("mode must be one of: create, append, overwrite, got: %s", input.Mode))
		return
	}

	// Create parent directories if requested
	if input.MkdirAll {
		dir := filepath.Dir(input.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			sendErrorResponse("Failed to create parent directories", err.Error())
			return
		}
	}

	// Handle different write modes
	var file *os.File
	var fileInfo os.FileInfo

	switch input.Mode {
	case "create":
		// Check if file exists
		if _, statErr := os.Stat(input.Path); statErr == nil {
			sendErrorResponse("File already exists", fmt.Sprintf("file %s already exists and mode is 'create'", input.Path))
			return
		}
		file, err = os.Create(input.Path)
	case "append":
		file, err = os.OpenFile(input.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	case "overwrite":
		file, err = os.Create(input.Path)
	}

	if err != nil {
		sendErrorResponse("Failed to open file", err.Error())
		return
	}
	defer file.Close()

	// Write content to file
	bytesWritten, err := file.WriteString(input.Content)
	if err != nil {
		sendErrorResponse("Failed to write content", err.Error())
		return
	}

	// Get file info for size
	fileInfo, err = file.Stat()
	if err != nil {
		log.Printf("Warning: Could not get file info: %v", err)
	}

	// Send success response
	output := ActionOutput{
		Success: true,
		Message: fmt.Sprintf("Successfully wrote %d bytes to %s", bytesWritten, input.Path),
		Path:    input.Path,
		Size:    fileInfo.Size(),
	}

	outputYAML, err := yaml.Marshal(output)
	if err != nil {
		log.Printf("Failed to marshal output: %v", err)
		os.Exit(1)
	}
	fmt.Print(string(outputYAML))
}

func sendErrorResponse(message, errorDetail string) {
	output := ActionOutput{
		Success: false,
		Message: message,
		Error:   errorDetail,
	}

	outputYAML, err := yaml.Marshal(output)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		fmt.Printf("success: false\nmessage: \"%s\"\nerror: \"%s\"\n", message, errorDetail)
	} else {
		fmt.Print(string(outputYAML))
	}
	os.Exit(1)
}
