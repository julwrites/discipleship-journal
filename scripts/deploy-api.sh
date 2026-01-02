#!/bin/bash
set -e

# Configuration
PROJECT_ID=${PROJECT_ID:-"discipleship-journal-pwa"}
# Use GCP_REGION from GitHub Actions if available, fallback to REGION or default
REGION=${GCP_REGION:-${REGION:-"asia-southeast1"}}
# Use GCP_SERVICE_NAME from GitHub Actions if available, fallback to default
SERVICE_NAME=${GCP_SERVICE_NAME:-"discipleship-journal-api"}
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

echo "Deploying Backend to Cloud Run ($PROJECT_ID)..."

# Build and Push Container Image
gcloud builds submit --tag $IMAGE_NAME api/

# Deploy to Cloud Run
# Note: Ensure DATABASE_URL and other secrets are set in Cloud Run environment variables or Secret Manager
DEPLOY_ARGS=""
if [ -n "$GCP_SERVICE_ACCOUNT" ]; then
  echo "Using service account: $GCP_SERVICE_ACCOUNT"
  DEPLOY_ARGS="--service-account $GCP_SERVICE_ACCOUNT"
fi

if [ -n "$CLOUD_SQL_INSTANCE" ]; then
  echo "Using Cloud SQL instance: $CLOUD_SQL_INSTANCE"
  DEPLOY_ARGS="$DEPLOY_ARGS --set-cloudsql-instances $CLOUD_SQL_INSTANCE"
fi

gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --project $PROJECT_ID \
  $DEPLOY_ARGS

echo "Backend deployment initiated."
