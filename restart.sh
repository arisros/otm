#!/usr/bin/env bash

set -euo pipefail

PORT=8080
ENV_FILE=".env"
LOG_FILE="otm.log"
BINARY="./otm-server"
SOURCE_DIR="./cmd/otm"

echo "🔄 Restarting server..."

# Step 1: Build the binary
echo "🛠️  Building Go binary..."
if ! go build -o "$BINARY" "$SOURCE_DIR"; then
    echo "❌ Build failed."
    exit 1
fi
echo "✅ Build successful."

# Step 2: Kill any process currently using the port
PID=$(lsof -ti :$PORT || true)
if [[ -n "$PID" ]]; then
    echo "⚠️  Killing process on port $PORT (PID: $PID)..."
    kill -9 "$PID"
else
    echo "✅ No process running on port $PORT."
fi

# Step 3: Load environment variables
if [[ -f "$ENV_FILE" ]]; then
    echo "📦 Loading environment from $ENV_FILE..."
    # shellcheck disable=SC2046
    export $(grep -v '^#' "$ENV_FILE" | xargs)
else
    echo "❌ Environment file not found: $ENV_FILE"
    exit 1
fi

# Step 4: Start the server in the background
echo "🚀 Starting server..."
nohup "$BINARY" > "$LOG_FILE" 2>&1 &

echo "✅ Server started on port $PORT."
echo "📄 Logs: $LOG_FILE"
