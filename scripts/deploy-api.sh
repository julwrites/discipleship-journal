#!/bin/bash
set -e

# Configuration
PROJECT_ID=${PROJECT_ID:-"discipleship-journal-pwa"}
REGION=${REGION:-"us-central1"}
SERVICE_NAME="discipleship-journal-api"
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

echo "Deploying Backend to Cloud Run ($PROJECT_ID)..."

# Build and Push Container Image
gcloud builds submit --tag $IMAGE_NAME api/

# Deploy to Cloud Run
# Note: Ensure DATABASE_URL and other secrets are set in Cloud Run environment variables or Secret Manager
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --project $PROJECT_ID

echo "Backend deployment initiated."
