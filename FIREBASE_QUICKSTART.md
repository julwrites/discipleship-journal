# Firebase Quick Start Guide

## ⚡ 5-Minute Setup for Local Development

### Step 1: Create Firebase Project
1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Click **"Add project"** → Name it `discipleship-journal`
3. Disable Google Analytics (for simplicity)
4. Click **"Create project"**

### Step 2: Register Web App
1. Click **"</>"** (Web) icon
2. App nickname: `Discipleship Journal`
3. Click **"Register app"**
4. **Copy the configuration values** (you'll need them next)

### Step 3: Enable Authentication
1. Go to **Build** → **Authentication** → **Get started**
2. Click **"Sign-in method"** tab
3. Enable **"Email/Password"**
4. Click **"Save"**

### Step 4: Update Environment File
Edit `.env` in the project root:

```bash
# Replace these with your actual Firebase values
VITE_FIREBASE_API_KEY=AIzaSy... (from Firebase config)
VITE_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
VITE_FIREBASE_MESSAGING_SENDER_ID=1234567890
VITE_FIREBASE_APP_ID=1:1234567890:web:abcdef123456

# Keep these as-is for local development
VITE_API_URL=http://localhost:8080/api
```

### Step 5: Generate Service Account Key (for Backend)
1. In Firebase Console: **Project Settings** → **Service accounts**
2. Click **"Generate new private key"**
3. Save as `service-account-key.json` in project root
4. Add to `.env`:
   ```bash
   GOOGLE_APPLICATION_CREDENTIALS=./service-account-key.json
   ```

### Step 6: Update Backend Project ID
Edit `api/middleware/auth.go` line 30:
```go
// Change from:
config := &firebase.Config{ProjectID: "discipleship-journal-pwa"}

// To your project ID:
config := &firebase.Config{ProjectID: "your-project-id"}
```

### Step 7: Launch the Application
```bash
# Build and start
./launch.sh

# Or manually
docker-compose up --build
```

### Step 8: Test Authentication
1. Open browser to `http://localhost:3000`
2. Click **"Sign In"**
3. Create a new account with email/password
4. You should be logged in and able to use the app!

## 🚀 Production Deployment Checklist

### Frontend (Vercel/Netlify/Firebase Hosting)
```bash
# Environment variables to set:
VITE_FIREBASE_API_KEY=your_real_key
VITE_FIREBASE_AUTH_DOMAIN=your-app.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_API_URL=https://your-api-domain.com/api
# ... other Firebase config
```

### Backend (Cloud Run)
```bash
# Environment variables:
GOOGLE_CLOUD_PROJECT=your-project-id
# No service account key needed (uses metadata server)
PORT=8080
DATABASE_URL=your_production_db_url
```

## 🔧 Troubleshooting Common Issues

### "Invalid API Key"
- Check `.env` values match Firebase Console exactly
- Ensure no trailing spaces in values

### "Auth domain not authorized"
1. Firebase Console → Authentication → Settings
2. Scroll to "Authorized domains"
3. Add `localhost` (for development) and your production domain

### Backend can't verify tokens
1. Verify service account has **Firebase Admin SDK** permissions
2. Check project ID matches in all places
3. For local dev: ensure `GOOGLE_APPLICATION_CREDENTIALS` points to valid JSON

### CORS Errors
- Backend already configured with CORS middleware
- Ensure `VITE_API_URL` matches exactly (no trailing slash)

## 🎯 Quick Test Without Firebase

For quick testing without setting up Firebase:

1. Use the default `.env` values (already mock values)
2. The app will use mock authentication
3. To simulate login, in browser console:
   ```javascript
   localStorage.setItem('E2E_TEST_USER', JSON.stringify({
     uid: 'test-user-123',
     email: 'test@example.com',
     displayName: 'Test User'
   }));
   location.reload();
   ```

## 📞 Need Help?

- **Firebase Issues**: Check [Firebase Documentation](https://firebase.google.com/docs)
- **Code Issues**: Check `docs/` directory or repository issues
- **Quick Fix**: Revert to mock auth by using default `.env` values

## 🔒 Security Notes

1. **Never commit** `.env` or `service-account-key.json` to git
2. Use environment variables in production
3. Set up Firebase Security Rules if using Firestore
4. Consider enabling **App Check** for production

---

**Time to complete**: ~5-10 minutes for basic setup
**Difficulty**: Beginner to Intermediate
**Authentication Methods Supported**: Email/Password, Google, more