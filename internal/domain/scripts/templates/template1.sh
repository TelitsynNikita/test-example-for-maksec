#!/bin/bash
# Template 1: Basic Monitoring Script
# Generated at: $(date)

LOG_FILE="/var/log/script-monitor/$(basename $0).log"
mkdir -p "$(dirname $LOG_FILE)"

echo "[$(date)] Script executed by user: $(whoami)" >> "$LOG_FILE"
echo "[$(date)] Script path: $0" >> "$LOG_FILE"
echo "[$(date)] Arguments: $@" >> "$LOG_FILE"

echo "Monitoring started..."

CPU_LOAD=$(uptime | awk -F'load average:' '{print $2}')
echo "[$(date)] CPU Load: $CPU_LOAD" >> "$LOG_FILE"

exit 0