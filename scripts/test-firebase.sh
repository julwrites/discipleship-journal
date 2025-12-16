#!/bin/bash

# Firebase Configuration Test Script
# This script helps validate your Firebase setup

set -e

echo "========================================="
echo "Firebase Configuration Test"
echo "========================================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ No .env file found. Please run ./launch.sh first."
    exit 1
fi

# Load environment variables
source .env 2>/dev/null || echo "⚠️  Could not source .env directly, checking variables..."

echo ""
echo "🔍 Checking Frontend Configuration..."

# Check frontend Firebase variables
check_var() {
    local var_name=$1
    local var_value=${!var_name}

    if [ -z "$var_value" ]; then
        echo "❌ $var_name: NOT SET"
        return 1
    elif [[ "$var_value" == *"mock"* ]] || [[ "$var_value" == *"example"* ]] || [[ "$var_value" == *"your-"* ]]; then
        echo "⚠️  $var_name: Using placeholder value: '$var_value'"
        return 2
    else
        echo "✅ $var_name: Set (value hidden for security)"
        return 0
    fi
}

echo ""
echo "Frontend Variables:"
check_var "VITE_FIREBASE_API_KEY"
check_var "VITE_FIREBASE_AUTH_DOMAIN"
check_var "VITE_FIREBASE_PROJECT_ID"
check_var "VITE_FIREBASE_STORAGE_BUCKET"
check_var "VITE_FIREBASE_MESSAGING_SENDER_ID"
check_var "VITE_FIREBASE_APP_ID"

echo ""
echo "🔍 Checking Backend Configuration..."

# Check backend configuration
if [ -n "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
    if [ -f "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
        echo "✅ GOOGLE_APPLICATION_CREDENTIALS: File exists at $GOOGLE_APPLICATION_CREDENTIALS"

        # Check if it's a valid JSON
        if jq empty "$GOOGLE_APPLICATION_CREDENTIALS" 2>/dev/null; then
            echo "✅ Service account key: Valid JSON"
            PROJECT_ID=$(jq -r '.project_id' "$GOOGLE_APPLICATION_CREDENTIALS" 2>/dev/null)
            if [ -n "$PROJECT_ID" ]; then
                echo "✅ Project ID in key: $PROJECT_ID"
            fi
        else
            echo "❌ Service account key: Invalid JSON"
        fi
    else
        echo "❌ GOOGLE_APPLICATION_CREDENTIALS: File not found at $GOOGLE_APPLICATION_CREDENTIALS"
    fi
else
    echo "⚠️  GOOGLE_APPLICATION_CREDENTIALS: Not set (OK for production/emulators)"
fi

echo ""
echo "🔍 Checking Project ID Consistency..."

# Check project ID consistency
FRONTEND_PROJECT_ID="$VITE_FIREBASE_PROJECT_ID"
if [ -n "$GOOGLE_APPLICATION_CREDENTIALS" ] && [ -f "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
    BACKEND_PROJECT_ID=$(jq -r '.project_id' "$GOOGLE_APPLICATION_CREDENTIALS" 2>/dev/null)
    if [ -n "$BACKEND_PROJECT_ID" ] && [ -n "$FRONTEND_PROJECT_ID" ]; then
        if [ "$BACKEND_PROJECT_ID" = "$FRONTEND_PROJECT_ID" ]; then
            echo "✅ Project IDs match: $FRONTEND_PROJECT_ID"
        else
            echo "❌ Project IDs don't match!"
            echo "   Frontend: $FRONTEND_PROJECT_ID"
            echo "   Backend:  $BACKEND_PROJECT_ID"
        fi
    fi
fi

echo ""
echo "🔍 Checking API Configuration..."

# Check API URL
if [ -n "$VITE_API_URL" ]; then
    echo "✅ VITE_API_URL: $VITE_API_URL"

    # Check if it's a local URL
    if [[ "$VITE_API_URL" == *"localhost"* ]] || [[ "$VITE_API_URL" == *"127.0.0.1"* ]]; then
        echo "ℹ️  Using local API URL (for development)"
    else
        echo "ℹ️  Using remote API URL (for production)"
    fi
else
    echo "❌ VITE_API_URL: Not set"
fi

echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="

echo ""
echo "📋 Next Steps:"

# Determine what needs to be done
if [[ "$VITE_FIREBASE_API_KEY" == *"mock"* ]]; then
    echo "1. ⚠️  You're using mock authentication"
    echo "   To use real Firebase:"
    echo "   - Update .env with real Firebase values"
    echo "   - See FIREBASE_QUICKSTART.md for instructions"
else
    echo "1. ✅ Real Firebase configuration detected"
    echo "   Test login at http://localhost:3000"
fi

if [ -z "$GOOGLE_APPLICATION_CREDENTIALS" ] || [ ! -f "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
    echo "2. ⚠️  Backend authentication not configured"
    echo "   For local development:"
    echo "   - Generate service account key from Firebase Console"
    echo "   - Add GOOGLE_APPLICATION_CREDENTIALS to .env"
else
    echo "2. ✅ Backend authentication configured"
fi

echo ""
echo "3. 🚀 Launch application: ./launch.sh"
echo "4. 🌐 Open browser: http://localhost:3000"
echo "5. 🔐 Test authentication flow"

echo ""
echo "📚 Documentation:"
echo "   - FIREBASE_QUICKSTART.md - Quick setup guide"
echo "   - FIREBASE_SETUP.md - Detailed configuration guide"
echo "   - DOCKER_SETUP.md - Docker and environment setup"

echo ""
echo "========================================="
echo "Test Complete"
echo "========================================="