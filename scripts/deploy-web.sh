#!/bin/bash
set -e

echo "Deploying Frontend..."

# Build the frontend
cd web
npm install
npm run build
cd ..

# Deploy to Firebase Hosting
# Assumes firebase-tools is installed and authenticated
firebase deploy --only hosting

echo "Frontend deployed successfully."
