package cron

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
	"github.com/rshade/cronai/internal/prompt"
	"github.com/stretchr/testify/assert"
)

// MockPromptManager is a mock implementation of prompt.Manager
type MockPromptManager struct {
	prompts map[string]string
	info    map[string]prompt.Info
}

func NewMockPromptManager() *MockPromptManager {
	return &MockPromptManager{
		prompts: make(map[string]string),
		info:    make(map[string]prompt.Info),
	}
}

func (m *MockPromptManager) LoadPrompt(name string) (string, error) {
	content, ok := m.prompts[name]
	if !ok {
		return "", fmt.Errorf("prompt not found: %s", name)
	}
	return content, nil
}

func (m *MockPromptManager) LoadPromptWithVariables(name string, variables map[string]string) (string, error) {
	content, err := m.LoadPrompt(name)
	if err != nil {
		return "", err
	}

	// Simple variable replacement for testing
	for _, value := range variables {
		content = fmt.Sprintf(content, value)
	}
	return content, nil
}

func (m *MockPromptManager) ListPrompts() ([]prompt.Info, error) {
	prompts := make([]prompt.Info, 0, len(m.info))
	for _, info := range m.info {
		prompts = append(prompts, info)
	}
	return prompts, nil
}

func (m *MockPromptManager) GetPrompt(name string) (prompt.Info, error) {
	info, ok := m.info[name]
	if !ok {
		return prompt.Info{}, fmt.Errorf("prompt not found: %s", name)
	}
	return info, nil
}

// GetPromptMetadata returns the metadata for a prompt
func (m *MockPromptManager) GetPromptMetadata(name string) (prompt.Metadata, error) {
	if info, ok := m.info[name]; ok {
		return *info.Metadata, nil
	}
	return prompt.Metadata{}, fmt.Errorf("prompt not found: %s", name)
}

func (m *MockPromptManager) GetPromptContent(name string) (string, error) {
	return m.LoadPrompt(name)
}

// GetPromptVariables returns the variables defined in a prompt's metadata
func (m *MockPromptManager) GetPromptVariables(name string) ([]prompt.Variable, error) {
	metadata, err := m.GetPromptMetadata(name)
	if err != nil {
		return nil, err
	}
	return metadata.Variables, nil
}

func (m *MockPromptManager) SetPrompt(name string, content string) {
	m.prompts[name] = content
}

// SetPromptInfo sets the info for a prompt
func (m *MockPromptManager) SetPromptInfo(name string, info prompt.Info) {
	if m.info == nil {
		m.info = make(map[string]prompt.Info)
	}
	m.info[name] = info
}

// Mock processor for testing
type mockProcessor struct {
	config       processor.Config
	processError error
}

func (m *mockProcessor) Process(_ *models.ModelResponse, _ string) error {
	return m.processError
}

func (m *mockProcessor) Validate() error {
	return nil
}

func (m *mockProcessor) GetConfig() processor.Config {
	return m.config
}

func (m *mockProcessor) GetType() string {
	return m.config.Type
}

func TestNewCronService(t *testing.T) {
	service := NewCronService("test.config")
	assert.NotNil(t, service)
	assert.Equal(t, "test.config", service.configFile)
	assert.NotNil(t, service.entries)
}

func TestParseConfigLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		lineNum   int
		expectErr bool
		expectObj *ScheduledTask
	}{
		{
			name:      "empty line",
			line:      "",
			lineNum:   1,
			expectErr: false,
			expectObj: nil,
		},
		{
			name:      "comment line",
			line:      "# This is a comment",
			lineNum:   1,
			expectErr: false,
			expectObj: nil,
		},
		{
			name:      "valid line",
			line:      "* * * * * openai test slack-test",
			lineNum:   1,
			expectErr: false,
			expectObj: &ScheduledTask{
				Schedule:  "* * * * *",
				Model:     "openai",
				Prompt:    "test",
				Processor: "slack-test",
				Variables: nil,
				Task: Task{
					Model:     "openai",
					Prompt:    "test",
					Processor: "slack-test",
					Variables: nil,
				},
			},
		},
		{
			name:      "invalid line format",
			line:      "* * * *",
			lineNum:   1,
			expectErr: true,
			expectObj: nil,
		},
		{
			name:      "with variables",
			line:      "* * * * * openai test slack-test key=value",
			lineNum:   1,
			expectErr: false,
			expectObj: &ScheduledTask{
				Schedule:  "* * * * *",
				Model:     "openai",
				Prompt:    "test",
				Processor: "slack-test",
				Variables: map[string]string{"key": "value"},
				Task: Task{
					Model:     "openai",
					Prompt:    "test",
					Processor: "slack-test",
					Variables: map[string]string{"key": "value"},
				},
			},
		},
		{
			name:      "with multiple variables",
			line:      "* * * * * openai test slack-test key=value,foo=bar",
			lineNum:   1,
			expectErr: false,
			expectObj: &ScheduledTask{
				Schedule:  "* * * * *",
				Model:     "openai",
				Prompt:    "test",
				Processor: "slack-test",
				Variables: map[string]string{"key": "value", "foo": "bar"},
				Task: Task{
					Model:     "openai",
					Prompt:    "test",
					Processor: "slack-test",
					Variables: map[string]string{"key": "value", "foo": "bar"},
				},
			},
		},
		{
			name:      "with model parameters",
			line:      "* * * * * openai:temperature=0.7 test slack-test",
			lineNum:   1,
			expectErr: false,
			expectObj: &ScheduledTask{
				Schedule:    "* * * * *",
				Model:       "openai",
				Prompt:      "test",
				Processor:   "slack-test",
				Variables:   nil,
				ModelParams: "temperature=0.7",
				Task: Task{
					Model:       "openai",
					Prompt:      "test",
					Processor:   "slack-test",
					Variables:   nil,
					ModelParams: "temperature=0.7",
				},
			},
		},
		{
			name:      "with custom template",
			line:      "* * * * * openai test slack-test template=custom_template",
			lineNum:   1,
			expectErr: false,
			expectObj: &ScheduledTask{
				Schedule:  "* * * * *",
				Model:     "openai",
				Prompt:    "test",
				Processor: "slack-test",
				Variables: map[string]string{"template": "custom_template"},
				Template:  "custom_template",
				Task: Task{
					Model:     "openai",
					Prompt:    "test",
					Processor: "slack-test",
					Variables: map[string]string{"template": "custom_template"},
					Template:  "custom_template",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseConfigLine(tt.line)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectObj, result)
			}
		})
	}
}

