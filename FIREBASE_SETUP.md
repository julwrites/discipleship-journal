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
The backend needs to know which Firebase project to use. Update the auth middleware configuration in `api/middleware/auth.go:30`:

```go
// Change this line:
config := &firebase.Config{ProjectID: "discipleship-journal-pwa"}

// To match your Firebase project ID:
config := &firebase.Config{ProjectID: "your-project-id"}
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

### Environment Variables for Production
```bash
# Production .env
VITE_FIREBASE_API_KEY=AIzaSy... (real key)
VITE_FIREBASE_AUTH_DOMAIN=your-app.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=your-project-id
# ... other Firebase config

VITE_API_URL=https://your-api-domain.com/api

# Backend (Cloud Run)
GOOGLE_CLOUD_PROJECT=your-project-id
# No GOOGLE_APPLICATION_CREDENTIALS needed (uses metadata server)
```

### Security Considerations
1. **Never commit** service account keys or `.env` files to version control
2. Use secret management (Google Secret Manager, GitHub Secrets, etc.)
3. Set proper Firebase security rules
4. Enable App Check for additional protection
5. Configure authorized domains in Firebase Console

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