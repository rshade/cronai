---
name: System Health Check
description: Analyzes system health metrics and provides recommendations
author: CronAI Team
version: 1.0
category: system
tags: health, monitoring, metrics, analysis
variables:
  - name: cpu_usage
    description: Current CPU usage percentage
  - name: memory_usage
    description: Current memory usage percentage
  - name: disk_space
    description: Current disk space usage percentage
  - name: network_throughput
    description: Current network throughput in MB/s
  - name: active_users
    description: Number of currently active users
---

# System Health Check

Analyze the following system metrics and provide recommendations:

1. CPU Usage: {{cpu_usage}}%
2. Memory Usage: {{memory_usage}}%
3. Disk Space: {{disk_space}}%
4. Network Throughput: {{network_throughput}} MB/s
5. Active Users: {{active_users}}

Please identify any anomalies or concerning patterns in the data. Provide specific recommendations if any metrics are approaching critical thresholds. Compare these values to our baseline expectations and highlight any deviations that require attention.