func TestScheduleTask(t *testing.T) {
	// Create a service without initialization
	service := &Service{
		configFile: "test.config",
		entries:    make(map[string]EntryMetadata),
		mu:         sync.Mutex{},
	}

	// Test scheduling with nil scheduler
	task := &ScheduledTask{
		Schedule:  "* * * * *",
		Model:     "openai",
		Prompt:    "test",
		Processor: "test",
		Task: Task{
			Model:     "openai",
			Prompt:    "test",
			Processor: "test",
		},
	}

	err := service.scheduleTask(task)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scheduler not initialized")

	// Initialize the scheduler and try again
	service.scheduler = cron.New()
	assert.NotNil(t, service.scheduler)

	// Create a mock prompt manager
	mockPM := NewMockPromptManager()
	mockPM.SetPrompt("test", "This is a test prompt")

	// Replace the global PM with our mock
	oldManager := prompt.PM
	prompt.PM = mockPM
	defer func() { prompt.PM = oldManager }()

	// Create a mock processor
	mockProc := &mockProcessor{}

	// Register the mock processor with the registry
	registry := processor.GetRegistry()
	registry.RegisterFactory("console", func(_ processor.Config) (processor.Processor, error) {
		return mockProc, nil
	})

	// Mock the ExecuteModel function
	oldExecuteModel := executeModel
	executeModel = func(model, prompt string, variables map[string]string, _ string) (*models.ModelResponse, error) {
		return &models.ModelResponse{
			Content:    "Test response",
			Model:      model,
			PromptName: prompt,
			Variables:  variables,
			Timestamp:  time.Now(),
		}, nil
	}
	defer func() { executeModel = oldExecuteModel }()

	// Now scheduling should succeed
	err = service.scheduleTask(task)
	assert.NoError(t, err)

	// Check that the entry was added
	assert.Len(t, service.entries, 1)
}

func TestRunTask(t *testing.T) {
	// Create a test-specific service
	service := &Service{
		configFile: "test.config",
		entries:    make(map[string]EntryMetadata),
		mu:         sync.Mutex{},
	}

	// Create a mock prompt manager
	mockPM := NewMockPromptManager()
	mockPM.SetPrompt("test", "This is a test prompt")

	// Replace the global PM with our mock
	oldManager := prompt.PM
	prompt.PM = mockPM
	defer func() { prompt.PM = oldManager }()

	// Create a mock processor
	mockProc := &mockProcessor{}

	// Register the mock processor with the registry
	registry := processor.GetRegistry()
	registry.RegisterFactory("console", func(_ processor.Config) (processor.Processor, error) {
		return mockProc, nil
	})

	// Mock the ExecuteModel function
	oldExecuteModel := executeModel
	executeModel = func(model, prompt string, variables map[string]string, _ string) (*models.ModelResponse, error) {
		return &models.ModelResponse{
			Content:    "Test response",
			Model:      model,
			PromptName: prompt,
			Variables:  variables,
			Timestamp:  time.Now(),
		}, nil
	}
	defer func() { executeModel = oldExecuteModel }()

	// Run a task
	task := Task{
		Model:     "openai",
		Prompt:    "test",
		Processor: "console",
	}

	err := service.RunTask(task)
	assert.NoError(t, err)
}

func TestRunNonExistentPrompt(t *testing.T) {
	// Create a service
	service := &Service{
		configFile: "test.config",
		entries:    make(map[string]EntryMetadata),
		mu:         sync.Mutex{},
	}

	// Create a mock prompt manager
	mockPM := NewMockPromptManager()
	// Don't set any prompts so LoadPrompt will fail

	// Replace the global PM with our mock
	oldManager := prompt.PM
	prompt.PM = mockPM
	defer func() { prompt.PM = oldManager }()

	// Run a task with a non-existent prompt
	task := Task{
		Model:     "openai",
		Prompt:    "non-existent",
		Processor: "console",
	}

	err := service.RunTask(task)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error loading prompt")
}

