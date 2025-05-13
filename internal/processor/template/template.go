package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode"
)

// TemplateData contains data available to templates
type TemplateData struct {
	Content     string            // Model response content
	Model       string            // Model name
	Timestamp   time.Time         // Response timestamp
	PromptName  string            // Name of the prompt
	Variables   map[string]string // Custom variables
	ExecutionID string            // Unique execution identifier
	Metadata    map[string]string // Additional metadata
}

// Manager handles template operations
type Manager struct {
	templates map[string]*template.Template
	mutex     sync.RWMutex
}

// singleton instance
var (
	instance *Manager
	once     sync.Once
)

// GetManager returns the singleton template manager
func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			templates: make(map[string]*template.Template),
		}
		// Register default templates
		instance.registerDefaultTemplates()
	})
	return instance
}

// getTemplateFuncMap returns the function map for templates
func getTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		// Variable existence check
		"hasVar": func(v map[string]string, key string) bool {
			_, exists := v[key]
			return exists
		},
		// Variable access with default value
		"getVar": func(v map[string]string, key, defaultVal string) string {
			if val, exists := v[key]; exists {
				return val
			}
			return defaultVal
		},
		// String comparison operators
		"eq":        func(a, b string) bool { return a == b },
		"ne":        func(a, b string) bool { return a != b },
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		// Boolean operators
		"not": func(b bool) bool { return !b },
		// Numeric comparison (converts strings to numbers first)
		"lt": func(a, b string) bool {
			aVal, aErr := strconv.ParseFloat(a, 64)
			bVal, bErr := strconv.ParseFloat(b, 64)
			if aErr != nil || bErr != nil {
				return a < b // Fallback to string comparison
			}
			return aVal < bVal
		},
		"le": func(a, b string) bool {
			aVal, aErr := strconv.ParseFloat(a, 64)
			bVal, bErr := strconv.ParseFloat(b, 64)
			if aErr != nil || bErr != nil {
				return a <= b // Fallback to string comparison
			}
			return aVal <= bVal
		},
		"gt": func(a, b string) bool {
			aVal, aErr := strconv.ParseFloat(a, 64)
			bVal, bErr := strconv.ParseFloat(b, 64)
			if aErr != nil || bErr != nil {
				return a > b // Fallback to string comparison
			}
			return aVal > bVal
		},
		"ge": func(a, b string) bool {
			aVal, aErr := strconv.ParseFloat(a, 64)
			bVal, bErr := strconv.ParseFloat(b, 64)
			if aErr != nil || bErr != nil {
				return a >= b // Fallback to string comparison
			}
			return aVal >= bVal
		},
		// JSON utilities
		"marshalJSON": func(s string) string {
			bytes, err := json.Marshal(s)
			if err != nil {
				return fmt.Sprintf("\"%s\"", s) // Fallback
			}
			return string(bytes)
		},
		// Map utilities
		"isLast": func(m map[string]string, key string) bool {
			// Get all keys
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}

			// Sort keys for consistent behavior
			sort.Strings(keys)

			// Check if this is the last key
			return key == keys[len(keys)-1]
		},
		// Date utilities
		"now": time.Now,
		"formatDate": func(format string, t time.Time) string {
			return t.Format(format)
		},
		"addDays": func(days int, t time.Time) time.Time {
			return t.AddDate(0, 0, days)
		},
		// String utilities
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": func(s string) string {
			// Simple title implementation that capitalizes first letter of each word
			prev := ' '
			return strings.Map(
				func(r rune) rune {
					result := r
					if prev == ' ' || prev == '\t' || prev == '\n' || prev == '\r' {
						result = unicode.ToTitle(r)
					}
					prev = r
					return result
				},
				s,
			)
		},
		"trim": strings.TrimSpace,
	}
}

// RegisterTemplate adds or updates a template
func (m *Manager) RegisterTemplate(name, content string) error {
	// Parse the template to validate it
	tmpl, err := template.New(name).
		Funcs(getTemplateFuncMap()).
		Parse(content)
	if err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	// Add to the template map
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.templates[name] = tmpl
	return nil
}

// GetTemplate retrieves a template by name
func (m *Manager) GetTemplate(name string) (*template.Template, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	tmpl, exists := m.templates[name]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", name)
	}
	return tmpl, nil
}

// Execute applies a template with the given data
func (m *Manager) Execute(name string, data TemplateData) (string, error) {
	tmpl, err := m.GetTemplate(name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}

// SafeExecute attempts to execute a template with fallbacks
func (m *Manager) SafeExecute(name string, data TemplateData) string {
	// Try primary template
	result, err := m.Execute(name, data)
	if err == nil {
		return result
	}

	// Try fallback template
	parts := strings.Split(name, "_")
	if len(parts) > 1 {
		fallbackName := "default_" + parts[0]
		if fallbackResult, fallbackErr := m.Execute(fallbackName, data); fallbackErr == nil {
			return fallbackResult
		}
	}

	// Last resort: just return the raw content
	return data.Content
}

// Validate checks if a template is valid
func (m *Manager) Validate(content string) error {
	_, err := template.New("validation").
		Funcs(getTemplateFuncMap()).
		Parse(content)
	return err
}

// LoadTemplatesFromDir loads templates from a directory
func (m *Manager) LoadTemplatesFromDir(directory string) error {
	// Check if directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return fmt.Errorf("template directory does not exist: %s", directory)
	}

	// Find template files
	files, err := filepath.Glob(filepath.Join(directory, "*.tmpl"))
	if err != nil {
		return fmt.Errorf("failed to glob template directory: %w", err)
	}

	// Load each template file
	for _, file := range files {
		// Extract template name from filename
		name := strings.TrimSuffix(filepath.Base(file), ".tmpl")

		// Read template content
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", file, err)
		}

		// Register template
		if err := m.RegisterTemplate(name, string(content)); err != nil {
			return fmt.Errorf("failed to register template %s: %w", name, err)
		}
	}

	return nil
}

// TemplateExists checks if a template with the given name exists
func (m *Manager) TemplateExists(name string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.templates[name]
	return exists
}

// ValidateTemplate validates a template file
func (m *Manager) ValidateTemplate(filePath string) error {
	// Read template content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", filePath, err)
	}

	// Validate the template
	return m.Validate(string(content))
}

// ValidateTemplatesInDir validates all templates in a directory
func (m *Manager) ValidateTemplatesInDir(directory string) (map[string]error, error) {
	// Check if directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil, fmt.Errorf("template directory does not exist: %s", directory)
	}

	// Find template files
	files, err := filepath.Glob(filepath.Join(directory, "*.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob template directory: %w", err)
	}

	// Validate each template file
	results := make(map[string]error)
	for _, file := range files {
		name := filepath.Base(file)
		err := m.ValidateTemplate(file)
		results[name] = err
	}

	return results, nil
}
