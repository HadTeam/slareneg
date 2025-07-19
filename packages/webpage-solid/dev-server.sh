#!/usr/bin/env zsh

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to kill process on a specific port
kill_port() {
    local port=$1
    echo "${YELLOW}Checking for processes on port $port...${NC}"
    
    # Get PID of process using the port
    local pid=$(lsof -ti:$port)
    
    if [[ -n "$pid" ]]; then
        echo "${RED}Found process $pid on port $port, killing it...${NC}"
        kill -9 $pid
        sleep 1
        
        # Verify the process was killed
        if lsof -ti:$port > /dev/null 2>&1; then
            echo "${RED}Failed to kill process on port $port${NC}"
            return 1
        else
            echo "${GREEN}Successfully killed process on port $port${NC}"
        fi
    else
        echo "${GREEN}No process found on port $port${NC}"
    fi
    return 0
}

# Function to wait for a port to be available
wait_for_port() {
    local port=$1
    local max_attempts=10
    local attempt=0
    
    while (( attempt < max_attempts )); do
        if ! lsof -ti:$port > /dev/null 2>&1; then
            return 0
        fi
        echo "Waiting for port $port to be free... (attempt $((attempt + 1))/$max_attempts)"
        sleep 1
        ((attempt++))
    done
    
    echo "${RED}Port $port is still in use after $max_attempts attempts${NC}"
    return 1
}

# Main script
echo "${GREEN}=== Development Server Manager ===${NC}"
echo ""

# Kill processes on ports 5173 and 8080
echo "${YELLOW}Step 1: Clearing ports...${NC}"
kill_port 5173
kill_port 8080
echo ""

# Wait for ports to be free
wait_for_port 5173 || exit 1
wait_for_port 8080 || exit 1

# Build the server
echo "${YELLOW}Step 2: Building server...${NC}"
if [[ -f "package.json" ]]; then
    # Check if there's a build script for the server
    if grep -q '"build:server"' package.json; then
        pnpm run build:server
    elif grep -q '"build"' package.json; then
        pnpm run build
    else
        echo "${YELLOW}No build script found, skipping server build${NC}"
    fi
else
    echo "${RED}No package.json found in current directory${NC}"
    exit 1
fi
echo ""

# Start Vite dev server on port 5173
echo "${YELLOW}Step 3: Starting Vite dev server on port 5173...${NC}"
if grep -q '"dev"' package.json; then
    # Run Vite in background and save PID
    pnpm run dev --port 5173 --host &
    VITE_PID=$!
    echo "${GREEN}Vite started with PID: $VITE_PID${NC}"
else
    echo "${RED}No dev script found in package.json${NC}"
    exit 1
fi

# Give Vite time to start
sleep 3

# Start server on port 8080
echo "${YELLOW}Step 4: Starting server on port 8080...${NC}"
if grep -q '"start:server"' package.json; then
    # Run server in background and save PID
    PORT=8080 pnpm run start:server &
    SERVER_PID=$!
    echo "${GREEN}Server started with PID: $SERVER_PID${NC}"
elif grep -q '"start"' package.json; then
    PORT=8080 pnpm run start &
    SERVER_PID=$!
    echo "${GREEN}Server started with PID: $SERVER_PID${NC}"
else
    echo "${YELLOW}No server start script found${NC}"
fi

echo ""
echo "${GREEN}=== Servers Running ===${NC}"
echo "Vite dev server: http://localhost:5173"
echo "Backend server: http://localhost:8080"
echo ""
echo "${YELLOW}Press Ctrl+C to stop all servers${NC}"

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "${YELLOW}Shutting down servers...${NC}"
    
    if [[ -n "$VITE_PID" ]]; then
        kill $VITE_PID 2>/dev/null
        echo "${GREEN}Vite server stopped${NC}"
    fi
    
    if [[ -n "$SERVER_PID" ]]; then
        kill $SERVER_PID 2>/dev/null
        echo "${GREEN}Backend server stopped${NC}"
    fi
    
    # Double-check and clean up any remaining processes
    kill_port 5173 2>/dev/null
    kill_port 8080 2>/dev/null
    
    echo "${GREEN}Cleanup complete${NC}"
    exit 0
}

# Set up trap to handle Ctrl+C
trap cleanup INT TERM

# Wait for both processes
wait
