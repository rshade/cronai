package template

import "fmt"

// registerOrPanic registers a template and panics if there is an error
func registerOrPanic(m *Manager, name, content string) {
	if err := m.RegisterTemplate(name, content); err != nil {
		panic(fmt.Sprintf("Failed to register default template %s: %v", name, err))
	}
}

// registerDefaultTemplates adds built-in templates
func (m *Manager) registerDefaultTemplates() {
	// Email templates
	registerOrPanic(m, "default_email_subject", "Response from {{.Model}} - {{.PromptName}}")
	registerOrPanic(m, "default_email_html", `
<html>
<body>
<h1>AI Response: {{.PromptName}}</h1>
<p><strong>Model:</strong> {{.Model}}</p>
<p><strong>Time:</strong> {{.Timestamp.Format "Jan 02, 2006 15:04:05"}}</p>
<div>
{{.Content}}
</div>
</body>
</html>
`)
	registerOrPanic(m, "default_email_text", `
AI Response: {{.PromptName}}
Model: {{.Model}}
Time: {{.Timestamp.Format "Jan 02, 2006 15:04:05"}}

{{.Content}}
`)

	// Slack templates
	slackTemplate := `{"blocks":[{"type":"header","text":{"type":"plain_text","text":"AI Response: {{.PromptName}}"}},{"type":"section","fields":[{"type":"mrkdwn","text":"*Model:* {{.Model}}"},{"type":"mrkdwn","text":"*Time:* {{.Timestamp}}"}]},{"type":"section","text":{"type":"mrkdwn","text":"{{.Content}}"}}]}`
	registerOrPanic(m, "default_slack", slackTemplate)

	// Webhook templates
	webhookTemplate := `{"timestamp":"{{.Timestamp}}","model":"{{.Model}}","prompt":"{{.PromptName}}","content":"{{.Content}}"}`
	registerOrPanic(m, "default_webhook", webhookTemplate)

	// File templates
	fileContentTemplate := `# AI Response: {{.PromptName}}
Model: {{.Model}}
Time: {{.Timestamp}}

{{.Content}}
`
	registerOrPanic(m, "default_file_content", fileContentTemplate)

	filenameTemplate := `logs/{{.Model}}-{{.Timestamp}}.txt`
	registerOrPanic(m, "default_file_filename", filenameTemplate)
}
