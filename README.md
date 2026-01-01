# Discipleship Journal

A P2P Discipleship Journal, designed for privacy and focused mentoring.

## Features
- **Journaling**: Private notes with Bible context.
- **AI Assistant**: Bible-aware AI chat and Q&A.
- **Connections**: Connect with friends and mentors.
- **Social Groups**: Create and join groups for shared discipleship.
- **Sharing**: Share notes with groups for accountability.
- **PWA**: Installable Progressive Web App.

## Project Structure
- `api/`: Go Backend (Chi, Pgx, Postgres)
- `web/`: React Frontend (Vite, TypeScript, Tailwind)
- `docs/`: Documentation and Task Management
- `scripts/`: Utility scripts

## Getting Started

### Local Development (Docker)
To run the full stack locally:

```bash
docker-compose up --build
```

This will start:
- **Frontend**: http://localhost:3000 (Nginx proxy)
- **Backend API**: http://localhost:8080
- **Postgres DB**: localhost:5432
- **Migrations**: Auto-run on startup

### Manual Setup
See `web/README.md` for frontend-specific instructions.

## Documentation
See `docs/tasks/` for the task board and status.
See `docs/features/` for feature details.
See `docs/architecture/` for design docs.
See `SECURITY.md` for security practices and guidelines.

## Security
This project includes pre-commit hooks for secret scanning to prevent accidental commits of sensitive information. To set up:

```bash
./scripts/setup-pre-commit.sh
```

**Important**: Never commit API keys, passwords, or tokens. Always use environment variables.

## CI/CD Pipeline
This project uses GitHub Actions for continuous deployment to Google Cloud Run (backend) and Firebase Hosting (frontend).
