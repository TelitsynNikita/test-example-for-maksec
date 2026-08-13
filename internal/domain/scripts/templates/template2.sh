#!/bin/bash
# Template 2: Disk Space Monitor
# Generated at: $(date)

LOG_FILE="/var/log/script-monitor/disk-monitor.log"
mkdir -p "$(dirname $LOG_FILE)"

echo "[$(date)] Disk space check by user: $(whoami)" >> "$LOG_FILE"

df -h | grep -v "tmpfs" | while read line; do
    USAGE=$(echo $line | awk '{print $5}' | sed 's/%//')
    MOUNT=$(echo $line | awk '{print $6}')
    if [ $USAGE -gt 80 ]; then
        echo "[$(date)] WARNING: $MOUNT is $USAGE% full" >> "$LOG_FILE"
    fi
done

exit 0