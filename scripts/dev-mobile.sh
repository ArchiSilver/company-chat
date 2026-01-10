#!/bin/bash

# Script to run both backend and mobile app for development

echo "🚀 Starting Company Chat development environment..."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

# Function to cleanup background processes
cleanup() {
    echo "🛑 Shutting down development servers..."
    kill $(jobs -p) 2>/dev/null
    exit
}

# Set trap to cleanup on script exit
trap cleanup EXIT INT TERM

echo "📱 Starting mobile development server..."
cd mobile
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
nvm use --lts
npm start &
MOBILE_PID=$!

cd ..

echo "🖥️  Starting backend server..."
go run ./cmd/app &
BACKEND_PID=$!

echo "✅ Both servers are starting up..."
echo "📱 Mobile: http://localhost:8081 (Metro bundler)"
echo "🖥️  Backend: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop both servers"

# Wait for background processes
wait