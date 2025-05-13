# Weekly Project Status Report

{{if hasVar .Variables "project"}}
## Project: {{.Variables.project}}
{{else}}
## Project: General Status Report
{{end}}

{{if eq (getVar .Variables "audience" "general") "executive"}}
### Executive Summary
This is a high-level summary for executive stakeholders. Focus on strategic impact, ROI, and business outcomes.
{{else if eq (getVar .Variables "audience" "general") "technical"}}
### Technical Update
This is a detailed technical update for the engineering team. Include technical details, challenges, and implementation specifics.
{{else}}
### General Update
This is a general status update for all stakeholders.
{{end}}

{{if hasVar .Variables "reportingPeriod"}}
Reporting period: {{.Variables.reportingPeriod}}
{{else}}
Reporting period: Last 7 days
{{end}}

## Overview

Please provide a comprehensive weekly project status report covering the following areas:

{{if eq (getVar .Variables "includeProgress" "true") "true"}}
## Progress Updates
- Summarize key accomplishments from the past week
- List completed deliverables and milestones
- Highlight any significant achievements
{{end}}

{{if eq (getVar .Variables "includeMetrics" "false") "true"}}
## Key Metrics
{{if hasVar .Variables "specificMetrics"}}
Focus on these specific metrics: {{.Variables.specificMetrics}}
{{else}}
Include standard project metrics (velocity, burn rate, etc.)
{{end}}
{{end}}

## Current Status
{{if hasVar .Variables "currentStatus"}}
The project is currently in status: {{.Variables.currentStatus}}
{{if eq .Variables.currentStatus "at-risk"}}
Please provide a detailed explanation of risk factors and mitigation strategies.
{{else if eq .Variables.currentStatus "blocked"}}
Please clearly identify blockers and suggest approaches to resolve them.
{{else if eq .Variables.currentStatus "on-track"}}
Please confirm all milestones and deliverables are on schedule.
{{end}}
{{else}}
Assess the overall project status (on-track, at-risk, or blocked).
{{end}}

{{if gt (getVar .Variables "issueCount" "0") "0"}}
## Issues and Blockers
{{if gt .Variables.issueCount "5"}}
Multiple issues detected ({{.Variables.issueCount}}). Please prioritize and provide recommendations for the top 3 most critical issues.
{{else}}
{{.Variables.issueCount}} issues detected. Please provide recommendations for each issue.
{{end}}
{{end}}

## Next Steps
- Outline planned work for the coming week
- Identify upcoming milestones and deliverables
- Highlight any anticipated challenges

{{if eq (getVar .Variables "includeRisks" "false") "true"}}
## Risk Assessment
Identify potential risks and provide mitigation strategies.
{{end}}

{{if hasVar .Variables "additionalSections"}}
## Additional Requested Information
{{.Variables.additionalSections}}
{{end}}

## Conclusion
Provide a brief summary of the current state and outlook for the coming week.

Report generation: {{.Timestamp.Format "Jan 02, 2006"}}