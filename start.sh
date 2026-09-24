#!/usr/bin/env bash
# 一键启动后端(8080)与前端(5173) —— Git Bash / macOS / Linux
set -e
ROOT="$(cd "$(dirname "$0")" && pwd)"
export PATH="/c/mine/code/tool/Google/go/go1.22.1/sdk/bin:$PATH"

( cd "$ROOT/backend" && go run ./cmd/server ) &
( cd "$ROOT/frontend" && [ -d node_modules ] || npm install; npm run dev ) &

wait