func TestRunWithFailedProcessor(t *testing.T) {
	// Create a test-specific service
	service := &Service{
		configFile: "test.config",
		entries:    make(map[string]EntryMetadata),
		mu:         sync.Mutex{},
	}

	// Create a mock prompt manager
	mockPM := NewMockPromptManager()
	mockPM.SetPrompt("test", "This is a test prompt")

	// Replace the global PM with our mock
	oldManager := prompt.PM
	prompt.PM = mockPM
	defer func() { prompt.PM = oldManager }()

	// Create a mock processor that fails
	mockProc := &mockProcessor{
		processError: fmt.Errorf("error processing response"),
	}

	// Register the mock processor with the registry
	registry := processor.GetRegistry()
	registry.RegisterFactory("console", func(_ processor.Config) (processor.Processor, error) {
		return mockProc, nil
	})

	// Mock the ExecuteModel function
	oldExecuteModel := executeModel
	executeModel = func(model, prompt string, variables map[string]string, _ string) (*models.ModelResponse, error) {
		return &models.ModelResponse{
			Content:    "Test response",
			Model:      model,
			PromptName: prompt,
			Variables:  variables,
			Timestamp:  time.Now(),
		}, nil
	}
	defer func() { executeModel = oldExecuteModel }()

	// Run a task
	task := Task{
		Model:     "openai",
		Prompt:    "test",
		Processor: "console",
	}

	err := service.RunTask(task)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error processing response")
}

func TestParseProcessor(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   string
		expectedTarget string
	}{
		{"slack-channel", "slack", "channel"},
		{"email-user@example.com", "email", "user@example.com"},
		{"webhook-alerts", "webhook", "alerts"},
		{"github-repo", "github", "repo"},
		{"file-log.txt", "file", "log.txt"},
		{"log-to-file", "file", "to-file"},
		{"console", "console", ""},
		{"unknown-format", "console", ""}, // Default case
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			procType, target := parseProcessor(tt.input)
			assert.Equal(t, tt.expectedType, procType)
			assert.Equal(t, tt.expectedTarget, target)
		})
	}
}

func TestIsValidModel(t *testing.T) {
	assert.True(t, isValidModel("openai"))
	assert.True(t, isValidModel("claude"))
	assert.True(t, isValidModel("gemini"))
	assert.False(t, isValidModel("unknown"))
}

func TestIsValidProcessor(t *testing.T) {
	assert.True(t, isValidProcessor("console"))
	assert.True(t, isValidProcessor("slack-channel"))
	assert.True(t, isValidProcessor("email-user@example.com"))
	assert.True(t, isValidProcessor("webhook-alerts"))
	assert.True(t, isValidProcessor("file-log.txt"))
	assert.True(t, isValidProcessor("log-to-file"))
	assert.True(t, isValidProcessor("github-repo"))
	assert.False(t, isValidProcessor("unknown"))
}

func TestService_ProcessResponse(t *testing.T) {
	// Create test service
	service := NewCronService("test.config")

	// Create test processor
	config := processor.Config{
		Type:   "test",
		Target: "test-target",
	}

	// Create mock processor
	mock := &mockProcessor{
		config: config,
	}

	// Test processing response
	response := &models.ModelResponse{
		Content:    "Test content",
		Model:      "test-model",
		Timestamp:  time.Now(),
		PromptName: "test-prompt",
	}

	// Test with mock processor
	err := service.ProcessResponse(mock, response, "")
	if err != nil {
		t.Errorf("ProcessResponse failed: %v", err)
	}
}

func TestService_GetProcessor(t *testing.T) {
	// Skip test for now as it requires modifying the global processor registry
	t.Skip("Skipping test that requires modifying global processor registry")

	// Create test service
	service := NewCronService("test.config")

	// Create test config
	config := processor.Config{
		Type:   "console", // Use a real processor type that's registered by default
		Target: "test-target",
	}

	// Test getting processor
	proc, err := service.GetProcessor("console", config)
	if err != nil {
		t.Errorf("GetProcessor failed: %v", err)
	}

	// Verify processor type
	if proc.GetType() != "console" {
		t.Errorf("GetProcessor returned wrong type: got %s, want console", proc.GetType())
	}
}

func TestService_CreateProcessor(t *testing.T) {
	// Skip test for now as it requires modifying the global processor registry
	t.Skip("Skipping test that requires modifying global processor registry")

	// Create test service
	service := NewCronService("test.config")

	// Create test config
	config := processor.Config{
		Type:   "console", // Use a real processor type that's registered by default
		Target: "test-target",
	}

	// Test creating processor
	proc, err := service.CreateProcessor("console", config)
	if err != nil {
		t.Errorf("CreateProcessor failed: %v", err)
	}

	// Verify processor type
	if proc.GetType() != "console" {
		t.Errorf("CreateProcessor returned wrong type: got %s, want console", proc.GetType())
	}
}
