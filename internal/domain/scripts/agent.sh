#!/bin/bash
# Agent for monitoring script execution
# Installed on remote Linux host via SSH

set -e

AGENT_DIR="/opt/script-monitor/agent"
SCRIPTS_DIR="/opt/script-monitor/scripts"

CALLBACK_URL="${AGENT_CALLBACK_URL:-http://localhost:8081/callback}"
CALLBACK_TOKEN="${AGENT_CALLBACK_TOKEN:-secret}"

LOG_FILE="${AGENT_DIR}/agent.log"
LAST_EVENT_FILE="${AGENT_DIR}/last_event"

mkdir -p "$AGENT_DIR" 2>/dev/null
mkdir -p "$SCRIPTS_DIR" 2>/dev/null

log() {
    echo "[$(date -Iseconds)] $1" >> "$LOG_FILE"
}

# ===== ОБРАБОТКА АРГУМЕНТОВ =====
if [ "$1" = "setup-systemd" ]; then
    setup_systemd
    exit $?
fi

get_last_event_time() {
    if [ -f "$LAST_EVENT_FILE" ]; then
        cat "$LAST_EVENT_FILE"
    else
        echo "1970-01-01T00:00:00+00:00"
    fi
}

update_last_event_time() {
    echo "$1" > "$LAST_EVENT_FILE"
}

get_username_by_uid() {
    local uid="$1"
    if [ -z "$uid" ]; then
        echo ""
        return
    fi

    local username=""
    if command -v getent &> /dev/null; then
        username=$(getent passwd "$uid" 2>/dev/null | cut -d: -f1)
    fi

    if [ -z "$username" ]; then
        username=$(cat /etc/passwd 2>/dev/null | grep ":${uid}:" | cut -d: -f1)
    fi

    if [ -z "$username" ]; then
        username="uid_${uid}"
    fi

    echo "$username"
}

extract_uid_from_event() {
    local line="$1"
    local uid=""

    uid=$(echo "$line" | grep -o "uid=[0-9]*" | head -1 | cut -d= -f2)
    if [ -n "$uid" ]; then
        echo "$uid"
        return
    fi

    uid=$(echo "$line" | grep -o "auid=[0-9]*" | head -1 | cut -d= -f2)
    if [ -n "$uid" ]; then
        echo "$uid"
        return
    fi

    echo "0"
}

extract_script_path_from_execve() {
    local line="$1"
    local script_path=""

    for i in 0 1 2 3 4 5 6 7 8 9; do
        local arg=$(echo "$line" | grep -o "a${i}=[^ ]*" | head -1 | cut -d= -f2 | tr -d '"' | tr -d "'")
        if [ -z "$arg" ]; then
            continue
        fi

        if echo "$arg" | grep -q "^/opt/script-monitor/scripts/"; then
            script_path="$arg"
            break
        fi

        if echo "$arg" | grep -q "\.sh" && echo "$arg" | grep -q "/opt/script-monitor/scripts/"; then
            script_path="$arg"
            break
        fi
    done

    if [ -z "$script_path" ]; then
        for i in 0 1 2 3 4 5 6 7 8 9; do
            local arg=$(echo "$line" | grep -o "a${i}=[^ ]*" | head -1 | cut -d= -f2 | tr -d '"' | tr -d "'")
            if [ -z "$arg" ]; then
                continue
            fi

            if echo "$arg" | grep -q "/opt/script-monitor/scripts/"; then
                script_path="$arg"
                break
            fi
        done
    fi

    if [ -z "$script_path" ]; then
        for i in 0 1 2 3 4 5 6 7 8 9; do
            local arg=$(echo "$line" | grep -o "a${i}=[^ ]*" | head -1 | cut -d= -f2 | tr -d '"' | tr -d "'")
            if [ -z "$arg" ]; then
                continue
            fi

            if echo "$arg" | grep -q "\.sh$" && echo "$arg" | grep -q "/opt/script-monitor/"; then
                script_path="$arg"
                break
            fi
        done
    fi

    echo "$script_path"
}

extract_script_path_from_open() {
    local line="$1"
    local script_path=""

    local names=$(echo "$line" | grep -o "name=[^ ]*" | cut -d= -f2 | tr -d '"' | tr -d "'")
    for name in $names; do
        if echo "$name" | grep -q "^/opt/script-monitor/scripts/"; then
            script_path="$name"
            break
        fi
    done

    echo "$script_path"
}

detect_package_manager() {
    if command -v apt-get &> /dev/null; then
        echo "apt-get"
    elif command -v dnf &> /dev/null; then
        echo "dnf"
    elif command -v yum &> /dev/null; then
        echo "yum"
    elif command -v apk &> /dev/null; then
        echo "apk"
    elif command -v zypper &> /dev/null; then
        echo "zypper"
    else
        echo "unknown"
    fi
}

install_auditd() {
    local pm=$(detect_package_manager)
    log "Detected package manager: $pm"

    case "$pm" in
        "apt-get")
            log "Installing auditd via apt-get..."
            apt-get update -qq 2>/dev/null || true
            apt-get install -y -qq auditd audispd-plugins 2>/dev/null || true
            ;;
        "dnf")
            log "Installing audit via dnf..."
            dnf install -y audit audit-libs 2>/dev/null || true
            ;;
        "yum")
            log "Installing audit via yum..."
            yum install -y audit audit-libs 2>/dev/null || true
            ;;
        "apk")
            log "Installing audit via apk..."
            apk add audit 2>/dev/null || true
            ;;
        "zypper")
            log "Installing audit via zypper..."
            zypper install -y audit 2>/dev/null || true
            ;;
        *)
            log "Unknown package manager: $pm. Skipping auditd installation."
            log "Please install auditd manually: auditctl, ausearch"
            ;;
    esac
}

