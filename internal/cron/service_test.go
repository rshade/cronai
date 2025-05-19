package cron

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
	"github.com/rshade/cronai/internal/prompt"
	"github.com/rshade/cronai/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPromptManager for testing
type MockPromptManager struct {
	prompts map[string]string
	mu      sync.Mutex
}

func NewMockPromptManager() *MockPromptManager {
	return &MockPromptManager{
		prompts: make(map[string]string),
	}
}

func (m *MockPromptManager) LoadPrompt(promptName string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if content, ok := m.prompts[promptName]; ok {
		return content, nil
	}
	return "", fmt.Errorf("prompt not found: %s", promptName)
}

func (m *MockPromptManager) LoadPromptWithVariables(promptName string, variables map[string]string) (string, error) {
	content, err := m.LoadPrompt(promptName)
	if err != nil {
		return "", err
	}

	// Simple variable replacement for testing
	for key, value := range variables {
		content = fmt.Sprintf(content, value)
	}
	return content, nil
}

func (m *MockPromptManager) SetPrompt(name, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.prompts[name] = content
}

// MockProcessor for testing
type MockProcessor struct {
	processCalled bool
	lastResponse  *models.ModelResponse
	shouldFail    bool
	mu            sync.Mutex
}

func (m *MockProcessor) Process(response *models.ModelResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.processCalled = true
	m.lastResponse = response

	if m.shouldFail {
		return fmt.Errorf("mock processor error")
	}
	return nil
}

func (m *MockProcessor) GetName() string {
	return "mock"
}

func (m *MockProcessor) Validate() error {
	if m.shouldFail {
		return fmt.Errorf("mock validation error")
	}
	return nil
}

func (m *MockProcessor) SetOption(key, value string) error {
	return nil
}

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

	// Check the fourth task (with variables)
	if len(tasks[3].Variables) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(tasks[3].Variables))
	}
	if tasks[3].Variables["reportType"] != "Weekly" {
		t.Errorf("Expected reportType 'Weekly', got '%s'", tasks[3].Variables["reportType"])
	}

	// Test non-existent config file
	_, err = parseConfigFile("non_existent_config")
	if err == nil {
		t.Error("Expected error when parsing non-existent config file, got nil")
	}
}

func TestCronService_ScheduleTask(t *testing.T) {
	// Create a test service
	configFile := "test.config"
	service := &CronService{
		configFile: configFile,
		scheduler:  nil, // We'll create it in tests
		entries:    make(map[string]CronEntryMetadata),
		mu:         sync.Mutex{},
	}

	// Test scheduling with nil scheduler
	task := &processor.ScheduledTask{
		Schedule:  "* * * * *",
		Model:     "openai",
		Prompt:    "test",
		Processor: "test",
		Task: processor.Task{
			Model:     "openai",
			Prompt:    "test",
			Processor: "test",
		},
	}

	err := service.scheduleTask(task)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scheduler not initialized")
}

func TestCronService_ListTasks(t *testing.T) {
	// Create a test service with entries
	service := &CronService{
		entries: map[string]CronEntryMetadata{
			"entry1": {
				Model:     "openai",
				Prompt:    "test1",
				Processor: "processor1",
				Schedule:  "0 * * * *",
				Variables: map[string]string{"var": "value"},
			},
			"entry2": {
				Model:     "claude",
				Prompt:    "test2",
				Processor: "processor2",
				Schedule:  "30 * * * *",
				Variables: nil,
			},
		},
	}

	tasks := service.ListTasks()
	assert.Len(t, tasks, 2)

	// Verify task details
	foundEntry1 := false
	foundEntry2 := false
	for _, task := range tasks {
		if task.Model == "openai" && task.Prompt == "test1" {
			foundEntry1 = true
			assert.Equal(t, "processor1", task.Processor)
			assert.Equal(t, "0 * * * *", task.Schedule)
			assert.Equal(t, map[string]string{"var": "value"}, task.Variables)
		} else if task.Model == "claude" && task.Prompt == "test2" {
			foundEntry2 = true
			assert.Equal(t, "processor2", task.Processor)
			assert.Equal(t, "30 * * * *", task.Schedule)
			assert.Nil(t, task.Variables)
		}
	}

	assert.True(t, foundEntry1, "Entry1 not found in task list")
	assert.True(t, foundEntry2, "Entry2 not found in task list")
}

