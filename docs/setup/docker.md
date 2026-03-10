# Docker Setup for Local Development

This document explains how to run the Discipleship Journal application locally using Docker.

## Prerequisites

- Docker and Docker Compose installed
- (Optional) Firebase account for authentication
- (Optional) Bible API key for real Bible passages

## Quick Start

1. **Clone the repository** (if not already done)
2. **Run the launch script**:
   ```bash
   ./launch.sh
   ```
   This will:
   - Check for `.env` file and create one if needed
   - Build and start all Docker services
   - Show you the URLs to access the application

## Manual Setup

If you prefer to set up manually:

1. **Create/edit the `.env` file**:
   ```bash
   cp .env .env.local
   # Edit .env.local with your values
   ```

2. **Start the services**:
   ```bash
   docker-compose up --build
   ```

## Environment Variables

The `.env` file contains all necessary configuration. Key variables:

### Database
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`: PostgreSQL credentials
- `DATABASE_URL`: Full PostgreSQL connection string

### API Backend
- `PORT`: API server port (default: 8080)
- `APP_ENV`: Environment mode (`development` or `production`)
- `BIBLE_API_URL`, `BIBLE_API_KEY`: For real Bible API integration (optional)
- `TEST_REAL_BIBLE_API`: Set to `true` to test with real Bible API

### Web Frontend (Vite)
- `VITE_API_URL`: Backend API URL (default: `http://localhost:8080/api`)
- `VITE_FIREBASE_*`: Firebase configuration (required for authentication)

## Firebase Setup

For authentication to work, you need to set up Firebase:

1. **Create a Firebase project** at https://console.firebase.google.com/
2. **Enable Authentication**:
   - Go to "Authentication" > "Sign-in method"
   - Enable "Email/Password" and/or "Google" sign-in
3. **Get Firebase configuration**:
   - Go to "Project Settings" > "General"
   - Scroll to "Your apps" section
   - If no web app exists, click "Add app" > "Web"
   - Copy the configuration values into your `.env` file

## Services

When running, you'll have:

1. **PostgreSQL Database** (`db`)
   - Port: 5432
   - Database: `discipleship_journal`
   - Credentials: From `.env` file

2. **Go API Backend** (`api`)
   - Port: 8080
   - Endpoint: `http://localhost:8080/api`
   - Automatically connects to database

3. **React Frontend** (`web`)
   - Port: 3000 (maps to 80 in container)
   - URL: `http://localhost:3000`
   - Built with Vite, served via Nginx

4. **Database Migrator** (`migrator`)
   - Runs database migrations on startup
   - Uses `migrate/migrate` image

## Development Notes

### Without Firebase
If you don't set up Firebase, the application will use mock authentication (bypass mode). This is useful for development but not for production.

### Bible API
The application can work without a real Bible API. It has a mock fallback. To use real Bible passages:
1. Get an API key from https://scripture.api.bible/ or similar service
2. Add `BIBLE_API_URL` and `BIBLE_API_KEY` to `.env`
3. Set `TEST_REAL_BIBLE_API=true` to test

### Hot Reload
For frontend development with hot reload, you may want to run the web service outside Docker:
```bash
cd web
npm install
npm run dev
```

### Database Persistence
Database data is stored in a Docker volume (`db_data`). To reset:
```bash
docker-compose down -v
```

## Troubleshooting

### "Connection refused" errors
Ensure all services are running:
```bash
docker-compose ps
```

### Build failures
Clear Docker cache:
```bash
docker-compose build --no-cache
```

### Port conflicts
Change ports in `.env` file and update `docker-compose.yml` if needed.

### Firebase errors
Check Firebase configuration in `.env` and ensure authentication is enabled in Firebase Console.

## Stopping the Application

Press `Ctrl+C` in the terminal or run:
```bash
docker-compose down
```

To also remove volumes (clears database):
```bash
docker-compose down -v
```