ensure_auditd_installed() {
    if ! command -v auditctl &> /dev/null || ! command -v ausearch &> /dev/null; then
        log "auditctl or ausearch not found. Attempting to install..."
        install_auditd
    fi

    if command -v auditctl &> /dev/null && command -v ausearch &> /dev/null; then
        log "auditctl and ausearch are available"
        return 0
    else
        log "WARNING: auditctl or ausearch still not available after installation"
        log "Agent will run in limited mode (only logging)"
        return 1
    fi
}

build_json_payload() {
    local user="$1"
    local script="$2"
    local action="$3"
    local timestamp="$4"

    local escaped_user=$(echo "$user" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | tr -d '\n\r')
    local escaped_script=$(echo "$script" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | tr -d '\n\r')
    local escaped_action=$(echo "$action" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | tr -d '\n\r')
    local escaped_timestamp=$(echo "$timestamp" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | tr -d '\n\r')

    echo "{\"user\":\"$escaped_user\",\"script\":\"$escaped_script\",\"action\":\"$escaped_action\",\"time\":\"$escaped_timestamp\"}"
}

send_callback() {
    local user="$1"
    local script="$2"
    local action="$3"
    local timestamp="$4"

    if [ -z "$user" ] || [ -z "$script" ] || [ -z "$action" ]; then
        return
    fi

    log "Sending callback: user=$user, script=$script, action=$action"

    local payload=$(build_json_payload "$user" "$script" "$action" "$timestamp")

    if command -v jq &> /dev/null; then
        if echo "$payload" | jq . > /dev/null 2>&1; then
            curl -s -X POST "$CALLBACK_URL" \
                -H "Authorization: Bearer $CALLBACK_TOKEN" \
                -H "Content-Type: application/json" \
                -d "$payload" \
                2>/dev/null || true
        else
            log "ERROR: Invalid JSON payload, falling back to manual"
            local fallback_payload="{\"user\":\"$user\",\"script\":\"$script\",\"action\":\"$action\",\"time\":\"$timestamp\"}"
            curl -s -X POST "$CALLBACK_URL" \
                -H "Authorization: Bearer $CALLBACK_TOKEN" \
                -H "Content-Type: application/json" \
                -d "$fallback_payload" \
                2>/dev/null || true
        fi
    else
        curl -s -X POST "$CALLBACK_URL" \
            -H "Authorization: Bearer $CALLBACK_TOKEN" \
            -H "Content-Type: application/json" \
            -d "$payload" \
            2>/dev/null || true
    fi
}

setup_systemd() {
    local service_file="/etc/systemd/system/script-monitor-agent.service"

    log "Setting up systemd service..."

    if [ -f "$service_file" ]; then
        log "Systemd service already exists"
        return 0
    fi

    if ! command -v systemctl &> /dev/null; then
        log "systemctl not found, skipping systemd setup"
        return 1
    fi

    cat > "$service_file" << 'EOF'
[Unit]
Description=Script Monitor Agent
Documentation=https://github.com/TelitsynNikita/script-monitor
After=auditd.service network.target
Requires=auditd.service

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/opt/script-monitor/agent
ExecStart=/opt/script-monitor/agent/agent.sh
Restart=always
RestartSec=10
StandardOutput=append:/opt/script-monitor/agent/agent.log
StandardError=append:/opt/script-monitor/agent/agent.log

NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/opt/script-monitor/agent /opt/script-monitor/scripts /var/log/script-monitor

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable script-monitor-agent 2>/dev/null || true
    systemctl start script-monitor-agent 2>/dev/null || true

    if systemctl is-active --quiet script-monitor-agent; then
        log "Systemd service started successfully"
    else
        log "WARNING: Systemd service failed to start"
        systemctl status script-monitor-agent 2>/dev/null || true
        return 1
    fi

    return 0
}

log "Agent starting..."

ensure_auditd_installed

if command -v auditctl &> /dev/null; then
    log "Configuring audit rules..."
    auditctl -w "$SCRIPTS_DIR" -p wa -k script_monitor 2>/dev/null || true
    auditctl -e 1 2>/dev/null || true
fi

log "Starting event monitoring..."

while true; do
    if command -v ausearch &> /dev/null; then
        LAST_EVENT_TIME=$(get_last_event_time)

        ausearch -k script_monitor -f "$SCRIPTS_DIR" --raw -ts "$LAST_EVENT_TIME" 2>/dev/null | while read -r line; do
            EVENT_TIME=$(echo "$line" | grep -o "time=[0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}T[0-9]\{2\}:[0-9]\{2\}:[0-9]\{2\}" | head -1 | cut -d= -f2)

            if [ -n "$EVENT_TIME" ]; then
                update_last_event_time "$EVENT_TIME"
            fi

            UID=$(extract_uid_from_event "$line")
            USER=$(get_username_by_uid "$UID")

            if echo "$line" | grep -q "type=EXECVE"; then
                SCRIPT=$(extract_script_path_from_execve "$line")
                if [ -n "$SCRIPT" ]; then
                    ACTION="execute"
                    TIMESTAMP=$(date -Iseconds)
                    send_callback "$USER" "$SCRIPT" "$ACTION" "$TIMESTAMP"
                fi
            elif echo "$line" | grep -q "type=OPEN"; then
                SCRIPT=$(extract_script_path_from_open "$line")
                if [ -n "$SCRIPT" ]; then
                    ACTION="open"
                    TIMESTAMP=$(date -Iseconds)
                    send_callback "$USER" "$SCRIPT" "$ACTION" "$TIMESTAMP"
                fi
            fi
        done
    fi
    sleep 2
done