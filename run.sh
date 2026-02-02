#!/usr/bin/env bash
set -e

echo "[bus] starting service..."

cd "$(dirname "$0")"

echo "[bus] building Go binary..."
go build -o bus ./cmd/main.go

echo "[bus] starting Go process..."
./bus &

echo "[bus] starting e-paper..."
exec python3 e-paper/main.py
