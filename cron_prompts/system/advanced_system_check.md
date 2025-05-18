---
name: Advanced System Check
description: Enhanced system health check with extended sections
author: CronAI Team
version: 1.1
category: system
tags: advanced, monitoring, health
extends: templates/base_system_check
---

{{include "templates/common_header.md"}}

{{include "templates/base_system_check.md"}}

## Additional Checks

### Security Analysis

Check for any potential security vulnerabilities or unusual access patterns.

### Service Health

Verify that all critical services are running and responding as expected.

### Backup Status

Check recent backup status and verify backup integrity.

## Performance Trends

Analyze performance trends over the past {{if hasVar .Variables "timespan"}}{{.Variables.timespan}}{{else}}week{{end}} and highlight any concerning patterns.
