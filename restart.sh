
#!/bin/bash

set -euo pipefail

PORT=8080
ENV_FILE="config/.env"
LOG_FILE="otm.log"
BINARY="./otm-server"

echo "🔄 Restarting server..."

# Step 1: Kill any process currently using the port
PID=$(lsof -ti :$PORT || true)
if [[ -n "$PID" ]]; then
    echo "⚠️  Killing process on port $PORT (PID: $PID)..."
    kill -9 "$PID"
else
    echo "✅ No process running on port $PORT."
fi

# Step 2: Load environment variables
if [[ -f "$ENV_FILE" ]]; then
    echo "📦 Loading environment from $ENV_FILE..."
    # shellcheck disable=SC2046
    export $(grep -v '^#' "$ENV_FILE" | xargs)
else
    echo "❌ Environment file not found: $ENV_FILE"
    exit 1
fi

# Step 3: Check if binary exists
if [[ ! -x "$BINARY" ]]; then
    echo "❌ Binary not found or not executable: $BINARY"
    exit 1
fi

# Step 4: Start server in background and redirect output to log
echo "🚀 Starting server..."
nohup "$BINARY" > "$LOG_FILE" 2>&1 &

echo "✅ Server started in background on port $PORT."
echo "📄 Logs: $LOG_FILE"
