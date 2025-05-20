package template

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Functions that need to be added to getTemplateFuncMap()

// JoinFunction joins items with a separator
func JoinFunction(items interface{}, sep string) string {
	switch v := items.(type) {
	case []string:
		return strings.Join(v, sep)
	case string:
		// Convert space-separated string to array and join
		parts := strings.Fields(v)
		return strings.Join(parts, sep)
	default:
		return fmt.Sprintf("%v", items)
	}
}

// ReplaceFunction replaces old with new in a string
func ReplaceFunction(s, old, newStr string) string {
	return strings.ReplaceAll(s, old, newStr)
}

// JSONFunction marshals data to a JSON string
func JSONFunction(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// DefaultFunction returns default value if input is empty
func DefaultFunction(value, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	// Check for empty string
	if str, ok := value.(string); ok && str == "" {
		return defaultValue
	}
	// Check for zero values of various types
	val := reflect.ValueOf(value)
	if val.IsZero() {
		return defaultValue
	}
	return value
}

// CreateExecutionContext creates a context that merges Data and Variables
func CreateExecutionContext(data Data) map[string]interface{} {
	execContext := make(map[string]interface{})

	// Add all Data fields
	execContext["Content"] = data.Content
	execContext["Model"] = data.Model
	execContext["Timestamp"] = data.Timestamp
	execContext["PromptName"] = data.PromptName
	execContext["Variables"] = data.Variables
	execContext["ExecutionID"] = data.ExecutionID
	execContext["Metadata"] = data.Metadata
	execContext["Parent"] = data.Parent

	// Merge variables into the top level for direct access
	if data.Variables != nil {
		for k, v := range data.Variables {
			execContext[k] = v
		}
	}

	return execContext
}

// PreprocessBlockSyntax transforms simple block syntax to Go's template syntax
func PreprocessBlockSyntax(content string) string {
	// Replace {{block "name"}} with {{block "name" .}}
	blockPattern := regexp.MustCompile(`(?s)\{\{\s*block\s+"([^"]+)"\s*\}\}`)
	content = blockPattern.ReplaceAllString(content, `{{block "$1" .}}`)

	// Replace {{endblock}} with {{end}}
	endBlockPattern := regexp.MustCompile(`(?s)\{\{\s*endblock\s*\}\}`)
	content = endBlockPattern.ReplaceAllString(content, `{{end}}`)

	return content
}
