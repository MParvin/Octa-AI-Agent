package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// ActionInput represents the input structure for the httprequest action
type ActionInput struct {
	URL     string            `json:"url"`
	Method  string            `json:"method,omitempty"` // Default: GET
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Timeout int               `json:"timeout,omitempty"` // Timeout in seconds, default: 30
}

// ActionOutput represents the output structure for the httprequest action
type ActionOutput struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	StatusCode int               `json:"status_code,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
	Error      string            `json:"error,omitempty"`
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
	if input.URL == "" {
		sendErrorResponse("Missing required field", "url is required")
		return
	}

	// Set defaults
	if input.Method == "" {
		input.Method = "GET"
	}
	if input.Timeout == 0 {
		input.Timeout = 30
	}

	// Validate HTTP method
	validMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
	}
	method := strings.ToUpper(input.Method)
	if !validMethods[method] {
		sendErrorResponse("Invalid HTTP method", fmt.Sprintf("method must be one of: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, got: %s", input.Method))
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(input.Timeout) * time.Second,
	}

	// Prepare request body
	var bodyReader io.Reader
	if input.Body != "" {
		bodyReader = strings.NewReader(input.Body)
	}

	// Create HTTP request
	req, err := http.NewRequest(method, input.URL, bodyReader)
	if err != nil {
		sendErrorResponse("Failed to create HTTP request", err.Error())
		return
	}

	// Set headers
	for key, value := range input.Headers {
		req.Header.Set(key, value)
	}

	// Set default Content-Type if body is provided and no Content-Type is set
	if input.Body != "" && req.Header.Get("Content-Type") == "" {
		// Try to detect if it's JSON
		if strings.TrimSpace(input.Body)[0] == '{' || strings.TrimSpace(input.Body)[0] == '[' {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "text/plain")
		}
	}

	log.Printf("Making %s request to %s", method, input.URL)

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		sendErrorResponse("Failed to execute HTTP request", err.Error())
		return
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		sendErrorResponse("Failed to read response body", err.Error())
		return
	}

	// Convert response headers to map
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0] // Take first value if multiple
		}
	}

	// Send success response
	output := ActionOutput{
		Success:    true,
		Message:    fmt.Sprintf("HTTP %s request to %s completed successfully", method, input.URL),
		StatusCode: resp.StatusCode,
		Headers:    responseHeaders,
		Body:       string(bodyBytes),
	}

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(output); err != nil {
		log.Printf("Failed to encode output: %v", err)
		os.Exit(1)
	}

	log.Printf("Request completed with status code: %d", resp.StatusCode)
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
