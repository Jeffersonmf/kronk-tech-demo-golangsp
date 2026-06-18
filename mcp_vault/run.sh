#!/usr/bin/env bash
# Builds and runs both MCP servers (vault_dumb on :9001, vault_smart on :9002)
# side by side. Ctrl+C stops both.
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

GO="${GO:-/home/jmarchetti/.local/share/mise/installs/go/1.25.1/bin/go}"
unset GOTOOLCHAIN
export KRONK_PROCESSOR="${KRONK_PROCESSOR:-cpu}"

mkdir -p bin
"$GO" build -o bin/dumb ./cmd/dumb
"$GO" build -o bin/smart ./cmd/smart

./bin/dumb &
dumb_pid=$!

./bin/smart &
smart_pid=$!

trap 'kill "$dumb_pid" "$smart_pid" 2>/dev/null' EXIT INT TERM

wait "$dumb_pid" "$smart_pid"
