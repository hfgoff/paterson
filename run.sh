#!/usr/bin/env bash
set -e

echo "[bus] starting service..."

cd "$(dirname "$0")"

cleanup() {
  echo "[bus] shutting down..."
  kill 0
  wait
}

echo "[bus] building Go binary..."
go build -o bus ./cmd/main.go

echo "[bus] starting Go process..."
./bus &
BUS_PID=$!

echo "[bus] starting e-paper..."
python3 e-paper/main.py &
PY_PID=$!

wait
