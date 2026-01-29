# Deployment Guide

This repository supports automated deployments to two environments: **Production** and **Staging**. Both environments are managed via GitHub Actions workflows and mirror each other in terms of configuration, with the exception of the underlying infrastructure and secrets.

## Environments

### Production
- **Trigger**: Push to `main` branch.
- **Environment Variable**: `APP_ENV=production`
- **Infrastructure**: Uses Production GCP Project and Firebase Project.
- **Workflow File**: `.github/workflows/deploy.yml`

### Staging
- **Trigger**: Push to `staging` branch.
- **Environment Variable**: `APP_ENV=production` (Mirrors production behavior)
- **Infrastructure**: Uses Staging GCP Project and Firebase Project.
- **Workflow File**: `.github/workflows/deploy-staging.yml`

## Environment Variables & Secrets

The deployment workflows rely on GitHub Secrets to configure the environment. For Staging, all secrets are prefixed with `STG_`.

### Backend / Cloud Run

These secrets are used to build the Docker image and deploy it to Google Cloud Run.

| Production Secret | Staging Secret | Description |
|-------------------|----------------|-------------|
| `FIREBASE_PROJECT_ID` | `STG_FIREBASE_PROJECT_ID` | The Google Cloud / Firebase Project ID. |
| `GCP_REGION` | `STG_GCP_REGION` | The GCP region for resources (e.g., `us-central1`). |
| `GCP_SERVICE_NAME` | `STG_GCP_SERVICE_NAME` | The name of the Cloud Run service. |
| `GCP_ARTIFACT_REPOSITORY` | `STG_GCP_ARTIFACT_REPOSITORY` | The name of the Artifact Registry repository. |
| `GCP_SA_KEY` | `STG_GCP_SA_KEY` | JSON Key for the Google Service Account used by GitHub Actions. |
| `GCP_SERVICE_ACCOUNT` | `STG_GCP_SERVICE_ACCOUNT` | The email of the Service Account to run the Cloud Run service as. |
| `CLOUD_SQL_INSTANCE` | `STG_CLOUD_SQL_INSTANCE` | The Cloud SQL connection name (e.g., `project:region:instance`). |

**Note**: The backend application also loads application-level secrets (like database credentials, API keys) directly from Google Secret Manager at runtime. Ensure that the Secret Manager in the target project (Production or Staging) is populated with the required secrets (e.g., `DB_USERNAME`, `DB_PASSWORD`, `BIBLE_API_KEY`).

### Frontend / Firebase Hosting

These secrets are injected into the frontend build process or used for Firebase Hosting deployment.

| Production Secret | Staging Secret | Description |
|-------------------|----------------|-------------|
| `FIREBASE_API_KEY` | `STG_FIREBASE_API_KEY` | Firebase API Key for the web app. |
| `FIREBASE_AUTH_DOMAIN` | `STG_FIREBASE_AUTH_DOMAIN` | Firebase Auth Domain. |
| `FIREBASE_PROJECT_ID` | `STG_FIREBASE_PROJECT_ID` | Firebase Project ID (same as Backend). |
| `FIREBASE_STORAGE_BUCKET` | `STG_FIREBASE_STORAGE_BUCKET` | Firebase Storage Bucket URL. |
| `FIREBASE_MESSAGING_SENDER_ID` | `STG_FIREBASE_MESSAGING_SENDER_ID` | Firebase Cloud Messaging Sender ID. |
| `FIREBASE_APP_ID` | `STG_FIREBASE_APP_ID` | Firebase Web App ID. |
| `FIREBASE_MEASUREMENT_ID` | `STG_FIREBASE_MEASUREMENT_ID` | Google Analytics Measurement ID. |
| `GCP_API_URL` | `STG_GCP_API_URL` | The base URL for the backend API. |
| `FIREBASE_SERVICE_ACCOUNT_DISCIPLESHIP_JOURNAL` | `STG_FIREBASE_SERVICE_ACCOUNT_DISCIPLESHIP_JOURNAL` | Service Account JSON for Firebase Hosting deployment (via `firebase-tools`). |

## Setup Instructions for Staging

To set up a new Staging environment:

1.  **Create Resources**: Ensure a separate GCP/Firebase project exists for Staging.
2.  **Secret Manager**: Populate Google Secret Manager in the Staging project with the necessary application secrets (same names as Production, but Staging values).
3.  **GitHub Secrets**: Add all `STG_` prefixed secrets listed above to the repository's GitHub Secrets.
4.  **Deploy**: Push code to the `staging` branch to trigger the initial deployment.
