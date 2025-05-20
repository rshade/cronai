package template

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateExecutionContext(t *testing.T) {
	// Create test data
	testTime := time.Date(2024, 3, 15, 12, 30, 45, 0, time.UTC)
	data := Data{
		Content:     "Test content",
		Model:       "test-model",
		Timestamp:   testTime,
		PromptName:  "test-prompt",
		ExecutionID: "exec-123",
		Variables: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Metadata: map[string]string{
			"meta1": "metadata1",
		},
		Parent: "parent-template",
	}

	// Create execution context
	context := CreateExecutionContext(data)

	// Verify all Data fields are in the context
	assert.Equal(t, "Test content", context["Content"])
	assert.Equal(t, "test-model", context["Model"])
	assert.Equal(t, testTime, context["Timestamp"])
	assert.Equal(t, "test-prompt", context["PromptName"])
	assert.Equal(t, "exec-123", context["ExecutionID"])
	assert.Equal(t, data.Variables, context["Variables"])
	assert.Equal(t, data.Metadata, context["Metadata"])
	assert.Equal(t, "parent-template", context["Parent"])

	// Verify Variables are merged into the top level
	assert.Equal(t, "value1", context["key1"])
	assert.Equal(t, "value2", context["key2"])
}

func TestJoinFunction(t *testing.T) {
	tests := []struct {
		name     string
		items    interface{}
		sep      string
		expected string
	}{
		{
			name:     "string slice",
			items:    []string{"a", "b", "c"},
			sep:      ",",
			expected: "a,b,c",
		},
		{
			name:     "space-separated string",
			items:    "a b c",
			sep:      "-",
			expected: "a-b-c",
		},
		{
			name:     "other type",
			items:    123,
			sep:      ",",
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinFunction(tt.items, tt.sep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReplaceFunction(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		old      string
		new      string
		expected string
	}{
		{
			name:     "simple replacement",
			s:        "hello world",
			old:      "world",
			new:      "universe",
			expected: "hello universe",
		},
		{
			name:     "multiple replacements",
			s:        "hello hello hello",
			old:      "hello",
			new:      "hi",
			expected: "hi hi hi",
		},
		{
			name:     "no match",
			s:        "hello world",
			old:      "foo",
			new:      "bar",
			expected: "hello world",
		},
		{
			name:     "empty string",
			s:        "",
			old:      "foo",
			new:      "bar",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceFunction(tt.s, tt.old, tt.new)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJSONFunction(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "map",
			value:    map[string]string{"key": "value"},
			expected: `{"key":"value"}`,
			wantErr:  false,
		},
		{
			name:     "slice",
			value:    []string{"a", "b", "c"},
			expected: `["a","b","c"]`,
			wantErr:  false,
		},
		{
			name:     "string",
			value:    "hello",
			expected: `"hello"`,
			wantErr:  false,
		},
		{
			name:     "number",
			value:    123,
			expected: `123`,
			wantErr:  false,
		},
		{
			name:     "boolean",
			value:    true,
			expected: `true`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := JSONFunction(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDefaultFunction(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		defValue interface{}
		expected interface{}
	}{
		{
			name:     "non-empty string",
			value:    "hello",
			defValue: "default",
			expected: "hello",
		},
		{
			name:     "empty string",
			value:    "",
			defValue: "default",
			expected: "default",
		},
		{
			name:     "nil value",
			value:    nil,
			defValue: "default",
			expected: "default",
		},
		{
			name:     "zero integer",
			value:    0,
			defValue: 10,
			expected: 10,
		},
		{
			name:     "non-zero integer",
			value:    5,
			defValue: 10,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultFunction(tt.value, tt.defValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPreprocessBlockSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "block without dot",
			input:    `{{block "header"}}Content{{end}}`,
			expected: `{{block "header" .}}Content{{end}}`,
		},
		{
			name:     "endblock directive",
			input:    `{{block "header" .}}Content{{endblock}}`,
			expected: `{{block "header" .}}Content{{end}}`,
		},
		{
			name:     "multiple blocks",
			input:    `{{block "header"}}Header{{endblock}}{{block "footer"}}Footer{{endblock}}`,
			expected: `{{block "header" .}}Header{{end}}{{block "footer" .}}Footer{{end}}`,
		},
		{
			name:     "with spaces",
			input:    `{{ block "header" }}Content{{ endblock }}`,
			expected: `{{block "header" .}}Content{{end}}`,
		},
		{
			name:     "no changes needed",
			input:    `{{block "header" .}}Content{{end}}`,
			expected: `{{block "header" .}}Content{{end}}`,
		},
		{
			name:     "no blocks",
			input:    `Just a regular template`,
			expected: `Just a regular template`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PreprocessBlockSyntax(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
