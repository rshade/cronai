package cron

import (
	"fmt"
	"os"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	// Create a temporary test config file
	testConfigPath := "test_config.tmp"
	testConfigContent := `# Test config file
0 8 * * * claude product_manager slack-pm-channel
0 9 * * 1 openai weekly_report email-team@company.com
*/15 * * * * gemini monitoring_check log-to-file
0 12 * * * claude report_template email-execs@company.com reportType=Weekly,project=CronAI,team=Dev

# Invalid lines
invalid line
0 8 * * * # Missing fields
`

	// Create the test config file
	err := os.WriteFile(testConfigPath, []byte(testConfigContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer func() {
		if err := os.Remove(testConfigPath); err != nil {
			t.Logf("Warning: Failed to remove test config file: %v", err)
		}
	}()

	// Ensure test prompt files exist
	if err := os.MkdirAll("cron_prompts", 0755); err != nil {
		t.Fatalf("Failed to create cron_prompts directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll("cron_prompts"); err != nil {
			t.Logf("Warning: Failed to remove cron_prompts directory: %v", err)
		}
	}()

	// Create test prompt files
	prompts := []string{"product_manager", "weekly_report", "monitoring_check", "report_template"}
	for _, p := range prompts {
		if err := os.WriteFile(fmt.Sprintf("cron_prompts/%s.md", p), []byte("Test prompt content"), 0644); err != nil {
			t.Fatalf("Failed to create test prompt file: %v", err)
		}
		defer func(name string) {
			if err := os.Remove(fmt.Sprintf("cron_prompts/%s.md", name)); err != nil {
				t.Logf("Warning: Failed to remove test prompt file: %v", err)
			}
		}(p)
	}

	// Parse the config file
	tasks, err := parseConfigFile(testConfigPath)
	// We expect validation errors but still get some valid tasks
	if err == nil {
		t.Errorf("Expected some validation errors but got none")
	}

	// Verify the tasks - we expect 4 valid tasks due to the validation
	if len(tasks) != 4 {
		t.Errorf("Expected 4 tasks, got %d", len(tasks))
	}

	// Check the first task
	if tasks[0].Schedule != "0 8 * * *" {
		t.Errorf("Expected schedule '0 8 * * *', got '%s'", tasks[0].Schedule)
	}
	if tasks[0].Model != "claude" {
		t.Errorf("Expected model 'claude', got '%s'", tasks[0].Model)
	}
	if tasks[0].Prompt != "product_manager" {
		t.Errorf("Expected prompt 'product_manager', got '%s'", tasks[0].Prompt)
	}
	if tasks[0].Processor != "slack-pm-channel" {
		t.Errorf("Expected processor 'slack-pm-channel', got '%s'", tasks[0].Processor)
	}

	// Check the second task
	if tasks[1].Schedule != "0 9 * * 1" {
		t.Errorf("Expected schedule '0 9 * * 1', got '%s'", tasks[1].Schedule)
	}
	if tasks[1].Model != "openai" {
		t.Errorf("Expected model 'openai', got '%s'", tasks[1].Model)
	}
	if tasks[1].Prompt != "weekly_report" {
		t.Errorf("Expected prompt 'weekly_report', got '%s'", tasks[1].Prompt)
	}
	if tasks[1].Processor != "email-team@company.com" {
		t.Errorf("Expected processor 'email-team@company.com', got '%s'", tasks[1].Processor)
	}

	// Check the third task
	if tasks[2].Schedule != "*/15 * * * *" {
		t.Errorf("Expected schedule '*/15 * * * *', got '%s'", tasks[2].Schedule)
	}
	if tasks[2].Model != "gemini" {
		t.Errorf("Expected model 'gemini', got '%s'", tasks[2].Model)
	}
	if tasks[2].Prompt != "monitoring_check" {
		t.Errorf("Expected prompt 'monitoring_check', got '%s'", tasks[2].Prompt)
	}
	if tasks[2].Processor != "log-to-file" {
		t.Errorf("Expected processor 'log-to-file', got '%s'", tasks[2].Processor)
	}

	// Check the fourth task (with variables)
	if tasks[3].Schedule != "0 12 * * *" {
		t.Errorf("Expected schedule '0 12 * * *', got '%s'", tasks[3].Schedule)
	}
	if tasks[3].Model != "claude" {
		t.Errorf("Expected model 'claude', got '%s'", tasks[3].Model)
	}
	if tasks[3].Prompt != "report_template" {
		t.Errorf("Expected prompt 'report_template', got '%s'", tasks[3].Prompt)
	}
	if tasks[3].Processor != "email-execs@company.com" {
		t.Errorf("Expected processor 'email-execs@company.com', got '%s'", tasks[3].Processor)
	}

	// Check variables
	if len(tasks[3].Variables) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(tasks[3].Variables))
	}
	if tasks[3].Variables["reportType"] != "Weekly" {
		t.Errorf("Expected reportType 'Weekly', got '%s'", tasks[3].Variables["reportType"])
	}
	if tasks[3].Variables["project"] != "CronAI" {
		t.Errorf("Expected project 'CronAI', got '%s'", tasks[3].Variables["project"])
	}
	if tasks[3].Variables["team"] != "Dev" {
		t.Errorf("Expected team 'Dev', got '%s'", tasks[3].Variables["team"])
	}

	// Test non-existent config file
	_, err = parseConfigFile("non_existent_config")
	if err == nil {
		t.Error("Expected error when parsing non-existent config file, got nil")
	}
}
