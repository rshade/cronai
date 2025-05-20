package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	// Capture output
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Run main function
	main()

	// Restore stdout and read output
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}
	os.Stdout = old
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	output := string(out)

	// Check for expected output
	expectedStrings := []string{
		"Result as string:",
		"# Report for CronAI",
		"## Production Environment Status",
		"Current status: healthy",
		"All systems operational.",
		"Report generated on 2025-05-12",
		"Result as bytes:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}

	// Check that it doesn't contain test environment message
	if strings.Contains(output, "Test Environment Status") {
		t.Error("Output should not contain test environment status for production")
	}
}

func TestMainWithError(t *testing.T) {
	// This is harder to test since we can't mock the template manager easily
	// In a real scenario, we'd refactor main() to accept dependencies
	// For now, we'll test what we can

	// The main function as written will always succeed with the hardcoded template
	// This test documents that limitation
	t.Log("Main function error handling cannot be tested without refactoring to accept dependencies")
}
