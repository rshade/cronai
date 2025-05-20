# System Health Check with Conditional Logic

{{if eq (getVar .Variables "role" "analyst") "admin"}}

## Admin System Health Check

You are conducting a detailed system health check with admin privileges. This check should be thorough and include all system metrics.
{{else if eq (getVar .Variables "role" "analyst") "operator"}}

## Operator System Health Check

You are conducting a system health check with operator privileges. Focus on operational metrics and indicators.
{{else}}

## General System Health Check

You are conducting a basic system health check. Focus on high-level indicators and summary information.
{{end}}

{{if hasVar .Variables "timeframe"}}
Analyze system health for the {{.Variables.timeframe}} timeframe.
{{else}}
Analyze system health for the past 24 hours.
{{end}}

{{if hasVar .Variables "environment"}}
Target environment: {{.Variables.environment}}
{{end}}

{{if eq (getVar .Variables "includePerformance" "true") "true"}}

## Performance Analysis

Include a detailed analysis of system performance metrics:

- CPU utilization
- Memory usage
- Disk I/O
- Network throughput
- Response times
{{end}}

{{if eq (getVar .Variables "includeCapacity" "false") "true"}}

## Capacity Planning

Include capacity planning recommendations based on current usage trends.
{{end}}

{{if gt (getVar .Variables "criticalErrors" "0") "0"}}

## Critical Issues

{{if gt .Variables.criticalErrors "5"}}
URGENT: Multiple critical errors detected ({{.Variables.criticalErrors}}). Immediate investigation required.
{{else}}
ATTENTION: {{.Variables.criticalErrors}} critical errors detected. Investigate as soon as possible.
{{end}}
{{else}}
No critical errors detected in this period.
{{end}}

{{if hasVar .Variables "customMetrics"}}

## Custom Metrics

Include analysis of the following custom metrics: {{.Variables.customMetrics}}
{{end}}

## Required Analysis

1. Evaluate the current health status of all monitored systems
2. Identify any anomalies or deviations from normal operation
3. Provide actionable recommendations for any issues discovered
4. Prioritize recommendations based on impact and urgency

{{if hasVar .Variables "additionalInstructions"}}

## Additional Instructions

{{.Variables.additionalInstructions}}
{{end}}

Report generation time: {{.Timestamp.Format "Jan 02, 2006 15:04:05"}}
