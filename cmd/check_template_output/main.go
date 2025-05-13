package main

import (
	"fmt"

	"github.com/rshade/cronai/internal/processor/template"
)

func main() {
	// Create the template
	templateContent := `# Report for {{.Variables.project}}

{{if eq .Variables.environment "production"}}
## Production Environment Status
Current status: {{.Variables.status}}
{{if eq .Variables.status "healthy"}}
All systems operational.
{{else}}
Warning: System requires attention!
{{end}}
{{else}}
## Test Environment Status
This is a test environment.
{{end}}

Report generated on {{.Variables.date}}.`

	// Create the manager
	manager := template.GetManager()
	templateName := "example_output"
	err := manager.RegisterTemplate(templateName, templateContent)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// Create data
	data := template.TemplateData{
		Variables: map[string]string{
			"project":     "CronAI",
			"environment": "production",
			"status":      "healthy",
			"date":        "2025-05-12",
		},
	}

	// Execute
	result, err := manager.Execute(templateName, data)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// Print result with byte representation to see exact whitespace
	fmt.Println("Result as string:")
	fmt.Println(result)
	fmt.Println("\nResult as bytes:")
	for i, b := range []byte(result) {
		fmt.Printf("%d ", b)
		if i > 0 && i%20 == 0 {
			fmt.Println()
		}
	}
}
