#!/bin/bash
# Build script for Render deployment

set -e  # Exit on error

echo "📦 Downloading dependencies..."
go mod download

echo "✅ Verifying modules..."
go mod verify

echo "🔨 Building application..."
go build -tags netgo -ldflags '-s -w' -o main

echo "✅ Build complete!"