func TestCronService_executeTask(t *testing.T) {
	// Create a mock prompt manager
	mockPromptManager := NewMockPromptManager()
	mockPromptManager.SetPrompt("test_prompt", "This is a test prompt")

	// Replace the global prompt manager
	oldManager := prompt.PM
	prompt.PM = mockPromptManager
	defer func() { prompt.PM = oldManager }()

	// Create a mock processor
	mockProc := &MockProcessor{}

	// Mock the processor manager GetProcessor function
	oldGetProcessor := processor.GetProcessor
	processor.GetProcessor = func(config *processor.ProcessorConfig) (processor.Processor, error) {
		return mockProc, nil
	}
	defer func() { processor.GetProcessor = oldGetProcessor }()

	// Mock the ExecuteModel function
	oldExecuteModel := executeModel
	executeModel = func(model, prompt string, variables map[string]string, params string) (*models.ModelResponse, error) {
		return &models.ModelResponse{
			Content:   "Test response",
			Model:     model,
			Variables: variables,
		}, nil
	}
	defer func() { executeModel = oldExecuteModel }()

	// Create service
	service := &CronService{}

	// Test successful execution
	task := processor.Task{
		Model:     "openai",
		Prompt:    "test_prompt",
		Processor: "test",
		Variables: map[string]string{"key": "value"},
	}

	service.executeTask(task)

	// Verify processor was called
	assert.True(t, mockProc.processCalled)
	assert.NotNil(t, mockProc.lastResponse)
	assert.Equal(t, "Test response", mockProc.lastResponse.Content)
	assert.Equal(t, "openai", mockProc.lastResponse.Model)

	// Test execution with prompt error
	task.Prompt = "non_existent"
	service.executeTask(task)

	// Test execution with processor error
	mockProc.shouldFail = true
	task.Prompt = "test_prompt"
	service.executeTask(task)
}

func TestCronService_StartService(t *testing.T) {
	// Create test config file
	testConfigPath := "test_start_config.tmp"
	testConfigContent := `0 8 * * * claude test_prompt test-processor`

	err := os.WriteFile(testConfigPath, []byte(testConfigContent), 0644)
	require.NoError(t, err)
	defer os.Remove(testConfigPath)

	// Create prompt directory and file
	require.NoError(t, os.MkdirAll("cron_prompts", 0755))
	defer os.RemoveAll("cron_prompts")

	require.NoError(t, os.WriteFile("cron_prompts/test_prompt.md", []byte("Test content"), 0644))

	// Create service
	service := NewCronService(testConfigPath)

	// Start the service
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = service.StartService(ctx)
	assert.NoError(t, err)

	// Verify scheduler was started
	assert.NotNil(t, service.scheduler)

	// Cancel should stop the service
	cancel()

	// Give it time to clean up
	time.Sleep(50 * time.Millisecond)
}

func TestCronService_Stop(t *testing.T) {
	// Create service with scheduler
	service := &CronService{
		scheduler: make(chan struct{}), // Mock scheduler for testing
	}

	// Verify Stop works when scheduler exists
	done := make(chan bool)
	go func() {
		err := service.Stop()
		assert.NoError(t, err)
		done <- true
	}()

	// Simulate scheduler stop
	close(service.scheduler.(chan struct{}))

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Stop did not complete in time")
	}

	// Verify Stop returns error when no scheduler
	service.scheduler = nil
	err := service.Stop()
	assert.Error(t, err)
}

