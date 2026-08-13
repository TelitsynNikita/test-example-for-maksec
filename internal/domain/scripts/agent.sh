#!/bin/bash
# Agent for monitoring script execution
# Installed on remote Linux host via SSH

set -e

AGENT_DIR="/opt/script-monitor/agent"
SCRIPTS_DIR="/opt/script-monitor/scripts"
CALLBACK_URL="${CALLBACK_URL:-http://localhost:8081/callback}"
CALLBACK_TOKEN="${CALLBACK_TOKEN:-secret}"
LOG_FILE="${AGENT_DIR}/agent.log"

mkdir -p "$AGENT_DIR" 2>/dev/null
mkdir -p "$SCRIPTS_DIR" 2>/dev/null

log() {
    echo "[$(date -Iseconds)] $1" >> "$LOG_FILE"
}

log "Agent starting..."

# Install auditd if not present
if ! command -v auditctl &> /dev/null; then
    log "Installing auditd..."
    apt-get update -qq && apt-get install -y -qq auditd audispd-plugins 2>/dev/null || true
fi

# Ensure auditd is running
if command -v auditctl &> /dev/null; then
    log "Configuring audit rules..."
    auditctl -w "$SCRIPTS_DIR" -p wa -k script_monitor 2>/dev/null || true
    auditctl -e 1 2>/dev/null || true
fi

# Function to send callback
send_callback() {
    local user="$1"
    local script="$2"
    local action="$3"
    local timestamp="$4"

    if [ -z "$user" ] || [ -z "$script" ] || [z -z "$action" ]; then
        return
    fi

    log "Sending callback: user=$user, script=$script, action=$action"

    curl -s -X POST "$CALLBACK_URL" \
        -H "Authorization: Bearer $CALLBACK_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"user\":\"$user\",\"script\":\"$script\",\"action\":\"$action\",\"time\":\"$timestamp\"}" \
        2>/dev/null || true
}

# Monitor events using ausearch
log "Starting event monitoring..."

while true; do
    if command -v ausearch &> /dev/null; then
        # Query recent events
        ausearch -k script_monitor -f "$SCRIPTS_DIR" --raw -ts recent 2>/dev/null | while read -r line; do
            # Parse event type
            if echo "$line" | grep -q "type=EXECVE"; then
                # Execute event
                USER=$(whoami)
                SCRIPT=$(echo "$line" | grep -o "a0=[^ ]*" | head -1 | cut -d= -f2 | tr -d '"')
                ACTION="execute"
                TIMESTAMP=$(date -Iseconds)
                send_callback "$USER" "$SCRIPT" "$ACTION" "$TIMESTAMP"
            elif echo "$line" | grep -q "type=OPEN"; then
                # Open event
                USER=$(whoami)
                SCRIPT=$(echo "$line" | grep -o "name=[^ ]*" | head -1 | cut -d= -f2 | tr -d '"')
                ACTION="open"
                TIMESTAMP=$(date -Iseconds)
                send_callback "$USER" "$SCRIPT" "$ACTION" "$TIMESTAMP"
            fi
        done
    fi
    sleep 2
done