package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ActionInput represents the input structure for the writefile-json action
type ActionInput struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Mode     string `json:"mode,omitempty"`      // "create", "append", "overwrite" (default: create)
	MkdirAll bool   `json:"mkdir_all,omitempty"` // Create parent directories if they don't exist
}

// ActionOutput represents the output structure for the writefile-json action
type ActionOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Path    string `json:"path"`
	Size    int64  `json:"size,omitempty"`
	Error   string `json:"error,omitempty"`
}

func main() {
	// Read JSON input from stdin
	var input ActionInput
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&input); err != nil {
		sendErrorResponse("Failed to parse JSON input", err.Error())
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
	var err error
	var file *os.File
	var fileInfo os.FileInfo

	switch input.Mode {
	case "create":
		// Check if file exists
		if _, err := os.Stat(input.Path); err == nil {
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

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(output); err != nil {
		log.Printf("Failed to encode output: %v", err)
		os.Exit(1)
	}
}

func sendErrorResponse(message, errorDetail string) {
	output := ActionOutput{
		Success: false,
		Message: message,
		Error:   errorDetail,
	}

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(output); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
	os.Exit(1)
}
