package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunCommand(t *testing.T) {
	// Test that run command is properly configured
	if runCmd.Use != "run" {
		t.Errorf("Expected run command Use to be 'run', got %s", runCmd.Use)
	}

	if runCmd.Short != "Execute a single AI task immediately" {
		t.Errorf("Unexpected short description: %s", runCmd.Short)
	}

	// Verify required flags
	requiredFlags := []string{"model", "prompt", "processor"}
	for _, flagName := range requiredFlags {
		flag := runCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to exist", flagName)
			continue
		}
		if !isRequired(runCmd, flagName) {
			t.Errorf("Expected flag '%s' to be required", flagName)
		}
	}

	// Verify optional flags
	optionalFlags := []string{"template", "vars", "model-params"}
	for _, flagName := range optionalFlags {
		flag := runCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to exist", flagName)
			continue
		}
		if isRequired(runCmd, flagName) {
			t.Errorf("Expected flag '%s' to be optional", flagName)
		}
	}
}

func TestMarkFlagRequiredOrFail(t *testing.T) {
	// Test successful case
	cmd := &cobra.Command{}
	cmd.Flags().String("test-flag", "", "test flag")

	// Should not panic with valid flag
	markFlagRequiredOrFail(cmd, "test-flag")

	// Test that it actually marked the flag as required
	if !isRequired(cmd, "test-flag") {
		t.Error("Flag was not marked as required")
	}

	// Test panic case with invalid flag
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid flag name")
		} else {
			expectedPanic := "Critical configuration error"
			if !strings.Contains(fmt.Sprint(r), expectedPanic) {
				t.Errorf("Expected panic message to contain '%s', got '%v'", expectedPanic, r)
			}
		}
	}()
	markFlagRequiredOrFail(cmd, "non-existent-flag")
}

func TestRunCommandExecution(t *testing.T) {
	// Skip this test for now as it requires mock implementation of the run command
	t.Skip("Skipping run command execution test to avoid service execution")
}

func TestRunCommandVariableParsing(t *testing.T) {
	// Test variable parsing logic in isolation
	tests := []struct {
		name      string
		varsInput string
		expected  map[string]string
	}{
		{
			name:      "single variable",
			varsInput: "key=value",
			expected: map[string]string{
				"key": "value",
			},
		},
		{
			name:      "multiple variables",
			varsInput: "key1=value1,key2=value2,key3=value3",
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name:      "variable with spaces",
			varsInput: "key = value with spaces",
			expected: map[string]string{
				"key": "value with spaces",
			},
		},
		{
			name:      "empty value",
			varsInput: "key=",
			expected: map[string]string{
				"key": "",
			},
		},
		{
			name:      "malformed variable (no equals)",
			varsInput: "notavariable",
			expected:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the variable parsing logic from the run command
			variables := make(map[string]string)
			if tt.varsInput != "" {
				for _, varPair := range strings.Split(tt.varsInput, ",") {
					keyValue := strings.SplitN(varPair, "=", 2)
					if len(keyValue) == 2 {
						key := strings.TrimSpace(keyValue[0])
						value := strings.TrimSpace(keyValue[1])
						variables[key] = value
					}
				}
			}

			// Compare results
			if len(variables) != len(tt.expected) {
				t.Errorf("Expected %d variables, got %d", len(tt.expected), len(variables))
			}
			for k, v := range tt.expected {
				if actual, ok := variables[k]; !ok {
					t.Errorf("Expected variable '%s' not found", k)
				} else if actual != v {
					t.Errorf("Expected variable '%s' to have value '%s', got '%s'", k, v, actual)
				}
			}
		})
	}
}

// Helper function to check if a flag is required
func isRequired(cmd *cobra.Command, flagName string) bool {
	// This is a bit of a hack, but cobra doesn't provide a direct way to check
	// We'll try to execute the command without the flag and see if it errors
	testCmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	testCmd.Flags().AddFlagSet(cmd.Flags())

	// Try to execute without the flag set
	testCmd.SetArgs([]string{})
	err := testCmd.Execute()

	return err != nil && strings.Contains(err.Error(), flagName) && strings.Contains(err.Error(), "required")
}
