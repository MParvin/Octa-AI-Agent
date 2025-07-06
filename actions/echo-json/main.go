package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// EchoInput represents the expected input structure
type EchoInput struct {
	Message string `json:"message"`
	Prefix  string `json:"prefix,omitempty"`
}

// EchoOutput represents the output structure
type EchoOutput struct {
	EchoedMessage string      `json:"echoed_message"`
	OriginalInput interface{} `json:"original_input"`
}

// ErrorOutput represents error response structure
type ErrorOutput struct {
	Error           string      `json:"error"`
	OriginalRequest interface{} `json:"original_request,omitempty"`
}

func main() {
	// Read JSON input from stdin
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		outputError("Failed to read input from stdin", nil)
		return
	}

	// Parse input JSON
	var input EchoInput
	if err := json.Unmarshal(inputData, &input); err != nil {
		outputError("Invalid input JSON format", string(inputData))
		return
	}

	// Validate required fields
	if input.Message == "" {
		outputError("Missing required field: message", input)
		return
	}

	// Process the echo operation
	prefix := input.Prefix
	if prefix == "" {
		prefix = ""
	}

	echoedMessage := prefix + input.Message

	// Create output
	output := EchoOutput{
		EchoedMessage: echoedMessage,
		OriginalInput: input,
	}

	// Marshal and output result
	outputJSON, err := json.Marshal(output)
	if err != nil {
		outputError("Failed to marshal output JSON", input)
		return
	}

	fmt.Print(string(outputJSON))
}

// outputError outputs an error response and exits
func outputError(message string, originalRequest interface{}) {
	errorOutput := ErrorOutput{
		Error:           message,
		OriginalRequest: originalRequest,
	}

	outputJSON, err := json.Marshal(errorOutput)
	if err != nil {
		// Fallback error output
		fmt.Printf(`{"error": "Failed to marshal error response: %s"}`, message)
		os.Exit(1)
	}

	fmt.Print(string(outputJSON))
	os.Exit(0) // Exit 0 for graceful error reporting
}
