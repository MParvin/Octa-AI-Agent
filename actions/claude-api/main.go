package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// ActionInput represents the input structure for the claude-api action
type ActionInput struct {
	APIKey       string  `json:"api_key,omitempty"`       // Claude API key (can also be set via CLAUDE_API_KEY env var)
	Model        string  `json:"model,omitempty"`         // Claude model to use (default: claude-3-sonnet-20240229)
	Prompt       string  `json:"prompt"`                  // The prompt/message to send to Claude
	MaxTokens    int     `json:"max_tokens,omitempty"`    // Maximum tokens to generate (default: 1000)
	Temperature  float64 `json:"temperature,omitempty"`   // Temperature for response generation (default: 0.7)
	SystemPrompt string  `json:"system_prompt,omitempty"` // System prompt for Claude
	Timeout      int     `json:"timeout,omitempty"`       // Timeout in seconds (default: 60)
}

// ActionOutput represents the output structure for the claude-api action
type ActionOutput struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Response string `json:"response,omitempty"` // Claude's response
	Model    string `json:"model,omitempty"`    // Model used
	Usage    Usage  `json:"usage,omitempty"`    // Token usage information
	Error    string `json:"error,omitempty"`
}

// Usage represents token usage information from Claude API
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ClaudeMessage represents a message in the Claude API format
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeRequest represents the request body for Claude API
type ClaudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature,omitempty"`
	System      string          `json:"system,omitempty"`
	Messages    []ClaudeMessage `json:"messages"`
}

// ClaudeResponse represents the response from Claude API
type ClaudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Model   string `json:"model"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        Usage  `json:"usage"`
}

// ClaudeErrorResponse represents error response from Claude API
type ClaudeErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
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
	if input.Prompt == "" {
		sendErrorResponse("Missing required field", "prompt is required")
		return
	}

	// Get API key from input or environment variable
	apiKey := input.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("CLAUDE_API_KEY")
	}
	if apiKey == "" {
		sendErrorResponse("Missing API key", "api_key must be provided in input or CLAUDE_API_KEY environment variable must be set")
		return
	}

	// Set defaults
	if input.Model == "" {
		input.Model = "claude-3-sonnet-20240229"
	}
	if input.MaxTokens == 0 {
		input.MaxTokens = 1000
	}
	if input.Temperature == 0 {
		input.Temperature = 0.7
	}
	if input.Timeout == 0 {
		input.Timeout = 60
	}

	// Validate model
	validModels := map[string]bool{
		"claude-3-sonnet-20240229": true,
		"claude-3-opus-20240229":   true,
		"claude-3-haiku-20240307":  true,
		"claude-2.1":               true,
		"claude-2.0":               true,
	}
	if !validModels[input.Model] {
		log.Printf("Warning: Unknown model %s, proceeding anyway", input.Model)
	}

	// Create Claude API request
	claudeReq := ClaudeRequest{
		Model:       input.Model,
		MaxTokens:   input.MaxTokens,
		Temperature: input.Temperature,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: input.Prompt,
			},
		},
	}

	// Add system prompt if provided
	if input.SystemPrompt != "" {
		claudeReq.System = input.SystemPrompt
	}

	// Marshal request to JSON
	requestBody, err := json.Marshal(claudeReq)
	if err != nil {
		sendErrorResponse("Failed to marshal request", err.Error())
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(input.Timeout) * time.Second,
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(requestBody))
	if err != nil {
		sendErrorResponse("Failed to create HTTP request", err.Error())
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Anthropic-Version", "2023-06-01")

	log.Printf("Making request to Claude API with model %s", input.Model)

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		sendErrorResponse("Failed to execute request to Claude API", err.Error())
		return
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		sendErrorResponse("Failed to read response body", err.Error())
		return
	}

	// Check for API errors
	if resp.StatusCode != http.StatusOK {
		var errorResp ClaudeErrorResponse
		if err := json.Unmarshal(bodyBytes, &errorResp); err != nil {
			sendErrorResponse("Claude API error", fmt.Sprintf("Status: %d, Body: %s", resp.StatusCode, string(bodyBytes)))
			return
		}
		sendErrorResponse("Claude API error", fmt.Sprintf("Status: %d, Type: %s, Message: %s", resp.StatusCode, errorResp.Error.Type, errorResp.Error.Message))
		return
	}

	// Parse Claude response
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(bodyBytes, &claudeResp); err != nil {
		sendErrorResponse("Failed to parse Claude response", err.Error())
		return
	}

	// Extract response text
	var responseText string
	if len(claudeResp.Content) > 0 && claudeResp.Content[0].Type == "text" {
		responseText = claudeResp.Content[0].Text
	} else {
		responseText = "No text content in response"
	}

	// Send success response
	output := ActionOutput{
		Success:  true,
		Message:  fmt.Sprintf("Successfully generated response using %s", claudeResp.Model),
		Response: responseText,
		Model:    claudeResp.Model,
		Usage:    claudeResp.Usage,
	}

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(output); err != nil {
		log.Printf("Failed to encode output: %v", err)
		os.Exit(1)
	}

	log.Printf("Request completed successfully. Input tokens: %d, Output tokens: %d", claudeResp.Usage.InputTokens, claudeResp.Usage.OutputTokens)
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
