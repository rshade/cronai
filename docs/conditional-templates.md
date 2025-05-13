# Conditional Logic in Templates

CronAI now supports conditional logic in prompt templates, allowing for dynamic content based on variables and conditions.

## Basic Syntax

### If-Else Statements

```
{{if <condition>}}
   Content when condition is true
{{else}}
   Content when condition is false
{{end}}
```

### If-ElseIf-Else Statements

```
{{if <condition1>}}
   Content when condition1 is true
{{else if <condition2>}}
   Content when condition2 is true
{{else}}
   Content when no conditions are true
{{end}}
```

## Variable Functions

### Check if Variable Exists

```
{{if hasVar .Variables "variableName"}}
   Variable exists
{{else}}
   Variable does not exist
{{end}}
```

### Get Variable with Default Value

```
{{getVar .Variables "variableName" "defaultValue"}}
```

## Comparison Operators

### String Comparisons

- Equal: `{{if eq .Variables.name "value"}}`
- Not Equal: `{{if ne .Variables.name "value"}}`
- Contains: `{{if contains .Variables.text "substring"}}`
- Has Prefix: `{{if hasPrefix .Variables.text "start"}}`
- Has Suffix: `{{if hasSuffix .Variables.text "end"}}`

### Numeric Comparisons

These operators attempt to convert strings to numbers before comparison:

- Less Than: `{{if lt .Variables.count "10"}}`
- Less Than or Equal: `{{if le .Variables.count "10"}}`
- Greater Than: `{{if gt .Variables.count "10"}}`
- Greater Than or Equal: `{{if ge .Variables.count "10"}}`

## Logical Operators

### AND Operator

```
{{if and (condition1) (condition2)}}
   Both conditions are true
{{end}}
```

### OR Operator

```
{{if or (condition1) (condition2)}}
   At least one condition is true
{{end}}
```

### NOT Operator

```
{{if not (condition)}}
   Condition is false
{{end}}
```

## Nested Conditionals

You can nest conditional blocks for more complex logic:

```
{{if condition1}}
  {{if condition2}}
    Both condition1 and condition2 are true
  {{else}}
    Only condition1 is true
  {{end}}
{{else}}
  condition1 is false
{{end}}
```

## Examples

### Environment-specific Content

```markdown
# System Analysis

{{if eq .Variables.environment "production"}}
## Production Environment
This is a production system. Be cautious with recommendations.
{{else if eq .Variables.environment "staging"}}
## Staging Environment
This is a staging system. You can suggest changes.
{{else}}
## Development Environment
This is a development system. Feel free to experiment.
{{end}}
```

### Optional Sections

```markdown
# Report Template

## Core Analysis
Always include this section...

{{if hasVar .Variables "includePerformance"}}
## Performance Analysis
This section only appears if includePerformance variable exists.
{{end}}

{{if hasVar .Variables "includeSecurityScan"}}
## Security Analysis
This section only appears if includeSecurityScan variable exists.
{{end}}
```

### Complex Conditions

```markdown
{{if and (hasVar .Variables "systemType") (eq .Variables.systemType "critical")}}
  {{if gt .Variables.errorCount "5"}}
    ## CRITICAL ALERT: Multiple Errors in Critical System
    Immediate attention required!
  {{else}}
    ## Critical System Status
    Monitoring required.
  {{end}}
{{else}}
  ## Standard System Report
  Regular monitoring protocol.
{{end}}
```

## Configuration Example

To use conditional logic in your cronai configuration:

```
# Run a prompt with variables that affect the conditional logic
0 9 * * * claude conditional_example email-team@company.com environment=production,systemType=database,errorCount=3,userQuestion=How can we optimize database performance?
```

See the `cron_prompts/conditional_example.md` file for a complete example of using conditionals in prompts.