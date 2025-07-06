package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// EchoInput represents the expected input structure
type EchoInput struct {
	Message string `yaml:"message"`
	Prefix  string `yaml:"prefix,omitempty"`
}

// EchoOutput represents the output structure
type EchoOutput struct {
	EchoedMessage string      `yaml:"echoed_message"`
	OriginalInput interface{} `yaml:"original_input"`
}

// ErrorOutput represents error response structure
type ErrorOutput struct {
	Error           string      `yaml:"error"`
	OriginalRequest interface{} `yaml:"original_request,omitempty"`
}

func main() {
	// Read YAML input from stdin
	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		outputError("Failed to read input from stdin", nil)
		return
	}

	// Parse input YAML
	var input EchoInput
	if err := yaml.Unmarshal(inputData, &input); err != nil {
		outputError("Invalid input YAML format", string(inputData))
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
	outputYAML, err := yaml.Marshal(output)
	if err != nil {
		outputError("Failed to marshal output YAML", input)
		return
	}

	fmt.Print(string(outputYAML))
}

// outputError outputs an error response and exits
func outputError(message string, originalRequest interface{}) {
	errorOutput := ErrorOutput{
		Error:           message,
		OriginalRequest: originalRequest,
	}

	outputYAML, err := yaml.Marshal(errorOutput)
	if err != nil {
		// Fallback error output
		fmt.Printf("error: \"Failed to marshal error response: %s\"\n", message)
		os.Exit(1)
	}

	fmt.Print(string(outputYAML))
	os.Exit(0) // Exit 0 for graceful error reporting
}
