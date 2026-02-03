#!/usr/bin/env bash
set -e

echo "[bus] starting service..."

cd "$(dirname "$0")"

cleanup() {
  echo "[bus] shutting down..."

  if [[ -n "$BUS_PID" ]]; then
    echo "[bus] stopping Go process ($BUS_PID)"
    kill "$BUS_PID" 2>/dev/null
  fi

  if [[ -n "$PY_PID" ]]; then
    echo "[bus] stopping e-paper ($PY_PID)"
    kill "$PY_PID" 2>/dev/null
  fi

  wait
  echo "[bus] shutdown complete"
}

# Run cleanup on Ctrl-C, SIGTERM, and script exit
trap cleanup INT TERM EXIT

echo "[bus] building Go binary..."
go build -o bus ./cmd/main.go

echo "[bus] starting Go process..."
./bus &
BUS_PID=$!

echo "[bus] starting e-paper..."
python3 e-paper/main.py &
PY_PID=$!

wait
