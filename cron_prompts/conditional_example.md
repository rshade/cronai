# Conditional Prompt Example

{{if hasVar .Variables "systemType"}}

## System: {{.Variables.systemType}}

{{else}}

## System: Unknown

{{end}}

{{if eq (getVar .Variables "environment" "") "production"}}
You are analyzing a PRODUCTION environment. Please be extra careful with your recommendations.
{{else if eq (getVar .Variables "environment" "") "staging"}}
You are analyzing a STAGING environment. You can suggest more experimental approaches.
{{else}}
You are analyzing a DEVELOPMENT environment. Feel free to suggest any improvements.
{{end}}

{{if hasVar .Variables "includeMetrics"}}
Please include detailed metrics in your analysis.
{{end}}

{{if and (hasVar .Variables "criticalSystem") (eq .Variables.criticalSystem "true")}}
THIS IS A CRITICAL SYSTEM. Any downtime has severe business impact.
{{end}}

{{if or (eq (getVar .Variables "priority" "low") "high") (eq (getVar .Variables "severity" "low") "high")}}
This issue has been marked as HIGH PRIORITY.
{{end}}

{{if gt (getVar .Variables "errorCount" "0") "5"}}
Multiple errors detected ({{.Variables.errorCount}}). This suggests a systemic issue.
{{else if gt (getVar .Variables "errorCount" "0") "0"}}
A few errors detected ({{.Variables.errorCount}}). This could be an isolated issue.
{{else}}
No errors detected.
{{end}}

Your task is to analyze the following information and provide recommendations:

{{.Variables.userQuestion}}

Please structure your response with:

1. Analysis of current state
2. Identified issues
3. Recommendations
4. Next steps
