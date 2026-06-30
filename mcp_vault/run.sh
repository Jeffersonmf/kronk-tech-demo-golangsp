#!/usr/bin/env bash
# Builds and runs both MCP servers (vault_dumb on :9001, vault_smart on :9002)
# side by side. Ctrl+C stops both.
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

# -- cleanup: kill anything that might be leftover from a previous run --
echo "==> limpando processos anteriores..."
pkill -f "bin/dumb"  2>/dev/null || true
pkill -f "bin/smart" 2>/dev/null || true
pkill kronk          2>/dev/null || true
# wait a moment for ports to be released
sleep 1

GO="${GO:-$(command -v go)}"
unset GOTOOLCHAIN
export KRONK_PROCESSOR="${KRONK_PROCESSOR:-metal}"

echo "==> compilando servidores MCP..."
mkdir -p bin
"$GO" build -o bin/dumb ./cmd/dumb
"$GO" build -o bin/smart ./cmd/smart

echo "==> iniciando vault_dumb em :9001..."
./bin/dumb &
dumb_pid=$!

echo "==> iniciando vault_smart em :9002..."
./bin/smart &
smart_pid=$!

trap 'echo "==> encerrando..."; kill "$dumb_pid" "$smart_pid" 2>/dev/null' EXIT INT TERM

echo "==> pronto! Ctrl+C para parar."
wait "$dumb_pid" "$smart_pid"
