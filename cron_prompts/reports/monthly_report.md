---
name: Monthly Report
description: Generates a comprehensive monthly report
author: CronAI Team
version: 1.0
category: report
tags: monthly, report, summary
variables:
  - name: month
    description: The month for the report (e.g., "January")
  - name: year
    description: The year for the report (e.g., "2025")
  - name: team
    description: The team this report is for
---

{{include "templates/common_header.md"}}

# Monthly Report: {{month}} {{year}}

## {{team}} Team Performance Summary

Please generate a comprehensive monthly report for the {{team}} team covering the month of {{month}} {{year}}. Include the following sections:

1. Key Accomplishments
2. Challenges Faced
3. Performance Metrics
4. Goals for Next Month
5. Resource Requirements

Each section should be detailed and include specific recommendations where appropriate. The report should be professional in tone and suitable for presentation to executive leadership.