#!/bin/bash

# Discipleship Journal Docker Launch Script
# This script helps launch the application using Docker Compose

set -e

echo "========================================="
echo "Discipleship Journal - Docker Launch"
echo "========================================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  No .env file found. Creating from template..."
    cp .env.example .env
    echo "✅ Created .env file."
    echo ""
    echo "Please edit the .env file with your configuration:"
    echo "1. Update Firebase credentials (get from Firebase Console)"
    echo "2. Optionally add Bible API credentials"
    echo "3. Review other settings as needed"
    echo ""
    echo "For Firebase setup:"
    echo "1. Go to https://console.firebase.google.com/"
    echo "2. Create a new project or use existing"
    echo "3. Enable Authentication (Email/Password or Google)"
    echo "4. Go to Project Settings > General"
    echo "5. Copy the Firebase configuration values"
    echo ""
    read -p "Press Enter to continue or Ctrl+C to cancel..."
fi

echo ""
echo "🔍 Checking Firebase configuration..."
./scripts/test-firebase.sh

echo ""
echo "🚀 Starting Docker Compose..."
echo ""

# Build and start services
docker-compose up --build

echo ""
echo "✅ Application should be running at:"
echo "   Frontend: http://localhost:3000"
echo "   Backend API: http://localhost:8080"
echo "   Database: localhost:5432"
echo ""
echo "To stop the application, press Ctrl+C or run: docker-compose down"
