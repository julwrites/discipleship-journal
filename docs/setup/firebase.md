# Firebase Setup & Integration Guide

This guide walks you through setting up Firebase authentication for the Discipleship Journal application, both for local development and production.

## Table of Contents
1. [Firebase Project Setup](#firebase-project-setup)
2. [Authentication Configuration](#authentication-configuration)
3. [Frontend Configuration](#frontend-configuration)
4. [Backend Configuration](#backend-configuration)
5. [Local Development Setup](#local-development-setup)
6. [Production Deployment](#production-deployment)
7. [Testing & Troubleshooting](#testing--troubleshooting)

## Firebase Project Setup

### Step 1: Create Firebase Project
1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Click **"Add project"**
3. Enter project name: `discipleship-journal` (or your preferred name)
4. Enable Google Analytics (optional but recommended)
5. Choose Analytics location (if enabled)
6. Click **"Create project"**

### Step 2: Register Web App
1. In your Firebase project dashboard, click **"</>"** (Web) icon
2. Register your app:
   - App nickname: `Discipleship Journal Web`
   - Optional: Set up Firebase Hosting (not required for this app)
3. Click **"Register app"**
4. Copy the Firebase configuration object (you'll need this later)

### Step 3: Enable Authentication Methods
1. In Firebase Console, go to **Build** → **Authentication**
2. Click **"Get started"**
3. Go to **Sign-in method** tab
4. Enable the following providers:

#### Email/Password (Recommended)
- Click **"Email/Password"**
- Toggle **"Enable"**
- Click **"Save"**

#### Google Sign-In (Optional)
- Click **"Google"**
- Toggle **"Enable"**
- Choose a support email
- Click **"Save"**

#### Other Providers (Optional)
You can also enable:
- GitHub
- Microsoft
- Apple
- Anonymous (for guest access)

## Authentication Configuration

### Frontend Configuration Values
After registering your web app, you'll get a configuration object like:
```javascript
const firebaseConfig = {
  apiKey: "AIzaSy...",
  authDomain: "your-project.firebaseapp.com",
  projectId: "your-project-id",
  storageBucket: "your-project.appspot.com",
  messagingSenderId: "1234567890",
  appId: "1:1234567890:web:abcdef123456"
};
```

These values map to the following environment variables in your `.env` file:

| Firebase Config | Environment Variable | Example Value |
|----------------|---------------------|---------------|
| `apiKey` | `VITE_FIREBASE_API_KEY` | `AIzaSy...` |
| `authDomain` | `VITE_FIREBASE_AUTH_DOMAIN` | `your-project.firebaseapp.com` |
| `projectId` | `VITE_FIREBASE_PROJECT_ID` | `your-project-id` |
| `storageBucket` | `VITE_FIREBASE_STORAGE_BUCKET` | `your-project.appspot.com` |
| `messagingSenderId` | `VITE_FIREBASE_MESSAGING_SENDER_ID` | `1234567890` |
| `appId` | `VITE_FIREBASE_APP_ID` | `1:1234567890:web:abcdef123456` |

### Backend Configuration
For the Go backend to verify Firebase tokens, you need:

1. **Service Account Key** (for local development)
2. **Google Application Default Credentials** (for production/Cloud Run)

## Frontend Configuration

### Update Environment Variables
Edit your `.env` file with the Firebase values:

```bash
# Firebase Configuration
VITE_FIREBASE_API_KEY=AIzaSy...
VITE_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
VITE_FIREBASE_MESSAGING_SENDER_ID=1234567890
VITE_FIREBASE_APP_ID=1:1234567890:web:abcdef123456

# API URL (should point to your backend)
VITE_API_URL=http://localhost:8080/api  # For local development
# VITE_API_URL=https://your-api-domain.com/api  # For production
```

### How Authentication Works (Frontend)
1. User signs in via Firebase UI
2. Firebase returns an ID token
3. Frontend sends token to backend in `Authorization: Bearer <token>` header
4. Backend verifies token with Firebase Admin SDK
5. If valid, request proceeds; if invalid, returns 401 Unauthorized

## Backend Configuration

### Option A: Local Development with Service Account Key

#### Step 1: Generate Service Account Key
1. In Firebase Console, go to **Project Settings** → **Service accounts**
2. Click **"Generate new private key"**
3. Click **"Generate key"**
4. Save the JSON file securely (e.g., `service-account-key.json`)

#### Step 2: Configure Backend
Add the service account key to your `.env` file:

```bash
# Service Account Key (for local development)
# Copy the entire JSON content or provide path
GOOGLE_APPLICATION_CREDENTIALS=./service-account-key.json
# OR
SERVICE_ACCOUNT_KEY='{"type": "service_account", "project_id": "...", ...}'
```

#### Step 3: Update API Configuration
The backend needs to know which Firebase project to use. The code has been updated to read from the `GOOGLE_CLOUD_PROJECT` environment variable, with a fallback to a mock project ID.

If you need to change the project ID, you can either:
1. Set the `GOOGLE_CLOUD_PROJECT` environment variable to your Firebase project ID
2. Or modify the default value in `api/main.go:63`

The current configuration in `api/main.go:61-65`:
```go
firebaseProjectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
if firebaseProjectID == "" {
    firebaseProjectID = "YOUR_FIREBASE_PROJECT_ID" // Default to your actual project ID
}
firebaseService, err := services.NewFirebaseService(context.Background(), "", firebaseProjectID)
```

### Option B: Production (Cloud Run / GCP)
For production deployments on Google Cloud Platform:

1. **No service account key needed** - uses Application Default Credentials
2. Ensure the service account has **Firebase Admin SDK** permissions
3. Set the `GOOGLE_CLOUD_PROJECT` environment variable to your project ID

## Local Development Setup

### Quick Start with Mock Authentication
For quick local testing without Firebase:

1. Use the default `.env` values (already set to mock values)
2. The app will bypass Firebase authentication
3. You can simulate a logged-in user by setting in browser console:
   ```javascript
   localStorage.setItem('E2E_TEST_USER', JSON.stringify({
     uid: 'test-user-123',
     email: 'test@example.com',
     displayName: 'Test User'
   }));
   ```
4. Refresh the page

### Full Local Setup with Firebase

#### 1. Firebase Emulator Suite (Recommended for Development)
```bash
# Install Firebase CLI
npm install -g firebase-tools

# Login to Firebase
firebase login

# Initialize Firebase in your project
firebase init emulators

# Select Authentication emulator
# Choose a port (default: 9099)

# Start emulators
firebase emulators:start
```

#### 2. Update Environment Variables for Emulators
```bash
# Point to emulator
VITE_FIREBASE_AUTH_DOMAIN=localhost:9099
VITE_FIREBASE_API_KEY=fake-api-key-for-emulator

# Disable SSL for emulator (add to firebase.ts)
if (window.location.hostname === 'localhost') {
  connectAuthEmulator(auth, 'http://localhost:9099');
}
```

#### 3. Update Backend for Emulators
```go
// In auth.go, add emulator support
if os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") != "" {
    authClient, err := app.Auth(ctx)
    if err != nil {
        return nil, err
    }
    // Use emulator
}
```

## Production Deployment

### 🏗️ Architecture Overview: What Needs to be Hosted

This is a **full-stack application** with three main components:

#### 1. **Frontend (React App)** - Hosted on **Firebase Hosting**
- **What**: Static JavaScript/HTML/CSS files
- **Purpose**: User interface, authentication UI, note editing
- **Hosting**: Firebase Hosting (free tier available)
- **Data Storage**: None - it's just static files

#### 2. **Backend API (Go Server)** - Hosted on **Google Cloud Run**
- **What**: Go application that handles:
  - User management (syncs from Firebase)
  - Notes CRUD operations
  - Bible AI integration (calls external Bible API)
  - Groups, connections, sharing
  - Push notifications
- **Purpose**: Business logic, database access, external API calls
- **Hosting**: Google Cloud Run (serverless containers)
- **Database**: PostgreSQL (separate service)

#### 3. **Database (PostgreSQL)** - Hosted on **Cloud SQL**
- **What**: Relational database
- **Purpose**: Stores all application data (notes, users, groups, etc.)
- **Hosting**: Google Cloud SQL (managed PostgreSQL)

#### 4. **Authentication (Firebase Auth)** - Managed by **Firebase**
- **What**: Authentication service
- **Purpose**: User sign-in/sign-up, password management
- **Hosting**: Firebase (managed service)
- **Important**: This is NOT your database - it only handles auth

### 🔗 How They Work Together
```
User's Browser
     ↓
Firebase Hosting (React App) ← Static files only
     ↓ (API calls)
Google Cloud Run (Go Backend) ← Your custom business logic
     ↓
Cloud SQL (PostgreSQL) ← Your data storage
     ↓
External Bible API (scripture.api.bible) ← Third-party service
```

### Key Distinction: Firebase vs Your Backend
- **Firebase**: Authentication + static file hosting
- **Your Go Backend**: All your business logic + database
- **They are separate services** that communicate via API calls

### 🐳 Artifact Registry vs Container Registry

#### Why Use Artifact Registry?
- **Newer & actively developed** (Container Registry is in maintenance mode)
- **Better permissions model** (granular artifact-level permissions)
- **Multi-format support** (Docker, npm, Maven, Python, etc.)
- **Vulnerability scanning** (built-in security scanning)
- **Cleaner URLs**: `us-central1-docker.pkg.dev/project/repo/image`

#### Migration from Container Registry
The CI/CD pipeline is already configured for Artifact Registry. Changes made:
1. **Image tags**: `gcr.io/project/image` → `region-docker.pkg.dev/project/repo/image`
2. **Authentication**: `gcloud auth configure-docker` with region prefix
3. **Permissions**: `roles/artifactregistry.writer` instead of `roles/storage.objectCreator`

#### One-Time Artifact Registry Setup (REQUIRED BEFORE DEPLOYMENT)
```bash
# Create repository (use your actual repository name from GCP_ARTIFACT_REPOSITORY secret)
# This must be done BEFORE running the CI/CD pipeline!
gcloud artifacts repositories create discipleship-journal-repo \
  --repository-format=docker \
  --location=asia-southeast1  # Use your GCP_REGION

# Verify
gcloud artifacts repositories list --location=asia-southeast1

# Common error if skipped: "404 Not Found" when pushing Docker image
```

### 🔐 Required Secrets & Environment Variables

#### Frontend (Firebase Hosting) - Build-time Secrets
These are baked into the JavaScript during `npm run build`:

| Variable | Purpose | Where to Set |
|----------|---------|--------------|
| `VITE_FIREBASE_API_KEY` | Firebase API key | GitHub Secrets |
| `VITE_FIREBASE_AUTH_DOMAIN` | Firebase auth domain | GitHub Secrets |
| `VITE_FIREBASE_PROJECT_ID` | Firebase project ID | GitHub Secrets |
| `VITE_FIREBASE_STORAGE_BUCKET` | Firebase storage | GitHub Secrets |
| `VITE_FIREBASE_MESSAGING_SENDER_ID` | Firebase messaging | GitHub Secrets |
| `VITE_FIREBASE_APP_ID` | Firebase app ID | GitHub Secrets |
| `VITE_API_URL` | Backend API URL | GitHub Secrets |

**Example values** (from your `.env.local`):
```bash
VITE_FIREBASE_API_KEY=AIzaSyCQoQ4pAL4fa3IG_rmzWbWhXtgqBcTw2ns
VITE_FIREBASE_AUTH_DOMAIN=YOUR_FIREBASE_PROJECT_ID.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=YOUR_FIREBASE_PROJECT_ID
VITE_FIREBASE_STORAGE_BUCKET=YOUR_FIREBASE_PROJECT_ID.firebasestorage.app
VITE_FIREBASE_MESSAGING_SENDER_ID=995319008345
VITE_FIREBASE_APP_ID=1:995319008345:web:ba91b8f5eb3d521e98f548

# ⚠️ CRITICAL: VITE_API_URL must point to your deployed backend
# Local development:
VITE_API_URL=http://localhost:8080/api

# Production (after deploying to Cloud Run):
# Format: https://SERVICE-REGION-PROJECT_ID.a.run.app/api
# Example:
# VITE_API_URL=https://discipleship-journal-api-uc-a.run.app/api
```

#### Backend (Cloud Run) - Runtime Secrets
These are set as environment variables in Cloud Run:

| Variable | Purpose | Where to Set |
|----------|---------|--------------|
| `DATABASE_URL` | PostgreSQL connection string | Google Secret Manager |
| `GOOGLE_CLOUD_PROJECT` | GCP project ID | Cloud Run environment |
| `APP_ENV` | Set to `production` | Cloud Run environment |
| `BIBLE_API_KEY` | Bible API key (optional) | Google Secret Manager |

#### CI/CD (GitHub Actions) - Deployment Secrets
These secrets are used by GitHub Actions to deploy:

| Secret | Purpose | Example Value | Where to Set |
|--------|---------|---------------|--------------|
| `GCP_SA_KEY` | GCP Service Account JSON key | `{ "type": "service_account", ... }` | GitHub Secrets |
| `FIREBASE_SERVICE_ACCOUNT_DISCIPLESHIP_JOURNAL` | Firebase service account | `{ "type": "service_account", ... }` | GitHub Secrets |
| `GCP_PROJECT_ID` | GCP Project ID | `your-backend-gcp-project-id` (shared backend GCP project) | GitHub Secrets |
| `GCP_REGION` | GCP Region | `asia-southeast1` | GitHub Secrets |
| `GCP_SERVICE_NAME` | Cloud Run service name | `discipleship-journal-api` | GitHub Secrets |
| `GCP_ARTIFACT_REPOSITORY` | Artifact Registry repo | `discipleship-journal-repo` | GitHub Secrets |
| `VITE_FIREBASE_PROJECT_ID` | Firebase Project ID | `YOUR_FIREBASE_PROJECT_ID` | GitHub Secrets |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://...` | GitHub Secrets |

### 🚀 Step-by-Step Deployment Guide

#### 📋 Deployment Order (Important!)
You must deploy in this order:
1. **Backend first** (Cloud Run) - gets a URL like `https://service-region-project.a.run.app`
2. **Frontend second** (Firebase Hosting) - needs the backend URL to build

#### Step 1: Set Up GitHub Secrets
Go to your GitHub repository → Settings → Secrets and variables → Actions

Add these secrets:

1. **GCP Configuration Secrets** (needed first):
   - `GCP_PROJECT_ID`: Your GCP project ID (e.g., `discipleship-journal`)
   - `GCP_REGION`: GCP region (e.g., `asia-southeast1`)
   - `GCP_SERVICE_NAME`: Cloud Run service name (e.g., `discipleship-journal-api`)
   - `GCP_ARTIFACT_REPOSITORY`: Artifact Registry repo (e.g., `discipleship-journal-repo`)
   - `GCP_SA_KEY`: GCP service account JSON key (see permissions below)
   - `DATABASE_URL`: PostgreSQL connection string (from Cloud SQL)

2. **Firebase Configuration Secrets**:
   - `FIREBASE_SERVICE_ACCOUNT_DISCIPLESHIP_JOURNAL`: Firebase service account JSON
   - `VITE_FIREBASE_PROJECT_ID`: Firebase project ID (from Firebase Console)

3. **Frontend Build Secrets** (add AFTER backend is deployed):
   - Copy all `VITE_FIREBASE_*` values from your `.env.local` (except `VITE_FIREBASE_PROJECT_ID` - already set above)
   - **WAIT**: Don't add `VITE_API_URL` yet - you need the backend URL first!

#### Step 2: Set Up Google Cloud (Artifact Registry)
1. **Create Artifact Registry repository** (REQUIRED - one-time manual setup):
   ```bash
   # Replace with your actual values from secrets
   gcloud artifacts repositories create discipleship-journal-repo \
     --repository-format=docker \
     --location=asia-southeast1 \
     --description="Docker repository for Discipleship Journal API"

   # Verify it was created
   gcloud artifacts repositories list --location=asia-southeast1
   ```

   **⚠️ IMPORTANT**: The CI/CD pipeline will fail with "404 Not Found" if this repository doesn't exist!

2. **Create Cloud SQL PostgreSQL instance**
3. **Create Secret Manager secret** for `DATABASE_URL`
4. **Grant Cloud Run service account** access to the secret
5. **Deploy backend via GitHub Actions** (already configured for Artifact Registry)

#### Step 3: Find Your Backend URL
After backend deployment completes:

1. **Go to Google Cloud Console** → Cloud Run
2. **Find your service**: `discipleship-journal-api`
3. **Copy the URL**: Looks like `https://discipleship-journal-api-uc-a.run.app`
4. **Add `/api` to the end**: `https://discipleship-journal-api-uc-a.run.app/api`

#### Step 4: Add VITE_API_URL to GitHub Secrets
Now that you have the backend URL:

1. **Go back to GitHub Secrets**
2. **Add `VITE_API_URL`** with your actual URL:
   ```
   VITE_API_URL=https://discipleship-journal-api-uc-a.run.app/api
   ```
3. **Add all other frontend secrets** from `.env.local`

#### Step 5: Deploy Frontend
1. **Trigger frontend deployment** in GitHub Actions
2. **Or push to main branch** to trigger auto-deploy

#### Step 3: Configure Firebase
1. **Add authorized domains** in Firebase Console
2. **Enable App Check** for additional security
3. **Configure Firebase Hosting** (already in `firebase.json`)
4. **Set up custom domain** (optional - see below)

### 💰 Cost Considerations
- **Firebase Hosting**: Free tier (10GB storage, 360MB/day bandwidth)
- **Cloud Run**: Pay per request + compute time (free tier available)
- **Cloud SQL**: ~$10-50/month depending on size
- **Firebase Auth**: Free up to 50,000 monthly active users

### 📝 Summary: What You're Actually Hosting
1. **Static files** on Firebase Hosting (free)
2. **Go API server** on Cloud Run (pay per use)
3. **PostgreSQL database** on Cloud SQL (monthly fee)
4. **Authentication** handled by Firebase (free tier)

**Firebase is only one piece** - you need all three components running for the app to work!

### 🔗 Final URL Configuration Examples

#### Scenario 1: Default Firebase Domain
```
Frontend: https://YOUR_FIREBASE_PROJECT_ID.web.app
Backend API: https://discipleship-journal-api-uc-a.run.app/api
VITE_API_URL: https://discipleship-journal-api-uc-a.run.app/api
VITE_FIREBASE_AUTH_DOMAIN: YOUR_FIREBASE_PROJECT_ID.firebaseapp.com
```

#### Scenario 2: Custom Domain (journal.navteens.org)
```
Frontend: https://journal.navteens.org
Backend API: https://discipleship-journal-api-uc-a.run.app/api  (same as above)
VITE_API_URL: https://discipleship-journal-api-uc-a.run.app/api  (unchanged)
VITE_FIREBASE_AUTH_DOMAIN: journal.navteens.org  (changed!)
```

#### Scenario 3: Custom API Domain (requires Cloud Load Balancing)
```
Frontend: https://journal.navteens.org
Backend API: https://api.journal.navteens.org/api  (extra cost: ~$18/month)
VITE_API_URL: https://api.journal.navteens.org/api
VITE_FIREBASE_AUTH_DOMAIN: journal.navteens.org
```

### 🎯 Key Takeaway
**For most cases, use Scenario 2:**
- Custom frontend domain: `journal.navteens.org`
- Keep backend at Cloud Run URL
- Only change `VITE_FIREBASE_AUTH_DOMAIN`
- `VITE_API_URL` stays as Cloud Run URL

This gives you a professional frontend URL while keeping backend costs low.

## 🔐 GCP Service Account Permissions

The `GCP_SA_KEY` service account needs **minimal permissions** (not admin!). Here's what's actually needed:

### 🎯 Quick Recommendation (Artifact Registry)
Use these **4 roles** (NOT admin roles):
1. **`roles/run.developer`** - Can deploy/update, but NOT delete services
2. **`roles/iam.serviceAccountUser`** - Required for deployment
3. **`roles/artifactregistry.writer`** - For Artifact Registry (instead of storage.objectCreator)
4. **`roles/secretmanager.secretAccessor`** - Can read secrets only

### ❌ What You DON'T Need
- **`roles/run.admin`** - Excessive! Can delete services, manage IAM, modify billing
- **`roles/storage.objectAdmin`** - Excessive! Can delete all your Docker images
- **`roles/secretmanager.admin`** - Excessive! Can create/delete secrets

### Why Minimal Permissions Matter
- **Security**: If GitHub account is compromised, damage is limited
- **Compliance**: Follows principle of least privilege
- **Safety**: Can't accidentally delete production services
- **Auditability**: Clear what the service account can do

### Required IAM Roles (Minimal Set)
For security, use the **minimum permissions required**:

#### 1. **Cloud Run Developer** (`roles/run.developer`) - NOT Admin!
- **What it allows**: Deploy new revisions, update existing services
- **What it doesn't allow**: Delete services, manage IAM, modify billing
- **Perfect for CI/CD**: Can deploy but not destroy

#### 2. **Service Account User** (`roles/iam.serviceAccountUser`)
- **Required**: To impersonate the Cloud Run runtime service account
- **Without this**: Deployment fails with "permission denied"

#### 3. **Artifact Registry Writer** (`roles/artifactregistry.writer`)
- **For Artifact Registry** (recommended over Container Registry)
- **Allows**: Push/pull Docker images, manage artifacts
- **Better than storage.objectCreator**: More granular permissions for containers

#### 4. **Secret Manager Secret Accessor** (`roles/secretmanager.secretAccessor`)
- **Required**: To read `DATABASE_URL` from Secret Manager
- **Minimal**: Can only read, not create or delete secrets

#### Optional (if needed):
- **Cloud SQL Client** (`roles/cloudsql.client`): Only if using Cloud SQL Proxy
- **Storage Object Creator** (`roles/storage.objectCreator`): Only if using Container Registry (not recommended)

### Step-by-Step Service Account Creation

#### Option A: Using Google Cloud Console (Recommended - Minimal)
1. **Go to IAM & Admin** → **Service Accounts**
2. **Click "Create Service Account"**
3. **Name**: `github-actions-deployer`
4. **Description**: "GitHub Actions deployment (minimal permissions)"
5. **Click "Create and Continue"**
6. **Add Roles** (add one at a time - use these exact roles):
   - **Cloud Run Developer** (NOT Admin!)
   - **Service Account User**
   - **Artifact Registry Writer** (for Artifact Registry)
   - **Secret Manager Secret Accessor**
7. **Click "Done"**
8. **Go to the service account** → **Keys** → **Add Key** → **Create new key**
9. **Choose JSON** → **Create**
10. **Download the JSON file** - this is your `GCP_SA_KEY`

#### Option B: Using gcloud CLI (Minimal Permissions)
```bash
# Create service account
gcloud iam service-accounts create github-actions-deployer \
  --display-name="GitHub Actions Deployer (minimal)"

# Grant MINIMAL roles (not admin!)
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:github-actions-deployer@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.developer"  # NOT admin!

gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:github-actions-deployer@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:github-actions-deployer@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"  # For Artifact Registry

gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:github-actions-deployer@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"

# Generate key
gcloud iam service-accounts keys create gcp-sa-key.json \
  --iam-account=github-actions-deployer@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

### Even More Minimal (Custom Role)
If you want absolute minimum permissions, create a custom role:

```yaml
# Custom role name: github-actions-deployer
# These are the EXACT permissions needed (Artifact Registry):
- run.services.get
- run.services.update
- run.operations.get
- artifactregistry.repositories.uploadArtifacts
- artifactregistry.tags.list
- iam.serviceAccounts.actAs
- secretmanager.secrets.get
- secretmanager.versions.access

# What's NOT included (for safety):
# - run.services.delete (can't delete services)
# - run.services.create (can't create new services, only update existing)
# - storage.objects.delete (can't delete images)
# - storage.objects.list (can't list bucket contents)
# - secretmanager.secrets.create (can't create secrets)
```

**Note**: With this custom role, you must:
1. **Manually create the Cloud Run service first** (one-time setup)
2. **The service account can only update it**, not create new ones
3. **Even safer** than `roles/run.developer`

### Verification
Test the service account has correct permissions:
```bash
# Set the key as active
export GOOGLE_APPLICATION_CREDENTIALS=gcp-sa-key.json

# Test Cloud Run access
gcloud run services list --region=us-central1

# Test Storage access
gcloud storage ls gs://artifacts.YOUR_PROJECT_ID.appspot.com/

# Test Secret Manager access
gcloud secrets list
```

### Security Best Practices
1. **Use separate service accounts** for different environments (dev/staging/prod)
2. **Rotate keys regularly** (every 90 days recommended)
3. **Monitor usage** in Cloud Audit Logs
4. **Restrict by IP** if possible (GitHub Actions IP ranges)
5. **Delete unused service accounts**

### Troubleshooting Permission Errors
If deployment fails, check for these common issues:

- **"Permission 'run.services.update' denied"**: Missing `roles/run.developer` or custom role
- **"Permission 'artifactregistry.repositories.uploadArtifacts' denied"**: Missing `roles/artifactregistry.writer`
- **"Cannot impersonate service account"**: Missing `roles/iam.serviceAccountUser`
- **"Permission 'secretmanager.versions.access' denied"**: Missing `roles/secretmanager.secretAccessor`
- **"Service ... not found"**: Service doesn't exist yet (create it manually first if using custom role)

## 🌐 Custom Domain Setup (e.g., journal.navteens.org)

Firebase Hosting supports custom domains with automatic SSL certificates. Here's how to set it up:

### Step 1: Add Custom Domain in Firebase Console
1. Go to Firebase Console → Hosting
2. Click "Add custom domain"
3. Enter your domain: `journal.navteens.org`
4. Click "Continue"

### Step 2: Verify Domain Ownership
Firebase will provide two options:

#### Option A: TXT Record (Recommended)
1. Add a TXT record to your DNS:
   ```
   Name: @ (or leave empty)
   Type: TXT
   Value: firebase=your-project-id
   TTL: 3600 (or default)
   ```
2. Wait for DNS propagation (5-60 minutes)
3. Click "Verify" in Firebase Console

#### Option B: A Record + TXT Record
If you want to use A records instead of CNAME:
1. Add A records pointing to Firebase IPs:
   ```
   Name: journal.navteens.org
   Type: A
   Value: 199.36.158.100
   TTL: 3600
   ```
   (Add multiple A records for all Firebase IPs)
2. Add TXT record for verification

### Step 3: Update DNS Records
After verification, Firebase will provide the final DNS records:

#### For subdomain (journal.navteens.org):
```
Name: journal
Type: CNAME
Value: your-project-id.web.app
TTL: 3600
```

#### For root domain (navteens.org - if you want the whole domain):
```
Name: @
Type: A
Value: 199.36.158.100
TTL: 3600
```
(Add multiple A records for all Firebase IP addresses)

### Step 4: Update Environment Variables
Update your frontend environment variables:

```bash
# In GitHub Secrets / .env.local
VITE_FIREBASE_AUTH_DOMAIN=journal.navteens.org  # ← Changed!

# ⚠️ VITE_API_URL: Find your actual Cloud Run URL
# 1. Deploy backend first (GitHub Actions will do this)
# 2. Go to Google Cloud Console → Cloud Run
# 3. Find your service URL
# 4. Add /api to the end
# Example: https://discipleship-journal-api-uc-a.run.app/api
VITE_API_URL=https://YOUR-ACTUAL-CLOUD-RUN-URL/api
```

### Step 5: Update Firebase Authorized Domains
1. Go to Firebase Console → Authentication → Settings
2. Add `journal.navteens.org` to authorized domains
3. Remove old domain if no longer needed

### Step 6: SSL Certificate (Automatic)
Firebase automatically provisions and renews SSL certificates via Let's Encrypt. No action needed!

### DNS Configuration Examples

#### Cloudflare (Recommended):
```
Type    Name                Content                    TTL
CNAME   journal             your-project-id.web.app    Auto
TXT     @                   firebase=your-project-id   3600
```

#### Google Domains:
```
Host name: journal
Type: CNAME
Data: your-project-id.web.app
TTL: 3600
```

#### AWS Route 53:
Create record set:
- Name: journal.navteens.org
- Type: CNAME
- Value: your-project-id.web.app
- TTL: 300

### Testing Your Custom Domain
1. **Wait for DNS propagation** (up to 48 hours, usually <1 hour)
2. **Check SSL certificate**: `https://journal.navteens.org`
3. **Verify Firebase Auth works** on the new domain
4. **Test API calls** from the new domain

### Important Notes
1. **SSL certificates** are automatic and renew automatically
2. **HTTP/2 and HTTP/3** are enabled by default
3. **Global CDN** included with Firebase Hosting
4. **No additional cost** for custom domains
5. **Can have multiple domains** pointing to the same site

### Troubleshooting Custom Domains
- **DNS propagation delay**: Use `dig journal.navteens.org` to check
- **SSL certificate issues**: Wait 24 hours for auto-provisioning
- **CORS errors**: Update `VITE_FIREBASE_AUTH_DOMAIN` and Firebase authorized domains
- **Mixed content warnings**: Ensure all resources use `https://`

### Backend Considerations
If you also want a custom domain for your API (Cloud Run):
1. **Cloud Run custom domains**: Requires Cloud Load Balancing (~$18/month)
2. **Alternative**: Use API subdomain like `api.journal.navteens.org`
3. **Simpler**: Keep Cloud Run URL and update `VITE_API_URL`

## Testing & Troubleshooting

### Common Issues & Solutions

#### 1. "Firebase App named '[DEFAULT]' already exists"
- Ensure you're not initializing Firebase multiple times
- The code already handles this with singleton pattern

#### 2. "Invalid API key" or "Auth domain not authorized"
- Check your `.env` values match Firebase Console
- Ensure domain is in authorized domains list (Firebase Console → Authentication → Settings → Authorized domains)

#### 3. Backend can't verify tokens
- Verify service account has Firebase Admin SDK permissions
- Check project ID matches between frontend and backend
- For emulators: set `FIREBASE_AUTH_EMULATOR_HOST=localhost:9099`

#### 4. CORS Issues
- Backend already has CORS middleware configured
- Ensure `VITE_API_URL` matches the backend URL
- Check browser console for CORS errors

### Testing Authentication Flow

#### Manual Test
1. Start the application: `./launch.sh`
2. Open browser to `http://localhost:3000`
3. Click "Sign In"
4. Use test credentials (if using emulator) or real Firebase auth
5. Verify you can access protected routes

#### Automated Test
The project includes E2E tests with mock authentication:
```bash
cd web
npm test:e2e
```

### Debug Tools

#### Firebase Console Debug View
1. Go to Firebase Console → Authentication → Users
2. See active users and sign-in methods
3. Check logs for authentication events

#### Browser Developer Tools
```javascript
// Check Firebase initialization
console.log(window.firebase)

// Check current user
import { auth } from '@/lib/firebase'
console.log(auth.currentUser)

// Manually get ID token
const token = await auth.currentUser.getIdToken()
console.log('Token:', token)
```

#### Backend Logs
```bash
# Check backend logs for auth errors
docker-compose logs api

# Or if running locally
cd api
go run main.go
```

## Advanced Configuration

### Custom Claims & User Roles
To add custom claims (e.g., admin roles):

```javascript
// Using Firebase Admin SDK (backend)
await admin.auth().setCustomUserClaims(uid, {
  admin: true,
  premium: false
});
```

### Security Rules
Set up Firebase Security Rules for Firestore (if added later):

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /users/{userId} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
  }
}
```

### Email Templates
Customize authentication emails in Firebase Console:
- Email verification templates
- Password reset templates
- Email address change templates

## Migration from Mock to Real Firebase

### Step-by-Step Migration
1. Start with mock authentication (default)
2. Set up Firebase project and get configuration
3. Update `.env` with real values
4. Test authentication locally with emulators
5. Deploy to staging with real Firebase
6. Test thoroughly
7. Deploy to production

### Rollback Plan
If issues arise:
1. Revert to mock values in `.env`
2. The app will automatically use mock authentication
3. Debug Firebase configuration
4. Try again

## Support & Resources

- [Firebase Documentation](https://firebase.google.com/docs)
- [Firebase Admin SDK for Go](https://firebase.google.com/docs/admin/setup)
- [Firebase Emulator Suite](https://firebase.google.com/docs/emulator-suite)
- [Firebase Authentication Pricing](https://firebase.google.com/pricing) (Free tier available)

---

**Need help?** Check the project's `docs/` directory or create an issue in the repository.
