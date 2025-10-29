#!/bin/sh
set -e

# Function to handle shutdown signals
shutdown() {
    echo "Received shutdown signal, stopping services..."
    
    # Stop MCP server (if running)
    if [ -n "$MCP_PID" ]; then
        kill -TERM "$MCP_PID" 2>/dev/null || true
        wait "$MCP_PID" 2>/dev/null || true
    fi
    
    # Stop Valkey server
    if command -v valkey-cli >/dev/null 2>&1; then
        valkey-cli shutdown 2>/dev/null || true
    fi
    
    echo "Services stopped gracefully"
    exit 0
}

# Register signal handlers (using signal numbers for sh compatibility)
# SIGTERM=15, SIGINT=2
trap shutdown 15 2

echo "Starting Valkey server..."
valkey-server --daemonize yes --bind 127.0.0.1 --port 6379

echo "Waiting for Valkey to be ready..."
MAX_RETRIES=30
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if valkey-cli ping >/dev/null 2>&1; then
        echo "Valkey is ready!"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "Waiting for Valkey... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 1
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "ERROR: Valkey failed to start within timeout"
    exit 1
fi

echo "Starting MCP server..."
# Start MCP server in background to allow signal handling
/usr/local/bin/mcp-ruleset-server &
MCP_PID=$!

echo "MCP Ruleset Server is running (PID: $MCP_PID)"

# Wait for MCP server process
wait $MCP_PID
