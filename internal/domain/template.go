package domain

type ScriptTemplate struct {
	Name        string
	Description string
	Content     string
}

func GetTemplates() map[string]ScriptTemplate {
	return map[string]ScriptTemplate{
		"template1": {
			Name:        "template1",
			Description: "Базовый шаблон мониторинга",
			Content: `#!/bin/bash
# Template 1: Basic Monitoring Script
# Generated at: $(date)

LOG_FILE="/var/log/script-monitor/$(basename $0).log"
mkdir -p "$(dirname $LOG_FILE)"

echo "[$(date)] Script executed by user: $(whoami)" >> "$LOG_FILE"
echo "[$(date)] Script path: $0" >> "$LOG_FILE"
echo "[$(date)] Arguments: $@" >> "$LOG_FILE"

# Ваша логика здесь
echo "Monitoring started..."

# Пример: проверка CPU
CPU_LOAD=$(uptime | awk -F'load average:' '{print $2}')
echo "[$(date)] CPU Load: $CPU_LOAD" >> "$LOG_FILE"

exit 0
`,
		},
		"template2": {
			Name:        "template2",
			Description: "Шаблон для проверки дискового пространства",
			Content: `#!/bin/bash
# Template 2: Disk Space Monitor
# Generated at: $(date)

LOG_FILE="/var/log/script-monitor/disk-monitor.log"
mkdir -p "$(dirname $LOG_FILE)"

echo "[$(date)] Disk space check by user: $(whoami)" >> "$LOG_FILE"

# Проверка дискового пространства
df -h | grep -v "tmpfs" | while read line; do
    USAGE=$(echo $line | awk '{print $5}' | sed 's/%//')
    MOUNT=$(echo $line | awk '{print $6}')
    if [ $USAGE -gt 80 ]; then
        echo "[$(date)] WARNING: $MOUNT is $USAGE% full" >> "$LOG_FILE"
    fi
done

exit 0
`,
		},
	}
}
