package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode"
)

// Removed unused templateBlock struct

// templateInheritance tracks parent-child relationships between templates
type templateInheritance struct {
	Parent string
	Blocks map[string]string
}

// TemplateData contains data available to templates
type TemplateData struct {
	Content     string            // Model response content
	Model       string            // Model name
	Timestamp   time.Time         // Response timestamp
	PromptName  string            // Name of the prompt
	Variables   map[string]string // Custom variables
	ExecutionID string            // Unique execution identifier
	Metadata    map[string]string // Additional metadata
	Parent      interface{}       // Parent template data for inheritance
}

// Manager handles template operations
type Manager struct {
	templates   map[string]*template.Template
	inheritance map[string]*templateInheritance
	mutex       sync.RWMutex
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
			templates:   make(map[string]*template.Template),
			inheritance: make(map[string]*templateInheritance),
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
		// Template inheritance and composition functions
		"block": func(name string, content string) (string, error) {
			// This function is a placeholder for defining blocks
			// The actual implementation is handled by the template engine
			return content, nil
		},
		"extends": func(name string) (string, error) {
			// This function is a placeholder for template inheritance
			// The actual implementation is handled during template processing
			return "", nil
		},
		"super": func() (string, error) {
			// This function allows accessing parent block content
			return "", nil
		},
		"include": func(name string) (string, error) {
			// This function allows including other templates
			// The actual implementation is handled during template processing
			return "", nil
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

// ParseInheritance parses template content for inheritance directives
func (m *Manager) ParseInheritance(name, content string) (*templateInheritance, string, error) {
	// Check for {{extends "parent_template"}} directive
	extendsPattern := regexp.MustCompile(`(?m)\{\{\s*extends\s+"([^"]+)"\s*\}\}`)
	match := extendsPattern.FindStringSubmatch(content)

	if len(match) < 2 {
		// No inheritance found
		return nil, content, nil
	}

	parentName := match[1]
	contentWithoutExtends := extendsPattern.ReplaceAllString(content, "")

	// Extract blocks from the template - try both define and block patterns
	// First, look for {{define "name"}}content{{end}} style
	definePattern := regexp.MustCompile(`(?s)\{\{\s*define\s+"([^"]+)"\s*\}\}(.*?)\{\{\s*end\s*\}\}`)
	defineMatches := definePattern.FindAllStringSubmatch(contentWithoutExtends, -1)

	// Also look for {{block "name" .}}content{{end}} style for compatibility
	blockPattern := regexp.MustCompile(`(?s)\{\{\s*block\s+"([^"]+)"\s+\.\s*\}\}(.*?)\{\{\s*end\s*\}\}`)
	blockMatches := blockPattern.FindAllStringSubmatch(contentWithoutExtends, -1)

	// Combine results from both patterns
	blocks := make(map[string]string)

	// Process define blocks first
	for _, defineMatch := range defineMatches {
		if len(defineMatch) >= 3 {
			blockName := defineMatch[1]
			blockContent := defineMatch[2]
			blocks[blockName] = blockContent
		}
	}

	// Then process legacy block style (will overwrite if duplicates)
	for _, blockMatch := range blockMatches {
		if len(blockMatch) >= 3 {
			blockName := blockMatch[1]
			blockContent := blockMatch[2]
			blocks[blockName] = blockContent
		}
	}

	return &templateInheritance{
		Parent: parentName,
		Blocks: blocks,
	}, contentWithoutExtends, nil
}

// ProcessInheritance processes a template with inheritance
func (m *Manager) ProcessInheritance(name string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check if this template inherits from another
	inheritance, exists := m.inheritance[name]
	if !exists {
		// No inheritance, just return the template as is
		_, err := m.GetTemplate(name)
		if err != nil {
			return "", err
		}
		return name, nil
	}

	// Get the parent template
	parentName := inheritance.Parent
	_, err := m.GetTemplate(parentName)
	if err != nil {
		return "", fmt.Errorf("parent template not found: %s", parentName)
	}

	// Check if the parent itself has inheritance (handle nested inheritance)
	_, parentHasInheritance := m.inheritance[parentName]
	if parentHasInheritance {
		// Debug logging removed

		// Process parent's inheritance first
		_, err = m.ProcessInheritance(parentName)
		if err != nil {
			return "", fmt.Errorf("failed to process parent inheritance: %w", err)
		}

		// Process parent as a template with includes
		m.mutex.RUnlock()
		if err := m.RegisterTemplateWithIncludes(parentName, ""); err != nil {
			m.mutex.RLock()
			return "", fmt.Errorf("failed to process parent includes: %w", err)
		}
		m.mutex.RLock()

		// Re-get parent inheritance after processing
		parentInheritance, parentHasInheritance := m.inheritance[parentName]

		if parentHasInheritance && parentInheritance != nil {
			// Merge blocks from grandparent that weren't overridden
			// Debug logging removed
			// Update blocks that weren't overridden in the child template
			for blockName := range parentInheritance.Blocks {
				if _, overridden := inheritance.Blocks[blockName]; !overridden {
					// This block wasn't overridden in the child, so it should inherit from parent
					parentBlockContent, exists := parentInheritance.Blocks[blockName]
					if exists {
						// Ensure child inheritance includes all parent blocks
						m.mutex.RUnlock()
						m.mutex.Lock()
						if inheritance.Blocks == nil {
							inheritance.Blocks = make(map[string]string)
						}
						inheritance.Blocks[blockName] = parentBlockContent
						m.mutex.Unlock()
						m.mutex.RLock()
					}
				}
			}
		}
	}

	return name, nil
}

// RegisterTemplate adds or updates a template
func (m *Manager) RegisterTemplate(name, content string) error {
	// Check for template inheritance
	inheritance, processedContent, err := m.ParseInheritance(name, content)
	if err != nil {
		return fmt.Errorf("failed to parse inheritance: %w", err)
	}

	// Create the base template
	tmpl := template.New(name).Funcs(getTemplateFuncMap())

	// If this template extends another, store the inheritance information
	if inheritance != nil && inheritance.Parent != "" {
		m.mutex.Lock()
		m.inheritance[name] = inheritance
		m.mutex.Unlock()

		// Check if parent template exists
		parentName := inheritance.Parent
		parent, err := m.GetTemplate(parentName)
		if err == nil {
			// If parent exists, clone its parse tree into this template
			for _, t := range parent.Templates() {
				if t.Name() != parentName {
					// Clone all associated templates except the root one
					_, err = tmpl.AddParseTree(t.Name(), t.Tree)
					if err != nil {
						return fmt.Errorf("failed to add parse tree from parent template %s: %w", t.Name(), err)
					}
				}
			}
		}

		// Define block templates explicitly
		for blockName, blockContent := range inheritance.Blocks {
			blockTemplate := fmt.Sprintf(`{{define "%s"}}%s{{end}}`, blockName, blockContent)
			_, err = tmpl.Parse(blockTemplate)
			if err != nil {
				return fmt.Errorf("failed to parse block %s: %w", blockName, err)
			}
		}
	}

	// Parse the main template content
	_, err = tmpl.Parse(processedContent)
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

// ExecuteTemplate executes a named template within a template file
func (m *Manager) ExecuteTemplate(name, templateName string, data TemplateData) (string, error) {
	// Get the template
	tmpl, err := m.GetTemplate(name)
	if err != nil {
		return "", err
	}

	// Create a new template for execution that combines all templates
	execTmpl := template.New("exec_" + name).Funcs(getTemplateFuncMap())

	// Add the main template and all its associated templates
	for _, t := range tmpl.Templates() {
		// Debug logging removed
		_, err = execTmpl.AddParseTree(t.Name(), t.Tree)
		if err != nil {
			return "", fmt.Errorf("failed to add template %s: %w", t.Name(), err)
		}
	}

	// Check if this template extends another
	m.mutex.RLock()
	inheritance, hasInheritance := m.inheritance[name]
	m.mutex.RUnlock()

	if hasInheritance {
		parentName := inheritance.Parent
		parentTmpl, err := m.GetTemplate(parentName)
		if err != nil {
			return "", fmt.Errorf("parent template not found: %w", err)
		}

		// Add all templates from the parent
		for _, t := range parentTmpl.Templates() {
			if execTmpl.Lookup(t.Name()) == nil {
				// Debug logging removed
				_, err = execTmpl.AddParseTree(t.Name(), t.Tree)
				if err != nil {
					return "", fmt.Errorf("failed to add parent template %s: %w", t.Name(), err)
				}
			}
		}
	}

	// Add all other templates too for completeness
	m.mutex.RLock()
	for tName, t := range m.templates {
		if tName != name {
			for _, subT := range t.Templates() {
				if execTmpl.Lookup(subT.Name()) == nil {
					// Debug logging removed
					_, err = execTmpl.AddParseTree(subT.Name(), subT.Tree)
					if err != nil {
						m.mutex.RUnlock()
						return "", fmt.Errorf("failed to add template %s: %w", subT.Name(), err)
					}
				}
			}
		}
	}
	m.mutex.RUnlock()

	// Execute the specific named template
	var buf bytes.Buffer
	if err := execTmpl.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	result := buf.String()
	// Debug logging removed
	return result, nil
}

// Execute applies a template with the given data
func (m *Manager) Execute(name string, data TemplateData) (string, error) {
	// Check if this template extends another
	m.mutex.RLock()
	inheritance, hasInheritance := m.inheritance[name]
	m.mutex.RUnlock()

	if hasInheritance {
		// If this template extends another, we need special handling
		parentName := inheritance.Parent

		// Get the parent template
		parentTmpl, err := m.GetTemplate(parentName)
		if err != nil {
			return "", fmt.Errorf("parent template not found: %w", err)
		}

		// Get the child template too
		childTmpl, err := m.GetTemplate(name)
		if err != nil {
			return "", fmt.Errorf("child template not found: %w", err)
		}

		// Create a completely new template for execution
		execTmpl := template.New("exec_" + name).Funcs(getTemplateFuncMap())

		// First, add all the define blocks from the child template to override parent
		for blockName, blockContent := range inheritance.Blocks {
			// Debug logging removed
			blockDef := fmt.Sprintf(`{{define "%s"}}%s{{end}}`, blockName, blockContent)
			_, err = execTmpl.Parse(blockDef)
			if err != nil {
				return "", fmt.Errorf("failed to parse block %s: %w", blockName, err)
			}
		}

		// Next, add all templates from the parent (except the main one)
		for _, t := range parentTmpl.Templates() {
			if t.Name() != parentName && t.Name() != execTmpl.Name() {
				// Skip if already defined in our template (from block overrides)
				if execTmpl.Lookup(t.Name()) == nil {
					// Debug logging removed
					_, err = execTmpl.AddParseTree(t.Name(), t.Tree)
					if err != nil {
						return "", fmt.Errorf("failed to add template %s from parent: %w", t.Name(), err)
					}
				}
			}
		}

		// Add all other child templates
		for _, t := range childTmpl.Templates() {
			if t.Name() != name && t.Name() != execTmpl.Name() &&
				!strings.HasPrefix(t.Name(), "exec_") {
				// Skip if already defined (from block overrides)
				if execTmpl.Lookup(t.Name()) == nil {
					// Debug logging removed
					_, err = execTmpl.AddParseTree(t.Name(), t.Tree)
					if err != nil {
						return "", fmt.Errorf("failed to add template %s from child: %w", t.Name(), err)
					}
				}
			}
		}

		// Finally, add the parent's main template (which we'll execute)
		_, err = execTmpl.AddParseTree(parentName, parentTmpl.Tree)
		if err != nil {
			return "", fmt.Errorf("failed to add parent template: %w", err)
		}

		// List all available templates for debugging
		// Debug logging removed

		// Execute the parent template (which will use our overridden blocks)
		var buf bytes.Buffer
		if err := execTmpl.ExecuteTemplate(&buf, parentName, data); err != nil {
			return "", fmt.Errorf("template execution failed: %w", err)
		}

		result := buf.String()
		// Debug logging removed
		return result, nil
	}

	// For non-inherited templates, execute normally
	tmpl, err := m.GetTemplate(name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	// Always handle templates properly regardless of whether they have includes
	// Debug logging removed

	// Create a new template with all known templates
	execTmpl := template.New(name).Funcs(getTemplateFuncMap())

	// Add the main template
	_, err = execTmpl.AddParseTree(name, tmpl.Tree)
	if err != nil {
		return "", fmt.Errorf("failed to add main template: %w", err)
	}

	// Add all associated templates defined in this template
	for _, t := range tmpl.Templates() {
		if t.Name() != name {
			_, err = execTmpl.AddParseTree(t.Name(), t.Tree)
			if err != nil {
				return "", fmt.Errorf("failed to add template %s: %w", t.Name(), err)
			}
		}
	}

	// Also add all other related templates that might be referenced
	// This is especially important for includes
	m.mutex.RLock()
	for tName, t := range m.templates {
		if tName != name && execTmpl.Lookup(tName) == nil {
			for _, subT := range t.Templates() {
				if execTmpl.Lookup(subT.Name()) == nil {
					// Debug logging removed
					_, err = execTmpl.AddParseTree(subT.Name(), subT.Tree)
					if err != nil {
						m.mutex.RUnlock()
						return "", fmt.Errorf("failed to add related template %s: %w", subT.Name(), err)
					}
				}
			}
		}
	}
	m.mutex.RUnlock()

	// Print available templates for debugging
	// Debug logging removed

	// Try to execute the template
	if err := execTmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	result := buf.String()
	// Debug logging removed
	return result, nil
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

		// Register template with includes support
		if err := m.RegisterTemplateWithIncludes(name, string(content)); err != nil {
			return fmt.Errorf("failed to register template %s: %w", name, err)
		}
	}

	return nil
}

// IncludeTemplate processes the template inclusion directive {{template "name" .}}
func (m *Manager) IncludeTemplate(content string) (string, error) {
	// Regex to find {{template "name" .}} directives
	templatePattern := regexp.MustCompile(`\{\{\s*template\s+"([^"]+)"\s+\.\s*\}\}`)
	matches := templatePattern.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return content, nil // No includes found
	}

	// Create a map to track already processed includes
	processedIncludes := make(map[string]bool)

	// Process all includes
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		includeName := match[1]
		if processedIncludes[includeName] {
			continue // Skip if already processed
		}

		// Check if the included template exists
		_, err := m.GetTemplate(includeName)
		if err != nil {
			return "", fmt.Errorf("template include not found: %s", includeName)
		}

		// Mark as processed
		processedIncludes[includeName] = true
	}

	return content, nil
}

// RegisterTemplateWithIncludes registers a template and processes any includes
func (m *Manager) RegisterTemplateWithIncludes(name, content string) error {
	// First register the template normally
	err := m.RegisterTemplate(name, content)
	if err != nil {
		return err
	}

	// Process includes
	m.mutex.RLock()
	tmpl, exists := m.templates[name]
	m.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("template not found after registration: %s", name)
	}

	// Find all template includes in the content
	includePattern := regexp.MustCompile(`\{\{\s*template\s+"([^"]+)"\s+\.\s*\}\}`)
	matches := includePattern.FindAllStringSubmatch(content, -1)

	// Debug logging removed

	// Process all includes found
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		includeName := match[1]
		// Debug logging removed

		includedTmpl, err := m.GetTemplate(includeName)
		if err != nil {
			// Debug logging removed
			continue // Skip if not found
		}

		// Add the included template's parse tree to our template
		m.mutex.Lock()
		for _, t := range includedTmpl.Templates() {
			if t.Name() != name { // Avoid self-reference
				// Debug logging removed
				tmpl, err = tmpl.AddParseTree(t.Name(), t.Tree)
				if err != nil {
					m.mutex.Unlock()
					return fmt.Errorf("failed to add parse tree for included template %s: %w", t.Name(), err)
				}
			}
		}

		// Also add the includedTmpl itself if it has a different name
		if includeName != name {
			tmpl, err = tmpl.AddParseTree(includeName, includedTmpl.Tree)
			if err != nil {
				m.mutex.Unlock()
				return fmt.Errorf("failed to add parse tree for the included template %s: %w", includeName, err)
			}
		}

		m.templates[name] = tmpl
		m.mutex.Unlock()
	}

	// Debug logging removed
	return nil
}

// LoadLibraryTemplates loads all templates in the default template library
func (m *Manager) LoadLibraryTemplates() error {
	// Define paths to look for library templates
	paths := []string{
		"templates/library",
		"../templates/library",
		"../../templates/library",
	}

	// Try each path
	for _, path := range paths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return m.LoadTemplatesFromDir(path)
		}
	}

	// If no library directory found, we'll skip but not error
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