// Test helper functions
func TestIsValidModel(t *testing.T) {
	tests := []struct {
		model    string
		expected bool
	}{
		{"openai", true},
		{"claude", true},
		{"gemini", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			result := isValidModel(tt.model)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidProcessor(t *testing.T) {
	tests := []struct {
		processor string
		expected  bool
	}{
		{"email-test@example.com", true},
		{"slack-channel", true},
		{"webhook-https://example.com", true},
		{"file-/tmp/test.log", true},
		{"log-test", true},
		{"console", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.processor, func(t *testing.T) {
			result := isValidProcessor(tt.processor)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseVariables(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{
			input: "key1=value1,key2=value2",
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			input: "single=value",
			expected: map[string]string{
				"single": "value",
			},
		},
		{
			input:    "",
			expected: nil,
		},
		{
			input: "key=value=with=equals",
			expected: map[string]string{
				"key": "value=with=equals",
			},
		},
		{
			input:    "invalid",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseVariables(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCronService_RunTask(t *testing.T) {
	// Create a mock prompt manager
	mockPromptManager := NewMockPromptManager()
	mockPromptManager.SetPrompt("test_prompt", "This is a test prompt")

	// Replace the global prompt manager
	oldManager := prompt.PM
	prompt.PM = mockPromptManager
	defer func() { prompt.PM = oldManager }()

	// Create a mock processor
	mockProc := &MockProcessor{}

	// Mock the processor manager GetProcessor function
	oldGetProcessor := processor.GetProcessor
	processor.GetProcessor = func(config *processor.ProcessorConfig) (processor.Processor, error) {
		return mockProc, nil
	}
	defer func() { processor.GetProcessor = oldGetProcessor }()

	// Mock the ExecuteModel function
	oldExecuteModel := executeModel
	executeModel = func(model, prompt string, variables map[string]string, params string) (*models.ModelResponse, error) {
		return &models.ModelResponse{
			Content:   "Test response",
			Model:     model,
			Variables: variables,
		}, nil
	}
	defer func() { executeModel = oldExecuteModel }()

	// Create service
	service := &CronService{}

	// Test successful execution
	task := processor.Task{
		Model:     "openai",
		Prompt:    "test_prompt",
		Processor: "test",
		Variables: map[string]string{"key": "value"},
	}

	err := service.RunTask(task)
	assert.NoError(t, err)

	// Verify processor was called
	assert.True(t, mockProc.processCalled)
	assert.NotNil(t, mockProc.lastResponse)

	// Test with invalid prompt
	task.Prompt = "non_existent"
	err = service.RunTask(task)
	assert.Error(t, err)
}

func TestCronService_parseConfigLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantTask *processor.ScheduledTask
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid task without variables",
			line: "0 8 * * * claude test_prompt slack-channel",
			wantTask: &processor.ScheduledTask{
				Schedule:  "0 8 * * *",
				Model:     "claude",
				Prompt:    "test_prompt",
				Processor: "slack-channel",
				Task: processor.Task{
					Model:     "claude",
					Prompt:    "test_prompt",
					Processor: "slack-channel",
				},
			},
			wantErr: false,
		},
		{
			name: "valid task with variables",
			line: "0 9 * * * openai report email-test@example.com type=weekly,format=html",
			wantTask: &processor.ScheduledTask{
				Schedule:  "0 9 * * *",
				Model:     "openai",
				Prompt:    "report",
				Processor: "email-test@example.com",
				Variables: map[string]string{
					"type":   "weekly",
					"format": "html",
				},
				Task: processor.Task{
					Model:     "openai",
					Prompt:    "report",
					Processor: "email-test@example.com",
					Variables: map[string]string{
						"type":   "weekly",
						"format": "html",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "comment line",
			line:    "# This is a comment",
			wantErr: false,
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: false,
		},
		{
			name:    "invalid format - too few fields",
			line:    "0 8 * * *",
			wantErr: true,
			errMsg:  "invalid format",
		},
		{
			name:    "invalid model",
			line:    "0 8 * * * invalid_model test_prompt processor",
			wantErr: true,
			errMsg:  "invalid model",
		},
		{
			name:    "invalid processor",
			line:    "0 8 * * * claude test_prompt invalid_processor",
			wantErr: true,
			errMsg:  "invalid processor format",
		},
		{
			name:    "invalid variables format",
			line:    "0 8 * * * claude test_prompt email-test@example.com invalid_var",
			wantErr: true,
			errMsg:  "invalid variable format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parseConfigLine(tt.line)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.wantTask != nil {
					assert.Equal(t, tt.wantTask.Schedule, task.Schedule)
					assert.Equal(t, tt.wantTask.Model, task.Model)
					assert.Equal(t, tt.wantTask.Prompt, task.Prompt)
					assert.Equal(t, tt.wantTask.Processor, task.Processor)
					assert.Equal(t, tt.wantTask.Variables, task.Variables)
				}
			}
		})
	}
}
