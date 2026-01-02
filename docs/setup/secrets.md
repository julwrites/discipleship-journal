# Secrets Setup Guide

This guide explains what secrets you need to set up and where to put them for the Discipleship Journal application.

## Overview

We use individual database components (not a single DATABASE_URL). Secrets should be stored in:
1. **GitHub Secrets** - For CI/CD deployment (primary)
2. **Google Secret Manager** - Optional fallback (code has integration)
3. **Local `.env.local`** - For local development

## Required Secrets

### Database Secrets (Cloud SQL)

| Secret | Where to Find | Example Value |
|--------|---------------|---------------|
| `DB_USERNAME` | Cloud SQL → Users → Username | `myproductionuser` |
| `DB_PASSWORD` | Cloud SQL → Users → Set/Reset password | `SecurePassword123!` |
| `DB_NAME` | Cloud SQL → Databases → Database name | `discipleship_journal` |
| `CLOUD_SQL_INSTANCE` | Cloud SQL → Instance → Connection name | `myproject:asia-southeast1:myinstance` |

### Other Secrets

| Secret | Where to Find | Purpose |
|--------|---------------|---------|
| `BIBLE_API_URL` | Bible API provider | Bible API endpoint |
| `BIBLE_API_KEY` | Bible API provider | Bible API authentication |
| `GCP_SERVICE_ACCOUNT` | IAM → Service Accounts | Cloud Run service account |
| `GCP_SA_KEY` | IAM → Service Accounts → Keys | Service account JSON key |

## Step-by-Step Setup

### 1. Google Cloud SQL Setup

#### Find Cloud SQL Instance Connection Name:
1. Go to [Google Cloud Console](https://console.cloud.google.com/sql)
2. Click on your instance
3. Copy "Connection name" (format: `project:region:instance`)

#### Create Database User:
1. In Cloud SQL instance, go to "Users" tab
2. Click "Add user account"
3. Set username and password
4. Note these for `DB_USERNAME` and `DB_PASSWORD`

#### Create Database:
1. In Cloud SQL instance, go to "Databases" tab
2. Click "Create database"
3. Name it `discipleship_journal` (or your preferred name)
4. Note this for `DB_NAME`

### 2. GitHub Secrets Setup

#### Add Secrets to GitHub:
1. Go to your GitHub repository
2. Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Add each secret from the table above

#### Required GitHub Secrets:
```
DB_USERNAME=your_cloud_sql_username
DB_PASSWORD=your_cloud_sql_password
DB_NAME=discipleship_journal
CLOUD_SQL_INSTANCE=project:region:instance
BIBLE_API_URL=https://api.scripture.api.bible
BIBLE_API_KEY=your_bible_api_key
GCP_SERVICE_ACCOUNT=service-account-email@project.iam.gserviceaccount.com
GCP_SA_KEY={ "type": "service_account", ... }  # Full JSON key
```

### 3. Local Development Setup

#### Create `.env.local`:
```bash
# Copy from .env template
cp .env .env.local

# Edit .env.local with your local values
# For local Docker development, defaults should work:
DB_USERNAME=user
DB_PASSWORD=password
DB_NAME=discipleship_journal
DB_HOST=db
DB_PORT=5432
# CLOUD_SQL_INSTANCE=  # Leave empty for local
```

### 4. Verify Setup

#### Test Locally:
```bash
# Start services
docker-compose up --build

# Check health endpoint
curl http://localhost:8081/health
# Should return "OK" if database connects
```

#### Test Production Deployment:
1. Merge changes to main branch
2. GitHub Actions will deploy automatically
3. Check deployment logs in GitHub Actions
4. Verify health endpoint: `https://api.journal.tehj.io/health`

## Troubleshooting

### Common Issues:

1. **"DATABASE_CONNECTION_FAILED"**
   - Check all required secrets are set in GitHub
   - Verify Cloud SQL instance is running
   - Ensure Cloud Run service account has "Cloud SQL Client" role

2. **Authentication Errors**
   - Verify username/password in Cloud SQL Users
   - Reset password if needed

3. **Permission Denied**
   - Add Cloud Run service account as Cloud SQL user
   - Or grant "Cloud SQL Client" role to service account

4. **Instance Not Found**
   - Verify `CLOUD_SQL_INSTANCE` format: `project:region:instance`
   - Check project ID matches Firebase project

## Security Best Practices

1. **Never commit secrets** to version control
2. **Use strong passwords** for database users
3. **Regularly rotate** database passwords
4. **Minimal privileges** for database users
5. **Enable SSL/TLS** for production connections

## Connection Details

### Local Development Connection:
```
postgres://user:password@db:5432/discipleship_journal?sslmode=disable
```

### Production Cloud SQL Connection:
```
postgres://username:password@/discipleship_journal?host=/cloudsql/project:region:instance&sslmode=disable
```

The application automatically builds the correct connection string based on whether `CLOUD_SQL_INSTANCE` is set.