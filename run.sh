#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

# PIDs of child processes
BUS_PID=""
PY_PID=""

# Cleanup function for systemd STOP signals
cleanup() {
    echo "[bus] shutting down..."

    if [[ -n "$BUS_PID" ]]; then
        echo "[bus] stopping Go process ($BUS_PID)"
        kill -TERM "$BUS_PID" 2>/dev/null || true
    fi

    if [[ -n "$PY_PID" ]]; then
        echo "[bus] stopping e-paper process ($PY_PID)"
        kill -TERM "$PY_PID" 2>/dev/null || true
    fi

    # Give them a second to exit gracefully
    sleep 1

    # Force kill if still running
    [[ -n "$BUS_PID" ]] && kill -KILL "$BUS_PID" 2>/dev/null || true
    [[ -n "$PY_PID" ]] && kill -KILL "$PY_PID" 2>/dev/null || true

    wait
    echo "[bus] shutdown complete"
}

# Handle signals from systemd
trap cleanup TERM INT EXIT

echo "[bus] building Go binary..."
# Use absolute paths if needed
go build -o bus ./cmd/main.go

echo "[bus] starting Go process..."
./bus &
BUS_PID=$!

echo "[bus] starting e-paper..."
# unbuffered output for live logs
python3 -u e-paper/main.py &
PY_PID=$!

# Wait for either process to exit
wait

# If one process dies, shut down the other
cleanup

# Exit with non-zero if a process crashed
exit 1
